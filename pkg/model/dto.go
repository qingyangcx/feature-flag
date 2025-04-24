package model

import "feature-flag/pkg/constant"

type Error struct {
	Code   constant.Code `json:"code"`
	Reason string        `json:"reason"`
}

type CreateFeatureReq Feature
type CreateFeatureRsp struct {
	FeatureID *uint64 `json:"feature_id"`
	Error     Error   `json:"error"`
}

type ModifyFeatureReq Feature
type ModifyFeatureRsp struct {
	Error Error `json:"error"`
}

type ModifyFeatureValueReq FeatureValue

type ModifyFeatureValueRsp struct {
	Error Error `json:"error"`
}

type GetFeatureListReq struct {
	IDs       []uint64 `json:"ids"`
	Name      string   `json:"name"`
	Key       string   `json:"key"`
	Namespace string   `json:"namespace"`
	Valid     *bool    `json:"valid"`
}

type GetFeatureListRsp struct {
	Features []*Feature `json:"features"`
	Error    Error      `json:"error"`
}
type GetFeatureRsp struct {
	Feature *Feature `json:"feature"`
	Error   Error    `json:"error"`
}

type SplitReq struct {
	Namespace string  `json:"namespace"`
	Identity  string  `json:"identity"`
	FeatureID *uint64 `json:"feature_id"`
}

type SplitValue struct {
	FeatureID  uint64 `json:"feature_id"`
	ValueID    uint64 `json:"value_id"`
	Value      string `json:"value"`
	Reason     string `json:"reason"`
	ReasonType uint8  `json:"reason_type"`
}

func (v *SplitValue) SetValue(featureValue *FeatureValue) {
	v.FeatureID = *featureValue.FeatureID
	v.ValueID = featureValue.ID
	v.Value = featureValue.Value
}

type SplitRsp struct {
	Values []*SplitValue `json:"feature_values"`
	Error  Error         `json:"error"`
}

type CreateFeatureValueReq FeatureValue
type CreateFeatureValueRsp struct {
	FeatureValueID *uint64 `json:"feature_value_id"`
	Error          Error   `json:"error"`
}
