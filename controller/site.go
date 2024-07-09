package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
	"ops-api/dao"
	"ops-api/service"
	"ops-api/utils"
	"strconv"
)

var Site site

type site struct{}

// GetSiteList 获取站点列表
// @Summary 获取站点列表
// @Description 站点关接口
// @Tags 站点管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int true "分页"
// @Param limit query int true "分页大小"
// @Param name query string false "站点名称"
// @Success 200 {string} json "{"code": 0, "data": []}"
// @Router /api/v1/sites [get]
func (s *site) GetSiteList(c *gin.Context) {
	params := new(struct {
		Name  string `form:"name"`
		Page  int    `form:"page" binding:"required"`
		Limit int    `form:"limit" binding:"required"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	data, err := service.Site.GetSiteList(params.Name, params.Page, params.Limit)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": data,
	})
}

// AddGroup 创建分组
// @Summary 创建分组
// @Description 站点关接口
// @Tags 站点管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param group body service.SiteGroupCreate true "分组信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/site/group [post]
func (s *site) AddGroup(c *gin.Context) {
	var group = &service.SiteGroupCreate{}

	if err := c.ShouldBind(group); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	if err := service.Site.AddGroup(group); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": nil,
	})
}

// AddSite 创建站点
// @Summary 创建站点
// @Description 站点关接口
// @Tags 站点管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param group body service.SiteCreate true "分组信息"
// @Success 200 {string} json "{"code": 0, "msg": "创建成功", "data": nil}"
// @Router /api/v1/site [post]
func (s *site) AddSite(c *gin.Context) {
	var group = &service.SiteCreate{}

	if err := c.ShouldBind(group); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	if err := service.Site.AddSite(group); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": nil,
	})
}

// DeleteGroup 删除站点分组
// @Summary 删除站点分组
// @Description 站点关接口
// @Tags 站点管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "分组ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功", "data": nil}"
// @Router /api/v1/site/group/{id} [delete]
func (s *site) DeleteGroup(c *gin.Context) {

	// 对ID进行类型转换
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("ERROR：", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 执行删除
	if err := service.Site.DeleteGroup(groupID); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "删除成功",
		"data": nil,
	})
}

// DeleteSite 删除站点
// @Summary 删除站点
// @Description 站点关接口
// @Tags 站点管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "分组ID"
// @Success 200 {string} json "{"code": 0, "msg": "删除成功", "data": nil}"
// @Router /api/v1/site/{id} [delete]
func (s *site) DeleteSite(c *gin.Context) {

	// 对ID进行类型转换
	siteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("ERROR：", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 执行删除
	if err := service.Site.DeleteSite(siteID); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "删除成功",
		"data": nil,
	})
}

// UpdateGroup 更新分组信息
// @Summary 更新分组信息
// @Description 站点关接口
// @Tags 站点管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param group body service.SiteGroupUpdate true "分组信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/site/group [put]
func (s *site) UpdateGroup(c *gin.Context) {
	var data = &service.SiteGroupUpdate{}

	// 解析请求参数
	if err := c.ShouldBind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 更新站点分组信息
	if err := service.Site.UpdateGroup(data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "更新成功",
		"data": nil,
	})
}

// UpdateSite 更新站点信息
// @Summary 更新站点信息
// @Description 站点关接口
// @Tags 站点管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param group body dao.UpdateSite true "分组信息"
// @Success 200 {string} json "{"code": 0, "msg": "更新成功", "data": nil}"
// @Router /api/v1/site [put]
func (s *site) UpdateSite(c *gin.Context) {
	var data = &dao.UpdateSite{}

	// 解析请求参数
	if err := c.ShouldBind(&data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 更新站点信息
	if err := service.Site.UpdateSite(data); err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "更新成功",
		"data": nil,
	})
}

// UploadLogo 站点图片上传
// @Summary 站点图片上传
// @Description 站点关接口
// @Tags 站点管理
// @Param Authorization header string true "Bearer 用户令牌"
// @Param logo formData file true "头像"
// @Success 200 {string} json "{"code": 0, "path": logoPath}"
// @Router /api/v1/site/logoUpload [post]
func (s *site) UploadLogo(c *gin.Context) {
	// 获取上传的Logo
	logo, err := c.FormFile("icon")
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
		})
		return
	}

	// 打开上传的图片
	src, err := logo.Open()
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	// 上传图片到MinIO
	// 拼接存储的路径（此路径为临时路径，在表单提交时会将图片移动到实际位置）
	logoPath := fmt.Sprintf("site/logo/%v", logo.Filename)

	// 检查对象是否存在，err不为空是则表示对象已存在
	_, err = utils.StatObject(logoPath)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  "上传的对象已存在",
		})
		return
	}

	err = utils.FileUpload(logoPath, logo.Header.Get("Content-Type"), src, logo.Size)
	if err != nil {
		logger.Error("ERROR：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"path": logoPath,
		"msg":  "图片上传成功",
	})
}
