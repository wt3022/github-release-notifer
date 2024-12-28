package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wt3022/github-release-notifier/internal/db"
	"gorm.io/gorm"
)

func ListNotifications(ctx *gin.Context, dbClient *gorm.DB) {
	var notifications []db.Notification

	if err := dbClient.Find(&notifications).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}

func DetailNotification(ctx *gin.Context, dbClient *gorm.DB) {
	var notification db.Notification

	id := ctx.Param("id")
	if err := dbClient.First(&notification, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notification)
}

func CreateNotification(ctx *gin.Context, dbClient *gorm.DB) {
	var notification db.Notification

	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// プロジェクトが存在するか確認
	var project db.Project
	if err := dbClient.First(&project, notification.ProjectID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "プロジェクトが存在しません"})
		return
	}

	if err := dbClient.Create(&notification).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, notification)
}

func UpdateNotification(ctx *gin.Context, dbClient *gorm.DB) {
	var notification db.Notification

	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dbClient.Save(&notification).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notification)
}

func DeleteNotification(ctx *gin.Context, dbClient *gorm.DB) {
	var notification db.Notification

	id := ctx.Param("id")
	if err := dbClient.Delete(&notification, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notification)
}
