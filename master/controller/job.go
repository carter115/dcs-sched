package controller

import (
	"context"
	"dcs-sched/pkg/common"
	"dcs-sched/pkg/dto"
	"github.com/carter115/gslog"
	"github.com/gin-gonic/gin"
)

func JobGroupRegistry(group *gin.RouterGroup) {
	job := JobController{}
	ctx := context.Background()
	if err := InitJobMgr(); err != nil {
		gslog.Errorf(ctx, "init job mgr error %s", err.Error())
	}

	group.POST("/save", job.Save)
	group.GET("/list", job.List)
	group.POST("/delete", job.Delete)
	group.POST("/kill", job.Kill)
}

type JobController struct{}

// Save JobController godoc
// @Summary 保存任务
// @Tags 任务
// @Accept json
// @Produce json
// @Param data body dto.JobInput true "data"
// @Success 200 {object} common.Response{data=dto.JobOutput} "desc
// @Router /job/save [post]
func (j *JobController) Save(c *gin.Context) {
	input := dto.JobInput{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, common.NewResponse(common.InvalidParam, nil))
		return
	}
	job := common.NewJob(input.Name, input.Command, input.CronExpr)
	JobMgr.Save(c, job)
	c.JSON(200, common.SuccessResponse(nil))
}

// List JobController godoc
// @Summary 列出任务
// @Tags 任务
// @Accept json
// @Produce json
// @Success 200 {string} string ""
// @Router /job/list [get]
func (j *JobController) List(c *gin.Context) {

}

// Delete JobController godoc
// @Summary 删除任务
// @Tags 任务
// @Accept json
// @Produce json
// @Success 200 {string} string ""
// @Router /job/delete [post]
func (j *JobController) Delete(c *gin.Context) {

}

// Kill JobController godoc
// @Summary 杀死任务
// @Tags 任务
// @Accept json
// @Produce json
// @Success 200 {string} string ""
// @Router /job/kill [post]
func (j *JobController) Kill(c *gin.Context) {

}
