package handler

import (
	"context"
	"feature-flag/pkg/constant"
	"feature-flag/pkg/logger"
	"feature-flag/pkg/model"
	"feature-flag/pkg/service/cache"
	"feature-flag/pkg/service/feature_flag"
	"feature-flag/pkg/service/split"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const kTimeOutMilliSec int64 = 500

type Handler struct {
	featureFlagService *feature_flag.FeatureService
	splitService       *split.SplitService
	cache              *cache.Cache
}

func NewHandler() *Handler {
	return &Handler{featureFlagService: feature_flag.NewFeatureService(), splitService: split.NewSplitService(), cache: cache.NewCache()}
}

func (h *Handler) HandleCreateFeature(c *gin.Context) {
	req := model.CreateFeatureReq{}
	rsp := model.CreateFeatureRsp{}
	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Error.Code = constant.Error
		rsp.Error.Reason = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(kTimeOutMilliSec))
	defer cancel()
	logger.Logger.Infof("create feature=%+v", req)
	rsp.FeatureID, rsp.Error = h.featureFlagService.CreateFeatureFlag(ctx, (*model.Feature)(&req))
	if rsp.Error.Code != constant.OK {
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	c.JSON(http.StatusOK, rsp)
}

func (h *Handler) HandleModifyFeature(c *gin.Context) {
	req := model.ModifyFeatureReq{}
	rsp := model.ModifyFeatureRsp{}
	err1 := c.ShouldBindJSON(&req)
	id, err2 := strconv.ParseUint(c.Param("id"), 10, 64)

	if err1 != nil || err2 != nil {
		rsp.Error.Code = constant.Error
		if err1 != nil {
			rsp.Error.Reason = err1.Error()
		}
		if err2 != nil {
			rsp.Error.Reason += err2.Error()
		}
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	req.ID = id
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(kTimeOutMilliSec))
	defer cancel()
	logger.Logger.Infof("modify feature=%+v", req)
	rsp.Error = h.featureFlagService.ModifyFeatureFlag(ctx, (*model.Feature)(&req))
	if rsp.Error.Code != constant.OK {
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	c.JSON(http.StatusOK, rsp)
}

func (h *Handler) HandleGetFeatureList(c *gin.Context) {
	req := model.GetFeatureListReq{}
	rsp := model.GetFeatureListRsp{}
	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Error.Code = constant.Error
		rsp.Error.Reason = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(kTimeOutMilliSec))
	defer cancel()
	logger.Logger.Infof("get feature list, req=%+v", req)
	rsp.Features, rsp.Error = h.featureFlagService.GetFeatureList(ctx, &req)
	if rsp.Error.Code != constant.OK {
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	c.JSON(http.StatusOK, rsp)
}

func (h *Handler) HandleGetFeature(c *gin.Context) {
	rsp := model.GetFeatureRsp{}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		rsp.Error.Code = constant.Error
		rsp.Error.Reason = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(kTimeOutMilliSec))
	defer cancel()
	logger.Logger.Infof("get feature, id=%d", id)
	rsp.Feature, rsp.Error = h.featureFlagService.GetFeature(ctx, id)
	if rsp.Error.Code != constant.OK {
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	c.JSON(http.StatusOK, rsp)
}

func (h *Handler) HandleModifyFeatureValue(c *gin.Context) {
	req := model.ModifyFeatureValueReq{}
	rsp := model.ModifyFeatureValueRsp{}
	err1 := c.ShouldBindJSON(&req)
	valueID, err2 := strconv.ParseUint(c.Param("id"), 10, 64)

	if err1 != nil || err2 != nil {
		rsp.Error.Code = constant.Error
		if err1 != nil {
			rsp.Error.Reason = err1.Error()
		}
		if err2 != nil {
			rsp.Error.Reason += err2.Error()
		}
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	if req.FeatureID == nil {
		rsp.Error.Code = constant.Error
		rsp.Error.Reason = "lack of feature_id"
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	req.ID = valueID
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(kTimeOutMilliSec))
	defer cancel()
	logger.Logger.Infof("modify feature value=%+v", req)
	rsp.Error = h.featureFlagService.ModifyFeatureValue(ctx, (*model.FeatureValue)(&req))
	if rsp.Error.Code != constant.OK {
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	c.JSON(http.StatusOK, rsp)
}

func (h *Handler) HandleSplit(c *gin.Context) {
	req := model.SplitReq{}
	rsp := model.SplitRsp{Error: model.Error{Code: constant.OK}}
	defer func() {
		if len(rsp.Values) > 0 {
			logger.Logger.Infof("split finish, value[0]=%+v,err=%+v", rsp.Values[0], rsp.Error)
		} else {
			logger.Logger.Warnf("split finish,err=%+v", rsp.Error)
		}
	}()
	err := c.ShouldBindJSON(&req)
	if err != nil {
		rsp.Error.Code = constant.Error
		rsp.Error.Reason = err.Error()
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	logger.Logger.Infof("handle split,identity=%s, feature_id=%v, namespace=%s", req.Identity, func() interface{} {
		if req.FeatureID == nil {
			return nil
		}
		return *req.FeatureID
	}(), req.Namespace)
	if req.Namespace == "" {
		rsp.Error.Code = constant.Error
		rsp.Error.Reason = "namespace empty"
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	if req.Identity == "" {
		message := "identity empty"
		logger.Logger.Error(message)
		rsp.Error.Code = constant.Error
		rsp.Error.Reason = message
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	featureMap, ok := h.cache.Get()[req.Namespace]
	if !ok {
		message := fmt.Sprintf("not found features for namespace %s", req.Namespace)
		logger.Logger.Error(message)
		rsp.Error.Code = constant.Error
		rsp.Error.Reason = message
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	if req.FeatureID == nil {
		rsp.Values = h.splitService.SplitAll(req.Identity, featureMap)
	} else {
		feature, ok := featureMap[*req.FeatureID]
		if !ok {
			message := fmt.Sprintf("not found feature for id=%d", req.FeatureID)
			logger.Logger.Error(message)
			rsp.Error.Code = constant.Error
			rsp.Error.Reason = message
			c.JSON(http.StatusBadRequest, rsp)
			return
		}
		splitValue := h.splitService.SplitOne(req.Identity, feature)
		if splitValue != nil {
			rsp.Values = append(rsp.Values, splitValue)
		}
	}
	c.JSON(http.StatusOK, rsp)
}

func (h *Handler) HandleCreateFeatureValue(c *gin.Context) {
	req := model.CreateFeatureValueReq{}
	rsp := model.CreateFeatureValueRsp{}
	if err := c.ShouldBindJSON(&req); err != nil {
		rsp.Error.Code = constant.Error
		rsp.Error.Reason = err.Error()
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(kTimeOutMilliSec))
	defer cancel()
	logger.Logger.Infof("create feature value=%+v", req)
	rsp.Error = h.featureFlagService.CreateFeatureValue(ctx, (*model.FeatureValue)(&req))
	if rsp.Error.Code != constant.OK {
		c.JSON(http.StatusBadRequest, rsp)
		return
	}
	c.JSON(http.StatusOK, rsp)
}
