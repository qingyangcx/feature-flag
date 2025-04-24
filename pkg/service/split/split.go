package split

import (
	"feature-flag/pkg/constant"
	"feature-flag/pkg/logger"
	"feature-flag/pkg/model"
	"fmt"

	"github.com/spaolacci/murmur3"
)

type SplitService struct {
}

func NewSplitService() *SplitService {
	return &SplitService{}
}

func (s *SplitService) SplitAll(identity string, featureMap map[uint64]*model.Feature) []*model.SplitValue {
	splitValues := make([]*model.SplitValue, 0)
	for _, feature := range featureMap {
		splitValue := s.SplitOne(identity, feature)
		if splitValue == nil {
			logger.Logger.Errorf("split value nil, identity=%s,feature=%+v", identity, feature)
			continue
		}
		splitValues = append(splitValues, splitValue)
	}
	return splitValues
}

func (s *SplitService) SplitOne(identity string, feature *model.Feature) (splitValue *model.SplitValue) {
	if feature == nil || feature.DefaultValue == nil {
		return
	}
	splitValue = &model.SplitValue{}
	splitValue.SetValue(feature.DefaultValue)
	splitValue.Reason = "hit nothing"
	if s.blacklistFilter(identity, feature.BlacklistSet) {
		splitValue.Reason = fmt.Sprintf("hit blacklist=%v", feature.Blacklist)
		splitValue.ReasonType = constant.ReasonTypeBlacklist
		return
	}
	featureValue := s.whitelistFilter(identity, feature.FeatureValues)
	if featureValue != nil {
		splitValue.SetValue(featureValue)
		splitValue.Reason = fmt.Sprintf("hit whitelist=%v", featureValue.Whitelist)
		splitValue.ReasonType = constant.ReasonTypeWhitelist
		return
	}
	featureValue = s.bucketFilter(identity, feature.FeatureValues, feature.Seed)
	if featureValue != nil {
		splitValue.SetValue(featureValue)
		splitValue.Reason = "hit by hash"
		splitValue.ReasonType = constant.ReasonTypeHash
		return
	}
	return
}

func (s *SplitService) blacklistFilter(identity string, blacklistSet map[string]struct{}) bool {
	if _, ok := blacklistSet[identity]; ok {
		return true
	}
	return false
}

func (s *SplitService) whitelistFilter(identity string, featureValues []*model.FeatureValue) *model.FeatureValue {
	for _, featureValue := range featureValues {
		if _, ok := featureValue.WhitelistSet[identity]; ok {
			return featureValue
		}
	}
	return nil
}

func (s *SplitService) bucketFilter(identity string, featureValues []*model.FeatureValue, seed uint32) *model.FeatureValue {
	hashValue := murmur3.Sum32WithSeed([]byte(identity), seed) % constant.TotalTraffic
	var upperBound uint32 = 0
	for _, featureValue := range featureValues {
		upperBound += *featureValue.Traffic
		if hashValue < upperBound {
			return featureValue
		}
	}
	return nil
}
