package dao

import (
	"context"
	"feature-flag/pkg/config"
	"feature-flag/pkg/constant"
	"feature-flag/pkg/logger"
	"feature-flag/pkg/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Dao struct {
	db *gorm.DB
}

func NewDao() *Dao {
	database := config.GlobalConfig.Database
	dsn := fmt.Sprintf(constant.DB_CONNECT_TEMPLATE, database.User, database.Password, database.Url, database.Db)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Logger.Fatalf("connect db failed,%v", err)
	}
	return &Dao{db: db}
}

func (dao *Dao) CreateFeature(ctx context.Context, feature *model.Feature) error {
	return dao.db.Create(feature).Error
}

func (dao *Dao) ModifyFeature(ctx context.Context, feature *model.Feature) error {
	result := dao.db.Select("blacklist").Updates(feature)
	if err := result.Error; err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no row affected for feature id=%d", feature.ID)
	}
	return nil
}

func (dao *Dao) ModifyFeatureValue(ctx context.Context, featureValue *model.FeatureValue, conditions map[string]interface{}) error {
	if err := dao.db.Model(featureValue).Where(map[string]interface{}{"id": featureValue.ID, "feature_id": featureValue.FeatureID}).Updates(conditions).Error; err != nil {
		return err
	}
	return nil
}

func (dao *Dao) GetFeature(feature *model.Feature) error {
	return dao.db.Preload("FeatureValues").First(feature).Error
}

func (dao *Dao) GetFeatureList1(conditions map[string]interface{}) (features []*model.Feature, err error) {
	features = make([]*model.Feature, 0)
	err = dao.db.Preload("FeatureValues").Where(conditions).Find(&features).Error
	return
}

func (dao *Dao) GetFeatureList2(conditions [][2]interface{}) (features []*model.Feature, err error) {
	features = make([]*model.Feature, 0)
	tx := dao.db
	for _, condition := range conditions {
		if len(condition) != 2 {
			logger.Logger.Warnf("invalid condition: %v", condition)
			continue
		}
		tx = tx.Where(condition[0], condition[1])
	}
	err = tx.Preload("FeatureValues").Find(&features).Error
	return
}

func (dao *Dao) CreateFeatureValue(ctx context.Context, featureValue *model.FeatureValue) error {
	return dao.db.Create(featureValue).Error
}
