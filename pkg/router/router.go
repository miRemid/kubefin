/*
Copyright 2022 The KubeFin Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/kubefin/kubefin/api"
	"github.com/kubefin/kubefin/pkg/server/costs_handler"
	"github.com/kubefin/kubefin/pkg/server/metrics_handler"
)

func NewServerRouter() *gin.Engine {
	router := gin.New()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"*"}
	corsHandler := cors.New(config)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Use(corsHandler)

	initMetricsRouter(router, corsHandler)
	initCostAnalyzeRouter(router, corsHandler)

	return router
}

func initMetricsRouter(router *gin.Engine, corsHandler gin.HandlerFunc) {
	metricsGroup := router.Group("/api/v1/metrics")
	metricsGroup.GET("/summary", metrics_handler.ClustersMetricsSummaryHandler)
	metricsGroup.GET("/clusters/:cluster_id/summary", metrics_handler.ClusterMetricsSummaryHandler)
	metricsGroup.GET("/clusters/:cluster_id/cpu", metrics_handler.ClusterCPUMetricsHandler)
	metricsGroup.GET("/clusters/:cluster_id/memory", metrics_handler.ClusterMemoryMetricsHandler)
	metricsGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	metricsGroup.Use(corsHandler)
}

func initCostAnalyzeRouter(router *gin.Engine, corsHandler gin.HandlerFunc) {
	costsGroup := router.Group("/api/v1/costs")
	costsGroup.GET("/summary", costs_handler.ClustersCostsSummaryHandler)
	costsGroup.GET("/clusters/:cluster_id/summary", costs_handler.ClusterCostsSummaryHandler)
	costsGroup.GET("/clusters/:cluster_id/resource", costs_handler.ClusterResourceCostsHandler)
	costsGroup.GET("/clusters/:cluster_id/workload", costs_handler.ClusterWorkloadsCostsHandler)
	costsGroup.GET("/clusters/:cluster_id/namespace", costs_handler.ClusterNamespacesCostsHandler)
	costsGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	costsGroup.Use(corsHandler)
}
