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

// Package main is the entry point of the Donatello application.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Litebrowsers/donatello/internal/db"
	"github.com/Litebrowsers/donatello/internal/models"
	"github.com/Litebrowsers/donatello/internal/tasks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

var challengeExpiration time.Duration

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

func cleanupExpiredChallenges() {
	interval := challengeExpiration * 2
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Cleaning up expired challenges...")
		// Update expired challenges that are still pending (JavaScript is NULL)
		result := db.DB.Model(&models.Challenge{}).
			Where("expires_at < ? AND java_script IS NULL", time.Now()).
			Update("java_script", false)

		if result.Error != nil {
			log.Printf("Error cleaning up expired challenges: %v", result.Error)
			continue
		}

		if result.RowsAffected > 0 {
			log.Printf("Marked %d challenges as no-js.", result.RowsAffected)
		}
	}
}

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	err = db.DB.AutoMigrate(&models.Task{}, &models.Challenge{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Get challenge expiration from environment variable or use default
	expirationStr := os.Getenv("CHALLENGE_EXPIRATION")
	challengeExpiration = 500 * time.Millisecond // Default expiration

	if expirationStr != "" {
		parsedExpiration, err := time.ParseDuration(expirationStr)
		if err == nil {
			challengeExpiration = parsedExpiration
		} else {
			log.Printf("Invalid CHALLENGE_EXPIRATION format: %s. Using default 1m.", expirationStr)
		}
	}

	go cleanupExpiredChallenges()

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
		id := c.Query("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter is required"})
			return
		}

		var challenge models.Challenge
		result := db.DB.First(&challenge, "id = ?", id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
			return
		}

		numShapesFirstTask := rand.Intn(6) + 1
		randomShapesFirstTask := tasks.GenerateRandomEvenSizedPrimitives(canvasSize, numShapesFirstTask)
		firstTaskGenerator := tasks.NewTaskGenerator(randomShapesFirstTask...)

		// Server-side drawing
		canvas := tasks.NewCanvas(canvasSize, canvasSize)
		err = canvas.DrawShapes(randomShapesFirstTask)
		if err != nil {
			return
		}
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

		firstTask := firstTaskGenerator.GenerateTask()

		var secondTask models.Task
		result = db.DB.Where("name = ?", "secondTask").First(&secondTask)
		if result.Error != nil {
			numShapesSecondTask := rand.Intn(6) + 1
			randomShapesSecondTask := tasks.GenerateRandomShapes(canvasSize, numShapesSecondTask)
			secondTaskGenerator := tasks.NewTaskGenerator(randomShapesSecondTask...)
			secondTask.Value = secondTaskGenerator.GenerateTask()
			secondTask.Name = "secondTask"
			db.DB.Create(&secondTask)
		}

		secondTaskShapes, err := tasks.ParseTask(secondTask.Value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse second task"})
			return
		}

		secondTaskCanvas := tasks.NewCanvas(canvasSize, canvasSize)
		err = secondTaskCanvas.DrawShapes(secondTaskShapes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse second task"})
			return
		}
		secondTaskHashes, err := secondTaskCanvas.CalculateHashes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate hashes for second task"})
			return
		}

		secondTaskCombinedHash, err := secondTaskCanvas.CalculateCombinedHash(secondTaskHashes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate combined hash for second task"})
			return
		}

		challenge.Task = firstTask
		challenge.ExpectedHash = combinedHash
		challenge.Fingerprint = secondTaskCombinedHash

		result = db.DB.Save(&challenge)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task to cache"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":          id,
			"first_task":  firstTask,
			"second_task": secondTask.Value,
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

		var challenge models.Challenge
		result := db.DB.First(&challenge, "id = ?", answer.ID)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
			return
		}

		processingTime := time.Since(challenge.CreatedAt)
		noiseDetect := challenge.ExpectedHash != answer.FirstTaskHash

		challenge.NoiseDetected = noiseDetect

		// Create a map for selective update
		updateData := map[string]interface{}{
			"NoiseDetected":  noiseDetect,
			"ActualHash":     answer.FirstTaskHash,
			"Fingerprint":    answer.SecondTaskHash,
			"ProcessingTime": processingTime.Milliseconds(),
			"JavaScript":     true,
		}

		if answer.DiffTaskHash != nil {
			updateData["NoiseHash"] = *answer.DiffTaskHash
		}

		if err := db.DB.Model(&challenge).Updates(updateData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update challenge in cache"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":         "ok",
			"noise_detected": noiseDetect,
		})
	})

	router.GET("/", func(c *gin.Context) {
		// Generate new challenge ID
		id := uuid.New()
		challenge := models.Challenge{
			ID:        id.String(),
			ExpiresAt: time.Now().Add(challengeExpiration),
		}
		result := db.DB.Create(&challenge)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create challenge"})
			return
		}

		// Read index.html
		root, err := os.Getwd()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get working directory")
			return
		}
		filePath := filepath.Join(root, "resources", "index.html")
		htmlContent, err := os.ReadFile(filePath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read index.html")
			return
		}

		// Replace placeholder with the real ID
		newHTML := strings.Replace(string(htmlContent), "__CHALLENGE_ID__", id.String(), 1)

		// Serve the modified HTML
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(newHTML))
	})

	router.GET("/predictor.worker.js", func(c *gin.Context) {
		root, err := os.Getwd()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to get working directory")
			return
		}

		filePath := filepath.Join(root, "resources", "predictor.worker.js")

		c.Header("Content-Type", "application/javascript")
		c.File(filePath)
	})
	log.Printf("Server starting on port %s", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Printf("Server can't be started %s", err.Error())
		return
	}
}
