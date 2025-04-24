package model

import "time"

type Feature struct {
	ID            uint64          `json:"id"`
	Name          string          `json:"name"`
	Key           string          `json:"key" gorm:"column:feature_key"`
	Blacklist     []string        `json:"blacklist" gorm:"column:blacklist;serializer:json"`
	Valid         *bool           `json:"valid"`
	Namespace     string          `json:"namespace"`
	FeatureValues []*FeatureValue `json:"feature_values" gorm:"foreignKey:FeatureID"`
	CreateTime    time.Time       `json:"create_time" gorm:"autoCreateTime"`
	ModifyTime    time.Time       `json:"modify_time" gorm:"autoUpdateTime"`

	BlacklistSet map[string]struct{} `json:"-" gorm:"-"`
	DefaultValue *FeatureValue       `json:"-" gorm:"-"`
	Seed         uint32              `json:"-" gorm:"-"`
}

func (Feature) TableName() string {
	return "feature"
}

type FeatureValue struct {
	ID         uint64    `json:"id"`
	FeatureID  *uint64   `json:"feature_id"`
	Value      string    `json:"value" gorm:"column:feature_value"`
	Traffic    *uint32   `json:"traffic"`
	Whitelist  []string  `json:"whitelist" gorm:"column:whitelist;serializer:json"`
	Default    *bool     `json:"default" gorm:"column:default"`
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	ModifyTime time.Time `json:"modify_time" gorm:"autoUpdateTime"`

	WhitelistSet map[string]struct{} `json:"-" gorm:"-"`
}

func (FeatureValue) TableName() string {
	return "feature_value"
}
