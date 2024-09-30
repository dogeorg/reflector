package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type Database struct {
	cache *cache.Cache
	mu    sync.Mutex
}

func NewDatabase() (*Database, error) {
	c := cache.New(2*time.Hour, 10*time.Minute)
	return &Database{cache: c}, nil
}

func (db *Database) Close() error {
	// No need to close in-memory cache
	return nil
}

func (db *Database) SaveEntry(token, ip string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.cache.Set(token, ip, cache.DefaultExpiration)
	return nil
}

func (db *Database) DeleteEntry(token string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.cache.Delete(token)
	return nil
}

func (db *Database) GetIP(token string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	ip, found := db.cache.Get(token)
	if !found {
		return "", fmt.Errorf("IP not found for token: %s", token)
	}
	return ip.(string), nil
}
