package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/buy_candy", buy_candy)
	r.Run("127.0.0.1:3333")
}

func buy_candy(c *gin.Context) {
	var candy_order CandyOrder
	c.ShouldBindJSON(&candy_order);


	if !is_valid(candy_order.CandyType) || candy_order.CandyCount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candyType or candyCount"})
		return
	}

	total_price := total_price(candy_order.CandyType, candy_order.CandyCount)

	if candy_order.Money < total_price {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "You need more money!"})
		return
	}

	change := candy_order.Money - total_price

	c.JSON(http.StatusCreated, gin.H{"thanks": "Thank you!", "change": change})
}

func is_valid(candyType string) bool {
	CandyTypes := []string{"CE", "AA", "NT", "DE", "YR"}
	for _, validType := range CandyTypes {
		if candyType == validType {
			return true
		}
	}
	return false
}

func total_price(candyType string, candyCount int) int {
	prices := map[string]int{
		"CE": 10,
		"AA": 15,
		"NT": 17,
		"DE": 21,
		"YR": 23,
	}
	return prices[candyType] * candyCount
}

type CandyOrder struct {
	Money     int    `json:"money"`
	CandyType string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}
