package api

import (
	"sync/atomic"

	"github.com/crisp-coder/chirpy/internal/database"
)

type ApiConfig struct {
	Db             *database.Queries
	FileserverHits atomic.Int32
}
