package cache

import (
	"context"
	"sia/backend/lib"
	"sia/backend/types"
	"sync"
	"sync/atomic"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
)

// Constants for low-resource server
const (
	MaxEcgSamples        = 50000  // ~3 min at 250 Hz, ~400 KB
	RistrettoMaxCost     = 500000 // Reduced for 4 GB RAM
	RistrettoNumCounters = 1000   // Reduced for lower memory
	BatchSize            = 10     // Batch ECG updates to reduce lock/cache overhead
)

// Cache manages ECG data with a lock-free ring buffer
type Cache struct {
	cache     *cache.Cache[[]float64]
	keyPrefix string
	ringBuf   []float64    // Pre-allocated ring buffer
	ringPos   uint64       // Current position (atomic)
	ringLen   uint64       // Current length (atomic)
	batch     []float64    // Buffer for batching writes
	batchIdx  int          // Current batch index
	mu        sync.RWMutex // For cache Set operations
	ecgTotalSamples uint64 
}

// Config manages WebSocket configuration
type Config struct {
	cache        *cache.Cache[types.WebSocketConfigResponse]
	cachedConfig types.WebSocketConfigResponse // Local cache
	mu           sync.RWMutex                  // Protect cachedConfig
}

// CreateNewCache initializes Cache and Config
func CreateNewCache() (*Cache, *Config, error) {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: RistrettoNumCounters, // Lowered for memory efficiency
		MaxCost:     RistrettoMaxCost,     // Lowered for 4 GB RAM
		BufferItems: 64,
	})
	if err != nil {
		return nil, nil, err
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)

	cacheManager := cache.New[[]float64](ristrettoStore)
	configManager := cache.New[types.WebSocketConfigResponse](ristrettoStore)

	// Initialize ECG cache
	ringBuf := make([]float64, MaxEcgSamples)
	c := &Cache{
		cache:     cacheManager,
		keyPrefix: "ecg",
		ringBuf:   ringBuf,
		batch:     make([]float64, BatchSize),
	}

	// Initialize config
	defaultConfig := types.WebSocketConfigResponse{
		ChunksSize: lib.CHUNK_SIZE,
		Priotize:   types.PRIO_ECG,
	}

	cf := &Config{
		cache:        configManager,
		cachedConfig: defaultConfig,
	}

	// Set initial cache values
	if err := c.cache.Set(context.Background(), c.keyPrefix, ringBuf[:0]); err != nil {
		return nil, nil, err
	}
	if err := cf.cache.Set(context.Background(), "config", defaultConfig); err != nil {
		return nil, nil, err
	}

	go lib.Print(lib.CACHE_SERVICE, "Cache Manager created")
	return c, cf, nil
}

// AddIndexToEcg appends an ECG value with batching
func (c *Cache) AddIndexToEcg(ctx context.Context, index float64) ([]float64, error) {
	// Batch incoming values
	c.batch[c.batchIdx] = index
	c.batchIdx++

	if c.batchIdx < BatchSize {
		return c.getCurrentView(), nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for i := 0; i < c.batchIdx; i++ {
		ringPos := atomic.AddUint64(&c.ringPos, 1) - 1
		ringIdx := ringPos % uint64(len(c.ringBuf))
		c.ringBuf[ringIdx] = c.batch[i]
		if ringLen := atomic.AddUint64(&c.ringLen, 1); ringLen > uint64(len(c.ringBuf)) {
			atomic.StoreUint64(&c.ringLen, uint64(len(c.ringBuf)))
		}

		// Track total ECG samples and clear if needed
		total := atomic.AddUint64(&c.ecgTotalSamples, 1)
		if total >= 45000 {
			lib.Print(lib.CACHE_SERVICE, "Reached 45000 ECG samples, clearing cache")
			c.clearUnsafe(ctx)
			atomic.StoreUint64(&c.ecgTotalSamples, 0)
			break // Stop further processing this batch, as buffer was reset
		}
	}

	c.batchIdx = 0
	result := c.getCurrentView()
	if err := c.cache.Set(ctx, c.keyPrefix, result); err != nil {
		return nil, err
	}
	return result, nil
}

// getCurrentView returns a view of the ring buffer without copying
func (c *Cache) getCurrentView() []float64 {
	ringLen := atomic.LoadUint64(&c.ringLen)
	ringPos := atomic.LoadUint64(&c.ringPos)
	if ringLen == 0 {
		return c.ringBuf[:0]
	}
	if ringLen >= uint64(len(c.ringBuf)) {
		start := ringPos % uint64(len(c.ringBuf))
		return append(c.ringBuf[start:], c.ringBuf[:start]...) // Only copy if necessary
	}
	return c.ringBuf[:ringLen]
}

// GetEcgArray retrieves the current ECG array
func (c *Cache) GetEcgArray(ctx context.Context) ([]float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Try local view first
	result := c.getCurrentView()
	cacheVal, err := c.cache.Get(ctx, c.keyPrefix)
	if err != nil {
		return result, nil // Return local view on cache miss
	}
	return cacheVal, nil
}

// GetLength returns the current ECG array length
func (c *Cache) GetLength(ctx context.Context) (int, error) {
	return int(atomic.LoadUint64(&c.ringLen)), nil
}

// ClearValues resets the ECG cache
func (c *Cache) ClearValues(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clearUnsafe(ctx)
	return nil
}

// GetConfig retrieves the cached configuration
func (c *Config) GetConfig(ctx context.Context) (*types.WebSocketConfigResponse, error) {
	c.mu.RLock()
	config := c.cachedConfig
	c.mu.RUnlock()

	cacheConfig, err := c.cache.Get(ctx, "config")
	if err != nil {
		return &config, nil // Return local cache
	}
	if cacheConfig != config {
		c.mu.Lock()
		c.cachedConfig = cacheConfig
		c.mu.Unlock()
	}
	return &cacheConfig, nil
}

// UpdateConfig updates the configuration
func (c *Config) UpdateConfig(ctx context.Context, newConfig types.WebSocketConfigResponse) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.cache.Set(ctx, "config", newConfig); err != nil {
		return err
	}
	c.cachedConfig = newConfig
	return nil
}

func (c *Cache) clearUnsafe(ctx context.Context) {
	c.ringBuf = make([]float64, len(c.ringBuf))
	atomic.StoreUint64(&c.ringPos, 0)
	atomic.StoreUint64(&c.ringLen, 0)
	c.batchIdx = 0
	c.cache.Set(ctx, c.keyPrefix, c.ringBuf[:0])
}
