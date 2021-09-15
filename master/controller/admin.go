package controller

import "github.com/gin-gonic/gin"

type AdminController struct{}

func AdminGroupRegistry(group *gin.RouterGroup) {
	admin := AdminController{}
	group.GET("/login", admin.Login)
}

// Login AdminController godoc
// @Summary 后台登录
// @Tags 管理后台
// @Accept json
// @Produce json
// @Success 200 {string} string ""
// @Router /admin/login [post]
func (a *AdminController) Login(c *gin.Context) {

}
