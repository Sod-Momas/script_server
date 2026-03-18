package handlers

import (
	"datasources/models"
	"datasources/services"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// TestRequest 连通性测试请求结构
type TestRequest struct {
	DBType   string `json:"db_type" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// SaveRequest 保存数据源请求结构
type SaveRequest struct {
	Name     string `json:"name" binding:"required"`
	DBType   string `json:"db_type" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Index 首页处理器
func Index(c *gin.Context) {
	c.File("./static/index.html")
}

// TestConnectivity 测试连通性
func TestConnectivity(c *gin.Context) {
	var req TestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误: " + err.Error()})
		return
	}

	// 构建服务请求
	serviceReq := services.TestRequest{
		DBType:   req.DBType,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
	}

	// 执行测试
	result := services.TestConnectivity(serviceReq)

	c.JSON(http.StatusOK, gin.H{
		"success": result.Success,
		"message": result.Message,
	})
}

// SaveDatasource 保存数据源
func SaveDatasource(c *gin.Context) {
	var req SaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误: " + err.Error()})
		return
	}

	// 创建数据源模型
	ds := models.Datasource{
		Name:     req.Name,
		DBType:   req.DBType,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
	}

	// 保存到数据库
	if err := ds.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据源保存成功!",
		"data": gin.H{
			"id": ds.ID,
		},
	})
}

// GetDatasourceList 获取数据源列表
func GetDatasourceList(c *gin.Context) {
	datasources, err := models.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    datasources,
	})
}

// GetDatasourceDetail 获取数据源详情
func GetDatasourceDetail(c *gin.Context) {
	idStr := c.Param("id")
	var id int64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的数据源ID"})
		return
	}

	datasource, err := models.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    datasource,
	})
}

// GetStaticPath 获取静态文件目录
func GetStaticPath() string {
	return filepath.Join(".", "static")
}
