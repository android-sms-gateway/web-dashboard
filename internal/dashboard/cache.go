package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-core-fx/cachefx/cache"
	"go.uber.org/zap"
)

const trendsCacheTTL = 60 * time.Second

type trendsCache struct {
	cache  cache.Cache
	logger *zap.Logger
}

func newTrendsCache(c cache.Cache, logger *zap.Logger) *trendsCache {
	return &trendsCache{
		cache:  c,
		logger: logger,
	}
}

func (c *trendsCache) get(ctx context.Context, userID string, days int) (Trends, bool) {
	raw, err := c.cache.Get(ctx, trendsCacheKey(userID, days))
	if err != nil {
		if !errors.Is(err, cache.ErrKeyNotFound) {
			c.logger.Debug("failed to get trends from cache", zap.Error(err))
		}
		return Trends{}, false
	}

	var data Trends
	if err = json.Unmarshal(raw, &data); err != nil {
		c.logger.Debug("failed to unmarshal cached trends", zap.Error(err))
		return Trends{}, false
	}

	return data, true
}

func (c *trendsCache) set(ctx context.Context, userID string, days int, data Trends) {
	raw, err := json.Marshal(data)
	if err != nil {
		c.logger.Debug("failed to marshal trends", zap.Error(err))
		return
	}

	if err = c.cache.Set(ctx, trendsCacheKey(userID, days), raw, cache.WithTTL(trendsCacheTTL)); err != nil {
		c.logger.Debug("failed to set trends cache", zap.Error(err))
	}
}

func trendsCacheKey(userID string, days int) string {
	return userID + ":" + strconv.Itoa(days)
}
