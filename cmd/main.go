package main

import (
	"feature-flag/pkg/config"
	"feature-flag/pkg/handler"
	"feature-flag/pkg/logger"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	workDir := flag.String("work_dir", "", "path of work directory")
	configPath := flag.String("config_file", "", "path of config file")
	flag.Parse()
	err := config.LoadConfig(*workDir, *configPath)
	if err != nil {
		log.Fatal(err)
	}

	handler := handler.NewHandler()
	if handler == nil {
		logger.Logger.Fatal("handler init failed")
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/feature", handler.HandleCreateFeature)
	r.GET("/features", handler.HandleGetFeatureList)
	r.GET("/feature/:id", handler.HandleGetFeature)
	r.GET("/split", handler.HandleSplit)
	r.PATCH("/feature/:id", handler.HandleModifyFeature)
	r.PATCH("/feature_value/:id", handler.HandleModifyFeatureValue)
	r.POST("/feature_value", handler.HandleCreateFeatureValue)
	r.Run(fmt.Sprintf("0.0.0.0:%d", config.GlobalConfig.FeatureFlag.Port))
}
