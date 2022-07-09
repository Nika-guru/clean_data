package cache

import (
	"base/pkg/server"
	"strings"
)

// Initialize Function in Cache Package
func init() {
	// Cache Configuration Value
	switch strings.ToLower(server.Config.GetString("CACHE_DRIVER")) {
	case "redis":
		server.Config.SetDefault("CACHE_PORT", "6379")

		redisCfg.Host = server.Config.GetString("CACHE_HOST")
		redisCfg.Port = server.Config.GetString("CACHE_PORT")
		redisCfg.Password = server.Config.GetString("CACHE_PASSWORD")
		redisCfg.Name = server.Config.GetInt("CACHE_NAME")

		if len(redisCfg.Host) != 0 && len(redisCfg.Port) != 0 {

			// Do Redis Cache Connection
			Redis = redisConnect()
		}
	}
}
