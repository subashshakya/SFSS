package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/subashshakya/SFSS/constants"
	"github.com/subashshakya/SFSS/db/orms"
	"github.com/subashshakya/SFSS/models"
)

func init() {
	validate = validator.New()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func ShareSecureFile(c *gin.Context) {
	var shareFile models.FileSharing
	if err := c.ShouldBindJSON(&shareFile); err != nil {
		log.Println(constants.BadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.BadRequest})
		return
	}
	if err := validate.Struct(&shareFile); err != nil {
		log.Println(constants.ValidationError, err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.ValidationError})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShortTimeout)
	defer cancel()
	err := orms.ShareFile(ctx, &shareFile)
	if err != nil {
		log.Println("Transaction not successful: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Successfully shared file"})
}

func ShareSuperSecret(c *gin.Context) {
	var superSecret models.SecretSharing
	if err := c.ShouldBindJSON(&superSecret); err != nil {
		log.Println(constants.BadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.BadRequest})
		return
	}
	if err := validate.Struct(&superSecret); err != nil {
		log.Println(constants.ValidationError, err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.ValidationError})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShortTimeout)
	defer cancel()
	err := orms.ShareSecret(ctx, &superSecret)
	if err != nil {
		log.Println("Could not save in the db:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Successfully shared secret"})
}

func GetFileSharedOfAUser(c *gin.Context) {
	var userSecrets []*models.FileSharing
	senderId, err := strconv.ParseInt(c.Param("id"), 10, 0)
	if err != nil {
		log.Println("ID parsing error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.BadRequest})
		return
	}
	if senderId == 0 {
		log.Println("ID is zero")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.BadRequest})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShortTimeout)
	defer cancel()
	userSecrets, err = orms.GetFileSharesOfAUser(ctx, uint(senderId))
	if err != nil {
		log.Println("Could not get files: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Successfully fetched user files", "data": userSecrets})
}

func GetSecretSharedOfAUser(c *gin.Context) {
	var userSecrets []*models.SecretSharing
	senderId, err := strconv.ParseInt(c.Param("id"), 10, 0)
	if err != nil {
		log.Println("ID parsing error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.BadRequest})
		return
	}
	if senderId == 0 {
		log.Println("ID is zero")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.BadRequest})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShortTimeout)
	defer cancel()
	userSecrets, err = orms.GetSecretSharesOfAUser(ctx, uint(senderId))
	if err != nil {
		log.Println("Could not get secrets: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Successfully fetched user files", "data": userSecrets})
}
