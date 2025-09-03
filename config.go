package main

import "sync/atomic"

type apiConfig struct {
	FileserverHits atomic.Int32
}
