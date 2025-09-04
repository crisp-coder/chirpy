package main

import (
	"sync/atomic"

	"github.com/crisp-coder/chirpy/internal/database"
)

type apiConfig struct {
	db             *database.Queries
	FileserverHits atomic.Int32
}
