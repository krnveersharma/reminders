package dashboardcontroller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reminders/controllers/sse-dashboard/clients"
	"github.com/reminders/internal/dto"
	"github.com/reminders/middlewares"
	"github.com/reminders/models"
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

	setupMiddleware := middlewares.SetUpMiddleware(db, secret)

	router.GET("/reminder-count", dashboarRoute.getReminders)
	router.GET("/reminder-seven-days", dashboarRoute.getSevenDaysReminders)
	router.GET("/reminder-last-year", dashboarRoute.getLastYearMonthlyReminders)
	router.GET("/reminder-user-seven-days", setupMiddleware.UserAuth, dashboarRoute.getUserSevenDaysReminders)
	router.GET("/reminder-user-last-year", setupMiddleware.UserAuth, dashboarRoute.getUserLastYearMonthlyReminders)
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

func (d *dashBoardRoutes) getSevenDaysReminders(ctx *gin.Context) {

	var result []dto.DayCount

	query := "SELECT TO_CHAR(date, 'YYYY-MM-DD') AS day, COUNT(*) AS count FROM reminders WHERE date >= CURRENT_DATE - INTERVAL '6 days'  GROUP BY day ORDER BY day;"
	if err := d.DB.Raw(query).Scan(&result).Error; err != nil {
		fmt.Printf("Error is:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get data for 7 days",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": result,
	})
}

func (d *dashBoardRoutes) getLastYearMonthlyReminders(ctx *gin.Context) {

	var result []dto.MonthCount

	query := "SELECT TO_CHAR(date, 'YYYY-MM') AS month, COUNT(*) AS count FROM reminders WHERE date >= CURRENT_DATE - INTERVAL '1 year' GROUP BY TO_CHAR(date, 'YYYY-MM') ORDER BY month;"

	if err := d.DB.Raw(query).Scan(&result).Error; err != nil {
		fmt.Printf("Error is:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get data for 7 days",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": result,
	})
}

func (d *dashBoardRoutes) getUserSevenDaysReminders(ctx *gin.Context) {
	userVal, found := ctx.Get("user")
	if !found {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Please login first",
		})
		return
	}

	user, ok := userVal.(models.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user session",
		})
		return
	}

	var result []dto.DayCount

	query := "SELECT TO_CHAR(date, 'YYYY-MM-DD') AS day, COUNT(*) AS count FROM reminders WHERE date >= CURRENT_DATE - INTERVAL '6 days' and user_id = ?  GROUP BY day ORDER BY day;"
	if err := d.DB.Raw(query, user.ID).Scan(&result).Error; err != nil {
		fmt.Printf("Error is:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get data for 7 days",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": result,
	})
}

func (d *dashBoardRoutes) getUserLastYearMonthlyReminders(ctx *gin.Context) {
	userVal, found := ctx.Get("user")
	if !found {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Please login first",
		})
		return
	}

	user, ok := userVal.(models.User)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user session",
		})
		return
	}

	var result []dto.MonthCount

	query := "SELECT TO_CHAR(date, 'YYYY-MM') AS month, COUNT(*) AS count FROM reminders WHERE date >= CURRENT_DATE - INTERVAL '1 year' and user_id = ? GROUP BY TO_CHAR(date, 'YYYY-MM') ORDER BY month;"

	if err := d.DB.Raw(query, user.ID).Scan(&result).Error; err != nil {
		fmt.Printf("Error is:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get data for 7 days",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": result,
	})
}
