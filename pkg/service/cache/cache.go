package cache

import (
	"feature-flag/pkg/dao"
	"feature-flag/pkg/logger"
	"feature-flag/pkg/model"
	"sync/atomic"
	"time"

	"github.com/spaolacci/murmur3"
)

type Namespace2ID2Feature map[string]map[uint64]*model.Feature

type Cache struct {
	featureMap atomic.Pointer[Namespace2ID2Feature]
	dao        *dao.Dao
}

func NewCache() *Cache {
	cache := &Cache{featureMap: atomic.Pointer[Namespace2ID2Feature]{}, dao: dao.NewDao()}
	go cache.memoryFlush()
	return cache
}

func (c *Cache) Get() Namespace2ID2Feature {
	return *c.featureMap.Load()
}

func (c *Cache) Update(featureMap Namespace2ID2Feature) {
	c.featureMap.Store(&featureMap)
}

func (c *Cache) UpdateFromFeatures(features []*model.Feature) {
	featureMap := make(Namespace2ID2Feature)
	for _, feature := range features {
		feature.Seed = murmur3.Sum32([]byte(feature.Key))
		c.fillWhitelistAndBlacklist(feature)
		if _, ok := featureMap[feature.Namespace]; !ok {
			featureMap[feature.Namespace] = make(map[uint64]*model.Feature)
		}
		featureMap[feature.Namespace][feature.ID] = feature
	}
	c.Update(featureMap)
}

func (c *Cache) fillWhitelistAndBlacklist(feature *model.Feature) {
	feature.BlacklistSet = make(map[string]struct{})
	for _, black := range feature.Blacklist {
		feature.BlacklistSet[black] = struct{}{}
	}
	for _, featureValue := range feature.FeatureValues {
		if featureValue.Default != nil && *featureValue.Default {
			feature.DefaultValue = featureValue
		}
		featureValue.WhitelistSet = make(map[string]struct{})
		for _, white := range featureValue.Whitelist {
			featureValue.WhitelistSet[white] = struct{}{}
		}
	}
}

func (c *Cache) memoryFlush() {
	for {
		time.Sleep(time.Second)
		features, err := c.dao.GetFeatureList1(map[string]interface{}{"valid": true})
		if err != nil {
			logger.Logger.Errorf("get feature list failed: %v", err)
			continue
		}
		c.UpdateFromFeatures(features)
	}
}
