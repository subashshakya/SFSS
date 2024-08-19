package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/subashshakya/SFSS/db/orms"
	"github.com/subashshakya/SFSS/models"
	"github.com/subashshakya/SFSS/utils"
)

const internalServerError = "Internal Server Error"
const invalidInput = "Invalid data input"

func UserSignUp(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidInput})
		return
	}

	ctx := context.Background()

	existingUser, err := orms.GetUser(ctx, newUser.Id)
	if err == nil && existingUser.Id != 0 {
		c.JSON(http.StatusConflict, gin.H{"success": false, "error": "User already exists"})
		return
	}

	if _, err := orms.CreateUser(ctx, &newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Could not save user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User sign up successful", "data": newUser})
}

func UserSignIn(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidInput})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	isPresent, err := orms.CheckUserInDB(ctx, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": internalServerError})
	}
	if !isPresent {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Could not match user credentials."})
	}
	token, tokenError := utils.GenerateToken(user.Id)
	if tokenError != nil || token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not generate token", "error": tokenError})
	}
	c.JSON(http.StatusAccepted, gin.H{"success": true, "message": "Sign-In Successful", "token": token})
}

func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": invalidInput})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Could not parse user id"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dbUser, dbErr := orms.GetUser(ctx, uint(id))
	if dbErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": internalServerError})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "The requested user is fetched successfully", "data": dbUser})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": invalidInput})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	isUserInDB, _ := orms.CheckUserInDB(ctx, &user)
	if !isUserInDB {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found for modification"})
	}
	isUpdated, error := orms.UpdateUser(ctx, &user)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": internalServerError})
	}
	if !isUpdated {
		c.JSON(http.StatusNotModified, gin.H{"success": false, "message": "Data not modified"})
	}
	c.JSON(http.StatusAccepted, gin.H{"success": true, "message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": invalidInput})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	isUserInDB, _ := orms.CheckUserInDB(ctx, &user)
	if !isUserInDB {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found for modification"})
	}
	deleteSuccess, error := orms.DeleteUser(ctx, &user)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": internalServerError})
	}
	if !deleteSuccess {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "message": "Could not process the requested process"})
	}
	c.JSON(http.StatusAccepted, gin.H{"success": true, "message": "Successfully deleted the user"})
}
