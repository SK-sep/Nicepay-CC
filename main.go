package main

import (
	module "eska/Module"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	MerchantId := os.Getenv("NICEPAY_MERCHANT_ID")
	MerchantKey := os.Getenv("NICEPAY_MERCHANT_KEY")
	register_endpoint := os.Getenv("REGISTER_ENDPOINT")
	payment_endpoint := os.Getenv("PAYMENT_ENDPOINT")
	inq_endpoint := os.Getenv("INQ_ENDPOINT")

	router := gin.Default()
	router.POST("/register", module.Register(MerchantId, MerchantKey, register_endpoint))
	router.POST("/payment", module.Payment(MerchantId, MerchantKey, payment_endpoint))
	router.GET("/check-status", module.Status(MerchantId, MerchantKey, inq_endpoint))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"Statuscode": "404", "message": "Path Not Found, Please Check Your URL!"})
	})
	router.Run("localhost:8080")
}
