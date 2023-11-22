package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

// User represents a user in the system
type User struct {
	ID                int
	Name              string
	PhoneNumber       string
	OTP               string
	OTPExpirationTime time.Time
}

// GenerateOTPRequest represents the request to generate an OTP for a user
type GenerateOTPRequest struct {
	PhoneNumber string `json:"phone_number"`
}

// VerifyOTPRequest represents the request to verify an OTP for a user
type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
}

func initDB() {
	connStr := "postgres://amakhyon:87654321@localhost/Aqari?sslmode=disable"
	pool, err := pgxpool.Connect(nil, connStr)
	if err != nil {
		log.Fatal(err)
	}

	db = pool
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	r.POST("/api/users", createUserHandler)
	r.POST("/api/users/generateotp", generateOTPHandler)
	r.POST("/api/users/verifyotp", verifyOTPHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// createUserHandler handles the HTTP POST request to create a new user
func createUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if isPhoneNumberExists(user.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
		return
	}

	if err := createUserDB(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// generateOTPHandler handles the HTTP POST request to generate an OTP for a user
func generateOTPHandler(c *gin.Context) {
	var request GenerateOTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := getUserByPhoneNumber(request.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	otp := generateRandomOTP()
	expirationTime := time.Now().Add(1 * time.Minute)

	if err := updateOTP(user.ID, otp, expirationTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"otp": otp, "expiration_time": expirationTime})
}

// verifyOTPHandler handles the HTTP POST request to verify an OTP for a user
func verifyOTPHandler(c *gin.Context) {
	var request VerifyOTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := getUserByPhoneNumber(request.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.OTP != request.OTP || user.OTPExpirationTime.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

// isPhoneNumberExists checks if a phone number already exists in the database
func isPhoneNumberExists(phoneNumber string) bool {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE phone_number = $1"
	if err := db.QueryRow(context.Background(), query, phoneNumber).Scan(&count); err != nil {
		log.Println(err)
		return false
	}
	return count > 0
}

// createUserDB creates a new user in the database
func createUserDB(user User) error {
	res, err := db.Exec(context.Background(), CreateUser, user.Name, user.PhoneNumber)
	if err != nil {
		return err
	}

	return res.Err()
}

// generateRandomOTP generates a random 4-digit OTP
func generateRandomOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

// getUserByPhoneNumber retrieves a user from the database by phone number
func getUserByPhoneNumber(phoneNumber string) (User, error) {
	var user User
	err := db.QueryRow(context.Background(), GetUserByPhoneNumber, phoneNumber).
		Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.OTP, &user.OTPExpirationTime)
	return user, err
}

// updateOTP updates the OTP and its expiration time for a user in the database
func updateOTP(userID int, otp string, expirationTime time.Time) error {
	_, err := db.Exec(context.Background(), UpdateOTP, otp, expirationTime, userID)
	return err
}

const CreateUser = `-- name: CreateUser
INSERT INTO users (name, phone_number) VALUES ($1, $2)`

const GetUserByPhoneNumber = `-- name: GetUserByPhoneNumber
SELECT id, name, phone_number, otp, otp_expiration_time FROM users WHERE phone_number = $1`

const UpdateTOP = `-- name: UpdateOTP
UPDATE users SET otp = $1, otp_expiration_time = $2 WHERE id = $3`
