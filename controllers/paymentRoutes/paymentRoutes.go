package paymentroutes

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/reminders/models"
	"gorm.io/gorm"
)

type paymentInfo struct {
	KeyId     string
	KeySecret string
	DB        *gorm.DB
}

type orderRequestSchema struct {
	Amount   int64
	Currency string
}

func setupPaymentInfo(DB *gorm.DB, apiKey, keySecret string) paymentInfo {
	return paymentInfo{
		KeyId:     apiKey,
		KeySecret: keySecret,
	}
}

func SetupPaymentRoutes(route *gin.RouterGroup, db *gorm.DB, apiKey, keySecret string) {
	paymentInfo := setupPaymentInfo(db, apiKey, keySecret)
	route.POST("/create-order", paymentInfo.createOrderHandler)
}

func (p *paymentInfo) createOrderHandler(ctx *gin.Context) {
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
	var data orderRequestSchema
	if err := ctx.ShouldBindBodyWithJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "request data not correct",
		})
		return
	}

	body, _ := json.Marshal(data)
	//get razorpay order id
	req, err := http.NewRequest("POST", "https://api.razorpay.com/v1/orders", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return
	}
	req.SetBasicAuth(p.KeyId, p.KeySecret)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error making api request %v:", err)
		return
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Fatalf("Error decoding response: %v", err)
		return
	}
	orderId := res["id"].(string)

	query := "INSERT INTO orders(created_at,updated_at,order_id,amount,currency,user_id) VALUES (?,?,?,?,?,?)"
	result := p.DB.Exec(query, time.Now(), time.Now(), orderId, data.Amount, data.Currency, user.ID)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": orderId,
	})
}
