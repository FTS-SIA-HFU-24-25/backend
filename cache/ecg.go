package cache

import (
	"context"
	"sia/backend/lib"
	"sia/backend/types"
	"sync"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
)

type Cache struct {
	*cache.Cache[[]float64]
	keyPrefix string
	mu        sync.Mutex
}

type Config struct {
	*cache.Cache[types.WebSocketConfigResponse]
	ChunkSize int
}

func CreateNewCache() (*Cache, *Config) {
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10000,
		MaxCost:     1000000,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)

	cacheManager := cache.New[[]float64](ristrettoStore)
	configManager := cache.New[types.WebSocketConfigResponse](ristrettoStore)

	c := &Cache{
		Cache:     cacheManager,
		keyPrefix: "ecg",
	}

	cf := &Config{
		Cache:     configManager,
		ChunkSize: lib.CHUNK_SIZE,
	}

	err = c.Set(context.Background(), c.keyPrefix, make([]float64, 0))
	if err != nil {
		panic(err)
	}

	_ = cf.Set(context.Background(), "config", types.WebSocketConfigResponse{
		ChunksSize:       lib.CHUNK_SIZE,
		StartReceiveData: 0,
		FilterType:       0,
		MaxPass:          0,
		MinPass:          0,
	})

	lib.Print(lib.CACHE_SERVICE, "Cache Manager created")

	return c, cf
}

func (c *Cache) AddIndexToEcg(ctx context.Context, index float64) (*[]float64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	arr, err := c.GetEcgArray(ctx)
	if err != nil {
		return nil, err
	}

	newArr := append(*arr, index)
	if len(newArr) > lib.ECG_HZ*60*7 {
		newArr = newArr[1:]
	}
	err = c.Set(ctx, c.keyPrefix, newArr)
	if err != nil {
		return nil, err
	}

	return &newArr, nil
}

func (c *Cache) GetEcgArray(ctx context.Context) (*[]float64, error) {
	arr, err := c.Get(ctx, c.keyPrefix)
	if err != nil {
		lib.Print(lib.CACHE_SERVICE, "Find error")
		return nil, err
	}

	return &arr, nil
}

func (c *Cache) GetLength(ctx context.Context) (int, error) {
	arr, err := c.GetEcgArray(ctx)
	if err != nil {
		return 0, err
	}

	return len(*arr), nil
}

func (c *Cache) ClearValues(ctx context.Context) error {
	err := c.Delete(ctx, c.keyPrefix)
	if err != nil {
		return err
	}

	err = c.Set(context.Background(), c.keyPrefix, make([]float64, 0))
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) GetConfig(ctx context.Context) (*types.WebSocketConfigResponse, error) {
	config, err := c.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) UpdateConfig(ctx context.Context, newConfig types.WebSocketConfigResponse) error {
	err := c.Set(ctx, "config", newConfig)
	if err != nil {
		return err
	}

	c.ChunkSize = newConfig.ChunksSize

	return nil
}
