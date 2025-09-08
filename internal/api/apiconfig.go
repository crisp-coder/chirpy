package api

import (
	"sync/atomic"

	"github.com/crisp-coder/chirpy/internal/database"
)

type ApiConfig struct {
	Db             *database.Queries
	JWT_SECRET     string
	POLKA_KEY      string
	FileserverHits atomic.Int32
}
