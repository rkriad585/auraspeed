package speedtest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"auraspeed/internal/config"
	"auraspeed/internal/logging"

	st "github.com/showwin/speedtest-go/speedtest"
)

// defaultCacheDuration is the default time the server cache remains valid.
const defaultCacheDuration = 1 * time.Hour

// CacheTTL is the configurable cache TTL in seconds
var CacheTTL = 3600 // default 1 hour

var cacheMutex sync.RWMutex

// CachedServer represents a speed test server stored in cache.
type CachedServer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Country string `json:"country"`
	Sponsor string `json:"sponsor"`
	Host    string `json:"host"`
}

// ServerCache holds cached server list with timestamp.
type ServerCache struct {
	Timestamp time.Time      `json:"timestamp"`
	Servers   []CachedServer `json:"servers"`
}

// getCacheFilePath returns the path to the server cache file.
func getCacheFilePath() string {
	return filepath.Join(config.GetDataDir(), "servers.json")
}

// loadServerCache loads the server cache from disk.
// Returns nil if cache doesn't exist or is corrupted.
func loadServerCache() (*ServerCache, error) {
	cacheFile := getCacheFilePath()
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var cache ServerCache
	if err := json.Unmarshal(data, &cache); err != nil {
		os.Remove(cacheFile)
		return nil, err
	}

	return &cache, nil
}

// saveServerCache saves the server list to disk cache.
func saveServerCache(servers st.Servers) error {
	cache := ServerCache{
		Timestamp: time.Now(),
		Servers:   make([]CachedServer, len(servers)),
	}

	for i, s := range servers {
		cache.Servers[i] = CachedServer{
			ID:      s.ID,
			Name:    s.Name,
			URL:     s.URL,
			Country: s.Country,
			Sponsor: s.Sponsor,
			Host:    s.Host,
		}
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	cacheFile := getCacheFilePath()
	return os.WriteFile(cacheFile, data, 0644)
}

// isCacheValid checks if the cache is still within its validity period.
func isCacheValid(cache *ServerCache) bool {
	cacheMutex.RLock()
	ttl := time.Duration(CacheTTL) * time.Second
	cacheMutex.RUnlock()
	return time.Since(cache.Timestamp) < ttl
}

// convertToSpeedtestServers converts cached servers back to speedtest-go format.
func convertToSpeedtestServers(cachedServers []CachedServer) st.Servers {
	servers := make(st.Servers, len(cachedServers))
	for i, cs := range cachedServers {
		id := cs.ID
		servers[i] = &st.Server{
			ID:      id,
			Name:    cs.Name,
			URL:     cs.URL,
			Country: cs.Country,
			Sponsor: cs.Sponsor,
			Host:    cs.Host,
		}
	}
	return servers
}

// GetCachedServers returns cached servers if available and valid.
func GetCachedServers() (st.Servers, bool) {
	cache, err := loadServerCache()
	if err != nil || !isCacheValid(cache) {
		return nil, false
	}
	return convertToSpeedtestServers(cache.Servers), true
}

// SaveServersToCache saves servers to cache file.
func SaveServersToCache(servers st.Servers) error {
	return saveServerCache(servers)
}

// FetchServersWithCache returns servers from cache or fetches fresh ones.
// It saves to cache after fetching.
func FetchServersWithCache(client *st.Speedtest) (st.Servers, error) {
	cache, err := loadServerCache()
	if err == nil && isCacheValid(cache) {
		return convertToSpeedtestServers(cache.Servers), nil
	}

	servers, err := client.FetchServers()
	if err != nil {
		return nil, err
	}

	if err := saveServerCache(servers); err != nil {
		logging.Get().ErrorWithFields("Failed to save server cache", map[string]interface{}{"error": err})
	}

	return servers, nil
}
