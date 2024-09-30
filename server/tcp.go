package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sia/backend/db"
	"sia/backend/models"
	"sia/backend/models/types"
	"sia/backend/tools"
	"sia/backend/uuid"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	va = validator.New()
)

func OpenTCPServer() {
	http.HandleFunc("/connect", handleConnection)
	s := &http.Server{
		Addr: ":8080",
	}
	tools.Log("[TCP]", "Server running on port 8080")
	s.ListenAndServe()
}

type connectionRequest struct {
	HardwareID int64          `json:"hardware_id" validate:"required,min=1,max=255"`
	Type       types.DataType `json:"type" validate:"required,min=0"`
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	tools.Log("[TCP]", "Handling connection...")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req connectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		tools.Log("[TCP]", "Failed to parse request body: "+err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := va.Struct(req); err != nil {
		tools.Log("[TCP]", "Failed to validate request: "+err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uuid := uuid.GenerateUUID(req.HardwareID)
	con := models.Connection{
		Uuid:      uuid,
		CreatedAt: time.Now(),
		Data:      []byte{},
		Type:      req.Type,
	}

	err := db.RedisDB.Set(fmt.Sprintf("%x", uuid), con)
	if err != tools.OK {
		tools.Log("[TCP]", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return uuid
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		ID [8]byte `json:"id"`
	}{
		ID: uuid,
	})
}
