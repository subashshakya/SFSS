package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/subashshakya/SFSS/constants"
	"github.com/subashshakya/SFSS/db/orms"
	"github.com/subashshakya/SFSS/models"
	"github.com/subashshakya/SFSS/utils"

	"log"
	"strconv"
)

var validate *validator.Validate

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	validate = validator.New()
}

func checkInvalidToken(c *gin.Context) bool {
	if err := utils.TokenValid(c); err != nil {
		log.Println(constants.Unauthorized)
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": constants.Unauthorized})
		return false
	}
	return true
}

func GetUserFiles(c *gin.Context) {
	var userFiles []models.SecureFile
	tokenIsValid := checkInvalidToken(c)
	if !tokenIsValid {
		return
	}
	userId, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	if userId == 0 {
		log.Println("Id invalid")
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": constants.BadRequest})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, _ := orms.GetUser(ctx, uint(userId))
	if user.Id == 0 {
		log.Println("Could not find user")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Could not find user"})
		return
	}
	userFiles, isNilErr := orms.GetSecureFilesOfAUser(ctx, int(userId))
	if isNilErr != nil {
		log.Println(isNilErr)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": false, "message": "Successfully fetched user files", "data": userFiles})
}

func UpdateSecureFile(c *gin.Context) {
	var secureFile models.SecureFile
	tokenIsValid := checkInvalidToken(c)
	if !tokenIsValid {
		return
	}
	if err := c.ShouldBindJSON(&secureFile); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.BadRequest})
		return
	}
	if err := validate.Struct(&secureFile); err != nil {
		log.Println("Validation Error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.ValidationError})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	updatedFile, err := orms.UpdateFile(ctx, &secureFile)
	if err != nil {
		log.Println("Failed to update the file: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Updated the file successfully", "data": updatedFile})
}

func MakeSecureFile(c *gin.Context) {
	var secureFile models.SecureFile
	tokenIsValid := checkInvalidToken(c)
	if !tokenIsValid {
		log.Println("Token is not valid")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": constants.Unauthorized})
		return
	}
	if err := c.ShouldBindJSON(&secureFile); err != nil {
		log.Println("JSON Binding Error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.BadRequest})
		return
	}
	if err := validate.Struct(&secureFile); err != nil {
		log.Println("Validation Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.ValidationError})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	createSuccess, err := orms.CreateSecureFile(ctx, &secureFile)
	if err != nil || !createSuccess {
		log.Println("Could not save the file: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"success": false, "message": "Created File Successfully"})
}

func DeleteFile(c *gin.Context) {
	tokenIsValid := checkInvalidToken(c)
	if !tokenIsValid {
		log.Println("Token is invalid")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": constants.Unauthorized})
		return
	}
	fileId := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := orms.GetSecureFileById(ctx, fileId)
	if err != nil {
		log.Println("File not found", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "File not found"})
		return
	}
	deleteSuccess, err := orms.DeleteSecureFile(ctx, fileId)
	if err != nil && !deleteSuccess {
		log.Println("Could not delete file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	if err == nil && !deleteSuccess {
		log.Println("Delete not successful")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Delete not successful"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Successfully deleted file"})
}

func GetSecureFileByID(c *gin.Context) {
	tokenIsValid := checkInvalidToken(c)
	if !tokenIsValid {
		log.Println("Token is invalid")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": constants.Unauthorized})
		return
	}
	fileId := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	secureFile, err := orms.GetSecureFileById(ctx, fileId)
	if err != nil {
		log.Println("Could not find the file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	if secureFile == nil {
		log.Println("Not found")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "File not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": false, "message": "Successfully fetched file", "data": secureFile})
}
