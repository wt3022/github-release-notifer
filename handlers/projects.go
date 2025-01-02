package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wt3022/github-release-notifier/internal/db"
	"github.com/wt3022/github-release-notifier/internal/utils"
	"gorm.io/gorm"
)

func ListProjects(c *gin.Context, dbClient *gorm.DB) {
	/*
		プロジェクト一覧を取得します
		クエリパラメータ:
			name: プロジェクト名の部分一致
			created_at__gte: 特定の作成日より前
			created_at__lte: 特定の作成日より後
			updated_at__gte: 特定の更新日より前
			updated_at__lte: 特定の更新日より後
			page: ページ番号
			page_size: ページあたりのアイテム数
	*/
	var projects []db.Project

	name := c.Query("name")

	query := utils.BuildQuery(c, dbClient)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if err := query.Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func DetailProject(c *gin.Context, dbClient *gorm.DB) {
	/*
		プロジェクトの詳細を取得します
		その際に通知情報も展開します
	*/
	var project db.Project
	var notification db.Notification

	id := c.Param("id")
	if err := dbClient.First(&project, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := dbClient.First(&notification, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	project.Notification = &notification

	c.JSON(http.StatusOK, project)
}

func CreateProjects(c *gin.Context, dbClient *gorm.DB) {
	/* プロジェクトを作成します */
	var projectRequest db.Project

	if err := c.ShouldBindJSON(&projectRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := dbClient.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tx.Error.Error()})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 通知設定の作成後にプロジェクトを作成
	if err := tx.Create(&projectRequest.Notification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// プロジェクトを作成
	if err := tx.Create(&projectRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, projectRequest)
}

func UpdateProject(c *gin.Context, dbClient *gorm.DB) {
	var project db.Project

	// リクエストボディの内容を取得
	var updateData db.Project
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dbClient.Preload("Notification").First(&project, updateData.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "プロジェクトが見つかりません"})
		return
	}

	tx := dbClient.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// プロジェクトの更新
	if err := tx.Model(&project).Updates(map[string]interface{}{
		"name":        updateData.Name,
		"description": updateData.Description,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 通知設定の更新（既存の通知設定のIDを使用）
	if err := tx.Model(&project.Notification).Updates(map[string]interface{}{
		"type": updateData.Notification.Type,
	}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// トランザクション内で更新後のデータを取得
	var updatedProject db.Project
	if err := tx.Preload("Notification").First(&updatedProject, updateData.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProject)
}

func DeleteProject(c *gin.Context, dbClient *gorm.DB) {
	/* プロジェクトを削除します */
	var project db.Project

	id := c.Param("id")
	if err := dbClient.Delete(&project, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func BulkDeleteProjects(c *gin.Context, dbClient *gorm.DB) {
	/* プロジェクトを一括削除します */
	var projectIds []int

	if err := c.ShouldBindJSON(&projectIds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := dbClient.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	if err := tx.Delete(&db.Project{}, projectIds).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
