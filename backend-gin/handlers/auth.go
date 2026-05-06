package handlers

import (
	"context"
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

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

type forgetPassword struct {
	Email string `json:"email" binding:"required,email"`
}

func LoginHandler(c *gin.Context) {

	var req Login

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Collection("users")

	var user models.User

	// detect email or username
	filter := bson.M{}

	if _, err := mail.ParseAddress(req.Identifier); err == nil {
		filter = bson.M{"email":req.Identifier}
	}else {
		filter = bson.M{"username":req.Identifier}
	}

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}

func RegisterHandler(c *gin.Context) {

	var req Register

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Collection("users")

	var user models.User

	// check if user exists
	err := collection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": req.Username},
			{"email": req.Email},
		},
	}).Decode(&user)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or email already exists"})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	newUser := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	res, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	token, err := utils.GenerateToken(oid.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "registered successfully",
		"token":   token,
	})
}

func ForgetPasswordHandler(c *gin.Context) {

	var req forgetPassword

	if err :=c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var user models.User
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.DB.Collection("users")

	err := collection.FindOne(ctx, bson.M{
		"email" : req.Email,
	}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email not found"})
		return
	}

	// generate reset token
	resetToken, err := utils.GenerateResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate reset token"})
		return
	}

	//send reset email
	err = utils.SendResetEmail(user.Email, resetToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send reset email"})
		return
	}

	//Hash reset token
	hashedResetToken, err := bcrypt.GenerateFromPassword([]byte(resetToken), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash reset token"})
		return
	}

	// save reset token to user
	_, err = collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"reset_token": string(hashedResetToken),
			"reset_token_expires": time.Now().Add(10 * time.Second),
		},
	})	
}