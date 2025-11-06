package main

import (
	"log"
	"memory-vault-api/internal/database"
	"memory-vault-api/internal/handlers"
	"memory-vault-api/internal/queue"
	"memory-vault-api/internal/services"
	"memory-vault-api/internal/storage"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "memory_vault")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	jwtSecret := getEnv("JWT_SECRET", "default-secret-key")
	s3Endpoint := getEnv("S3_ENDPOINT", "http://localhost:9000")
	s3AccessKey := getEnv("S3_ACCESS_KEY", "minioadmin")
	s3SecretKey := getEnv("S3_SECRET_KEY", "minioadmin")
	s3Bucket := getEnv("S3_BUCKET", "memory-vault")

	// Initialize database
	db, err := database.NewPostgresDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient := queue.NewRedisClient(redisAddr)
	defer redisClient.Close()

	// Initialize S3 storage
	s3Client, err := storage.NewS3Client(s3Endpoint, s3AccessKey, s3SecretKey)
	if err != nil {
		log.Fatal("Failed to initialize S3 client:", err)
	}

	// Initialize services
	authService := services.NewAuthService(jwtSecret)
	fileService := services.NewFileService(db, redisClient, s3Client, s3Bucket)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, db)
	fileHandler := handlers.NewFileHandler(fileService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(authHandler.JWTMiddleware())
	{
		api.POST("/upload", fileHandler.UploadFile)
		api.GET("/files", fileHandler.GetFiles)
		api.GET("/files/:id", fileHandler.GetFile)
		api.DELETE("/files/:id", fileHandler.DeleteFile)
		api.GET("/ws", fileHandler.WebSocketHandler)
	}

	log.Println("Server starting on :8080")
	r.Run(":8080")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
