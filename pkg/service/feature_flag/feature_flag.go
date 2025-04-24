package feature_flag

import (
	"context"
	"encoding/json"
	"feature-flag/pkg/constant"
	"feature-flag/pkg/dao"
	"feature-flag/pkg/logger"
	"feature-flag/pkg/model"
	"fmt"
)

type FeatureService struct {
	dao *dao.Dao
}

func NewFeatureService() *FeatureService {
	return &FeatureService{dao: dao.NewDao()}
}

func (s *FeatureService) CreateFeatureFlag(ctx context.Context, feature *model.Feature) (*uint64, model.Error) {
	errModel := ValidateFeatureCreate(feature)
	if errModel.Code != constant.OK {
		return nil, errModel
	}
	err := s.dao.CreateFeature(ctx, feature)
	if err != nil {
		return nil, model.Error{Code: constant.Error, Reason: err.Error()}
	}
	return &feature.ID, model.Error{Code: constant.OK}
}

func (s *FeatureService) ModifyFeatureFlag(ctx context.Context, feature *model.Feature) model.Error {
	err := s.dao.ModifyFeature(ctx, feature)
	if err != nil {
		return model.Error{Code: constant.Error, Reason: err.Error()}
	}
	return model.Error{Code: constant.OK}
}

func (s *FeatureService) ModifyFeatureValue(ctx context.Context, featureValue *model.FeatureValue) model.Error {
	feature := &model.Feature{}
	feature.ID = *featureValue.FeatureID
	err := s.dao.GetFeature(feature)
	if err != nil {
		return model.Error{Code: constant.Error, Reason: err.Error()}
	}
	conditions := make(map[string]interface{})
	if featureValue.Traffic != nil {
		errModel := ValidateTraffic(featureValue, feature)
		if errModel.Code != constant.OK {
			return errModel
		}
		conditions["traffic"] = *featureValue.Traffic
	}
	if featureValue.Value != "" {
		conditions["feature_value"] = featureValue.Value
	}
	if featureValue.Whitelist != nil {
		errModel := ValidateWhitelist(featureValue, feature)
		if errModel.Code != constant.OK {
			return errModel
		}
		whitelist, err := json.Marshal(featureValue.Whitelist)
		if err != nil {
			return model.Error{Code: constant.Error, Reason: err.Error()}
		}
		conditions["whitelist"] = whitelist
	}
	err = s.dao.ModifyFeatureValue(ctx, featureValue, conditions)
	if err != nil {
		return model.Error{Code: constant.Error, Reason: err.Error()}
	}
	return model.Error{Code: constant.OK}
}

func (s *FeatureService) GetFeature(ctx context.Context, id uint64) (*model.Feature, model.Error) {
	feature := &model.Feature{}
	feature.ID = id
	err := s.dao.GetFeature(feature)
	if err != nil {
		return nil, model.Error{Code: constant.Error, Reason: err.Error()}
	}
	return feature, model.Error{Code: constant.OK}
}

func (s *FeatureService) GetFeatureList(ctx context.Context, req *model.GetFeatureListReq) (features []*model.Feature, errorModel model.Error) {
	errorModel.Code = constant.OK
	conditions := make([][2]interface{}, 0)
	if len(req.IDs) > 0 {
		var condition [2]interface{}
		condition[0] = "id in ?"
		condition[1] = req.IDs
		conditions = append(conditions, condition)
	}
	if req.Key != "" {
		var condition [2]interface{}
		condition[0] = "feature_key like ?"
		condition[1] = fmt.Sprintf("%%%s%%", req.Key)
		conditions = append(conditions, condition)
	}
	if req.Name != "" {
		var condition [2]interface{}
		condition[0] = "name like ?"
		condition[1] = fmt.Sprintf("%%%s%%", req.Name)
		conditions = append(conditions, condition)
	}
	if req.Valid != nil {
		var condition [2]interface{}
		condition[0] = "valid = ?"
		condition[1] = *req.Valid
		conditions = append(conditions, condition)
	}
	features, err := s.dao.GetFeatureList2(conditions)
	if err != nil {
		logger.Logger.Errorf("get feature list failed: %v", err)
		errorModel.Code = constant.Error
		errorModel.Reason = err.Error()
		return
	}
	return
}

func (s *FeatureService) CreateFeatureValue(ctx context.Context, featureValue *model.FeatureValue) (errorModel model.Error) {
	errorModel.Code = constant.OK
	if featureValue == nil || featureValue.FeatureID == nil {
		return model.Error{Code: constant.Error}
	}
	feature := &model.Feature{}
	feature.ID = *featureValue.FeatureID
	err := s.dao.GetFeature(feature)
	if err != nil {
		logger.Logger.Errorf("get feature failed,err=%v", err)
		errorModel.Code = constant.Error
		errorModel.Reason = err.Error()
		return
	}
	errorModel = ValidateFeatureValueCreate(feature, featureValue)
	if errorModel.Code != constant.OK {
		return
	}
	err = s.dao.CreateFeatureValue(ctx, featureValue)
	if err != nil {
		logger.Logger.Errorf("create feature value=%v failed: %v", featureValue, err)
		errorModel.Code = constant.Error
		errorModel.Reason = err.Error()
	}
	return
}
