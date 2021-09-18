package controller

import (
	"context"
	"dcs-sched/pkg/common"
	"dcs-sched/pkg/dto"
	"github.com/carter115/gslog"
	"github.com/gin-gonic/gin"
	"strconv"
)

func JobGroupRegistry(group *gin.RouterGroup) {
	job := JobController{}
	ctx := context.Background()
	if err := InitJobMgr(); err != nil {
		gslog.Errorf(ctx, "init job mgr error %s", err.Error())
	}

	group.POST("/save", job.Save)
	group.GET("/list", job.List)
	group.POST("/delete/:id", job.Delete)
	group.POST("/kill/:id", job.Kill)
}

type JobController struct{}

// Save JobController godoc
// @Summary 保存任务
// @Tags 任务
// @Accept json
// @Produce json
// @Param data body dto.JobInput true "data"
// @Success 200 {object} common.Response
// @Router /job/save [post]
func (j *JobController) Save(c *gin.Context) {
	input := dto.JobInput{}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, common.NewResponse(c, common.InvalidParam, nil))
		return
	}

	job := common.NewJob(input.Name, input.Command, input.CronExpr)
	// 参数包含Id，则修改已存在的Job
	if input.Id != 0 {
		job.Id = input.Id
	}

	if err := JobMgr.Save(c, job); err != nil {
		c.JSON(400, common.NewResponse(c, common.JobSaveError, job))
		return
	}
	c.JSON(200, common.SuccessResponse(c, nil))
}

// Delete JobController godoc
// @Summary 删除任务
// @Tags 任务
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} common.Response{data=common.Job}
// @Router /job/delete/{id} [post]
func (j *JobController) Delete(c *gin.Context) {
	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, common.NewResponse(c, common.InvalidParam, nil))
		return
	}

	oldJob, err := JobMgr.Delete(c, int64(aid))
	if err != nil {
		gslog.Error(c, err.Error())
		c.JSON(400, common.NewResponse(c, common.JobDeleteError, nil))
		return
	}
	c.JSON(200, common.SuccessResponse(c, oldJob))
}

// List JobController godoc
// @Summary 列出任务
// @Tags 任务
// @Accept json
// @Produce json
// @Success 200 {object} common.Response{data=[]common.Job}
// @Router /job/list [get]
func (j *JobController) List(c *gin.Context) {
	jobs, err := JobMgr.List(c)
	if err != nil {
		gslog.Error(c, err.Error())
		c.JSON(400, common.NewResponse(c, common.JobListError, nil))
		return
	}
	c.JSON(200, common.SuccessResponse(c, jobs))
}

// Kill JobController godoc
// @Summary 杀死任务
// @Tags 任务
// @Accept json
// @Produce json
// @Param id path int true "任务ID"
// @Success 200 {object} common.Response
// @Router /job/kill/{id} [post]
func (j *JobController) Kill(c *gin.Context) {
	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, common.NewResponse(c, common.InvalidParam, nil))
		return
	}
	if err = JobMgr.Kill(c, int64(aid)); err != nil {
		gslog.Error(c, err.Error())
		c.JSON(400, common.NewResponse(c, common.JobKillError, nil))
		return
	}
	c.JSON(200, common.SuccessResponse(c, nil))
}
