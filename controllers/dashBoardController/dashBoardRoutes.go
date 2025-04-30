package dashboardcontroller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reminders/controllers/sse-dashboard/clients"
	"gorm.io/gorm"
)

var (
	TotalReminders = 0
)

type dashBoardRoutes struct {
	DB     *gorm.DB
	Secret string
}

func SetUpDashBoardRoutes(db *gorm.DB, secret string, router *gin.RouterGroup) {
	dashboarRoute := dashBoardRoutes{
		DB:     db,
		Secret: secret,
	}

	router.GET("/reminder-count", dashboarRoute.getReminders)
}

func (d *dashBoardRoutes) getReminders(ctx *gin.Context) {

	var result int64
	query := "SELECT COUNT(*) from reminders"
	if err := d.DB.Raw(query).Scan(&result).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to Fetch count",
		})
		return
	}
	TotalReminders = int(result)
	clients.BroadCastMessage(TotalReminders)
	fmt.Printf("broadcasted message \n")
	ctx.JSON(http.StatusOK, gin.H{
		"message": result,
	})
}
