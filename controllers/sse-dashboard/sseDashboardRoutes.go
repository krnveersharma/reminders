package SSEDashboardController

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reminders/controllers/sse-dashboard/clients"
)

var TotalReminders = 0

func SetUpSSEDashboardRoutes(router *gin.RouterGroup) {
	router.GET("/get-total-reminders", SSEHandler)
}

func SSEHandler(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	flusher, ok := ctx.Writer.(http.Flusher)

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "streaming unsupported",
		})
	}

	msgChan := make(chan int)
	clients.AddClient(msgChan)

	defer clients.RemoveClient(msgChan)

	notify := ctx.Request.Context().Done()

	for {
		select {

		case <-notify:
			return

		case msg := <-msgChan:
			fmt.Fprintf(ctx.Writer, "data: %d\n\n", msg)
			flusher.Flush()
		}
	}
}
