/*
# Donatello

Copyright © 2025 Litebrowsers
Licensed under a Proprietary License

This software is the confidential and proprietary information of Litebrowsers
Unauthorized copying, redistribution, or use is prohibited.
For licensing inquiries, contact:
vera cohopie at gmail dot com
thor betson at gmail dot com
*/

package main

import (
	"github.com/Litebrowsers/donatello/canvas_tasks"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitMiddleware returns a gin.HandlerFunc that limits requests.
func RateLimitMiddleware(limit rate.Limit, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(limit, burst)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}

func main() {
	router := gin.Default()

	// Configure port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	canvasSize := 20
	canvasSizeStr := os.Getenv("CANVAS_SIZE")
	if canvasSizeStr != "" {
		canvasSize, _ = strconv.Atoi(canvasSizeStr)
	}

	// Apply Rate Limiter Middleware
	router.Use(RateLimitMiddleware(rate.Every(time.Second/5), 10))

	router.GET("/challenge", func(c *gin.Context) {
		numShapesFirstTask := rand.Intn(6) + 1
		randomShapesFirstTask := canvas_tasks.GenerateRandomEvenSizedPrimitives(canvasSize, numShapesFirstTask)
		firstTaskGenerator := canvas_tasks.NewTaskGenerator(randomShapesFirstTask...)

		// Server-side drawing
		canvas := canvas_tasks.NewCanvas(canvasSize, canvasSize)
		canvas.DrawShapes(randomShapesFirstTask)
		hashes, err := canvas.CalculateHashes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate hashes"})
			return
		}

		firstTask := firstTaskGenerator.GenerateTask()

		numShapesSecondTask := rand.Intn(6) + 1
		randomShapesSecondTask := canvas_tasks.GenerateRandomShapes(canvasSize, numShapesSecondTask)
		secondTaskGenerator := canvas_tasks.NewTaskGenerator(randomShapesSecondTask...)

		secondTask := secondTaskGenerator.GenerateTask()

		c.JSON(http.StatusOK, gin.H{
			"first_task":  firstTask,
			"second_task": secondTask,
			"hashes":      hashes,
		})
	})
	router.GET("/", func(c *gin.Context) {
		root, err := os.Getwd()
		if err != nil {
			c.String(500, "Failed to get working directory")
			return
		}

		filePath := filepath.Join(root, "resources", "index.html")

		c.File(filePath)
	})

	log.Printf("Server starting on port %s", port)
	err := router.Run(":" + port)
	if err != nil {
		log.Printf("Server can't be started %s", err.Error())
		return
	}
}
