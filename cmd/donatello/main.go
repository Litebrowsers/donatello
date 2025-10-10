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
	"fmt"
	"github.com/Litebrowsers/donatello/internal/canvas_tasks"
	"github.com/Litebrowsers/donatello/internal/models"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Litebrowsers/donatello/internal/cache"
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

	memoryCache := cache.NewInMemoryCache()

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

		combinedHash, err := canvas.CalculateCombinedHash(hashes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate hashes"})
			return
		}

		fmt.Println(combinedHash)
		id := uuid.New()

		firstTask := firstTaskGenerator.GenerateTask()

		challenge := cache.Challenge{
			Task:         firstTask,
			ExpectedHash: combinedHash,
			ExpiresAt:    time.Now().Add(1 * time.Minute),
		}

		err = memoryCache.Set(id.String(), challenge)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task to cache"})
			return
		}

		numShapesSecondTask := rand.Intn(6) + 1
		randomShapesSecondTask := canvas_tasks.GenerateRandomShapes(canvasSize, numShapesSecondTask)
		secondTaskGenerator := canvas_tasks.NewTaskGenerator(randomShapesSecondTask...)

		secondTask := secondTaskGenerator.GenerateTask()

		c.JSON(http.StatusOK, gin.H{
			"id":          id.String(),
			"first_task":  firstTask,
			"second_task": secondTask,
		})
	})
	router.POST("/challenge", func(c *gin.Context) {
		var answer models.ChallengeAnswer

		if err := c.ShouldBindJSON(&answer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid JSON: " + err.Error(),
			})
			return
		}

		fmt.Printf("Challenge ID: %s\n", answer.ID)
		fmt.Printf("TotalHash1: %s\n", answer.FirstTaskHash)
		fmt.Printf("TotalHash2: %s\n", answer.SecondTaskHash)

		challenge, exists, err := memoryCache.Get(answer.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task to cache"})
			return
		}

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
		}

		noiseDetect := challenge.ExpectedHash != answer.FirstTaskHash

		c.JSON(http.StatusOK, gin.H{
			"status":         "ok",
			"noise_detected": noiseDetect,
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
