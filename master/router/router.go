package router

import (
	"dcs-sched/master/config"
	"dcs-sched/master/controller"
	_ "dcs-sched/master/docs"
	"dcs-sched/pkg/common"
	"dcs-sched/pkg/middleware"
	"github.com/carter115/gslog"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func Server() error {
	r := gin.New()
	r.Use(gslog.GinLogger(), middleware.Recovery())

	// router
	r.GET("/", homeHandler)
	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

	// admin
	adminGroup := r.Group("/admin")
	controller.AdminGroupRegistry(adminGroup)

	// job
	jobGroup := r.Group("/job")
	controller.JobGroupRegistry(jobGroup)

	// run
	s := http.Server{
		Addr:           config.Config.Server.Addr,
		Handler:        r,
		ReadTimeout:    config.Config.Server.ReadTimeout,
		WriteTimeout:   config.Config.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	return s.ListenAndServe()
}

// homeHandler godoc
// @Summary 首页
// @Tags 首页
// @Accept json
// @Produce json
// @Success 200 {string} string ""
// @Router / [get]
func homeHandler(c *gin.Context) {
	c.JSON(200, common.SuccessResponse(c, nil))
}
