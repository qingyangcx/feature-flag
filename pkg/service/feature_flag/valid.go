package feature_flag

import (
	"feature-flag/pkg/constant"
	"feature-flag/pkg/model"
	"fmt"
)

func ValidateFeatureCreate(feature *model.Feature) model.Error {
	if feature.Name == "" {
		return model.Error{Code: constant.FeatureNameEmpty, Reason: "feature name empty"}
	}
	if feature.Key == "" {
		return model.Error{Code: constant.FeatureKeyEmpty, Reason: "feature key empty"}
	}

	if len(feature.FeatureValues) == 0 {
		return model.Error{Code: constant.FeatureValueNotExist, Reason: "feature values not exist"}
	}
	var trafficSum uint32 = 0
	var defaultValueExist bool = false
	var valueDefaultValue bool = false
	for _, value := range feature.FeatureValues {
		trafficSum += *value.Traffic
		if value.Default == nil {
			value.Default = &valueDefaultValue
		}
		if *value.Default {
			defaultValueExist = true
		}
	}
	if !defaultValueExist {
		return model.Error{Code: constant.LackDefaultFeatureValue, Reason: "need mark one feature value default"}
	}
	if trafficSum > constant.TotalTraffic {
		return model.Error{Code: constant.TrafficOverflow, Reason: fmt.Sprintf("total traffic=%d exceeds limit=%d", trafficSum, constant.TrafficOverflow)}
	}
	return model.Error{Code: constant.OK}
}

func ValidateTraffic(featureValue *model.FeatureValue, feature *model.Feature) model.Error {
	if featureValue == nil || feature == nil {
		return model.Error{Code: constant.Error, Reason: "invalid feature"}
	}
	var trafficSum uint32 = 0
	for _, value := range feature.FeatureValues {
		if value.ID == featureValue.ID {
			trafficSum += *featureValue.Traffic
		} else {
			trafficSum += *value.Traffic
		}
	}
	if trafficSum > constant.TotalTraffic {
		return model.Error{Code: constant.TrafficOverflow, Reason: fmt.Sprintf("total traffic=%d exceeds limit=%d", trafficSum, constant.TotalTraffic)}
	}
	return model.Error{Code: constant.OK}
}

func ValidateWhitelist(featureValue *model.FeatureValue, feature *model.Feature) model.Error {
	if featureValue == nil || feature == nil {
		return model.Error{Code: constant.Error, Reason: "invalid feature"}
	}
	whitelist := make(map[string]struct{})
	for _, value := range feature.FeatureValues {
		if value.ID != featureValue.ID {
			for _, w := range value.Whitelist {
				whitelist[w] = struct{}{}
			}
		}
	}
	for _, w := range featureValue.Whitelist {
		if _, ok := whitelist[w]; ok {
			return model.Error{Code: constant.Error, Reason: fmt.Sprintf("whitelist=%s conflicts,need remove", w)}
		}
	}
	return model.Error{Code: constant.OK}
}

func ValidateFeatureValueCreate(feature *model.Feature, featureValue *model.FeatureValue) (errorModel model.Error) {
	errorModel = model.Error{Code: constant.OK}
	if feature == nil || featureValue == nil || featureValue.Traffic == nil {
		errorModel.Code = constant.Error
		errorModel.Reason = "feature or feature value invalid"
		return
	}
	if featureValue.Default != nil && *featureValue.Default {
		errorModel.Code = constant.Error
		errorModel.Reason = "can not create default value"
		return
	}
	trafficCurr := uint32(0)
	for _, value := range feature.FeatureValues {
		if value.Value == featureValue.Value {
			errorModel.Code = constant.Error
			errorModel.Reason = fmt.Sprintf("feature value=%s already exists", value.Value)
			return
		}
		if value.Traffic != nil {
			trafficCurr += *value.Traffic
		}
	}
	if trafficCurr+*featureValue.Traffic > constant.TotalTraffic {
		return model.Error{Code: constant.TrafficOverflow, Reason: fmt.Sprintf("total traffic=%d exceeds limit=%d", trafficCurr, constant.TotalTraffic)}
	}
	return ValidateWhitelist(featureValue, feature)
}
