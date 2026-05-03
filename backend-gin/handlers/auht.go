package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"shortlink/config"
	"shortlink/models"
	"shortlink/utils"
)

type Register struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type Login struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {

	var req Login

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "all fields are required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Collection("users")

	var user models.User

	err := collection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": req.Identifier},
			{"email": req.Identifier},
		},
	}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}


	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"message": "login successful",
		"token":   token,
	})
}

func RegisterHandler(c *gin.Context) {

	var req Register

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "all fields are required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Collection("users")

	// check if user exists
	err := collection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": req.Username},
			{"email": req.Email},
		},
	}).Err()

	if err == nil {
		c.JSON(400, gin.H{"error": "username or email already exists"})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to hash password"})
		return
	}

	newUser := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	res, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}
	// convert id from object to string
	id := res.InsertedID.(primitive.ObjectID).Hex()


	token, err := utils.GenerateToken(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"message": "registered successfully",
		"token":   token,
	})
}