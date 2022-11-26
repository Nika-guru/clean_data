package cache

import (
	"review-service/pkg/server"
	"strings"
)

// Initialize Function in Cache Package
func init() {
	// Local Cache Configuration Value
	localCacheCfg.NumCounters = server.Config.GetInt64("LOCAL_CACHE_NUM_COUNTERS")
	localCacheCfg.MaxCost = server.Config.GetInt64("LOCAL_CACHE_MAX_COST")
	localCacheCfg.BufferItems = server.Config.GetInt64("LOCAL_CACHE_BUFFER_ITEMS")
	localCacheCfg.Metrics = server.Config.GetBool("LOCAL_CACHE_METRICS")

	if localCacheCfg.MaxCost != 0 && localCacheCfg.BufferItems != 0 && localCacheCfg.NumCounters != 0 {
		// Do Redis Cache Connection
		LocalCache = localCacheConnect()
	}

	// Remote Cache Configuration Value
	switch strings.ToLower(server.Config.GetString("REMOTE_CACHE_DRIVER")) {
	case "redis":
		server.Config.SetDefault("REMOTE_CACHE_PORT", "6379")

		redisCfg.Host = server.Config.GetString("REMOTE_CACHE_HOST")
		redisCfg.Port = server.Config.GetString("REMOTE_CACHE_PORT")
		redisCfg.Password = server.Config.GetString("REMOTE_CACHE_PASSWORD")
		redisCfg.Name = server.Config.GetInt("REMOTE_CACHE_NAME")

		if len(redisCfg.Host) != 0 && len(redisCfg.Port) != 0 {

			// Do Redis Cache Connection
			Redis = redisConnect()
		}
	}
}
