package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/subashshakya/SFSS/constants"
	"github.com/subashshakya/SFSS/db/orms"
	"github.com/subashshakya/SFSS/models"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	validate = validator.New()
}

func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func CreateSuperSecret(c *gin.Context) {
	var superSecret models.SuperSecret
	if isTokenValid := checkInvalidToken(c); !isTokenValid {
		log.Println("Token not valid")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": constants.Unauthorized})
		return
	}
	if err := c.ShouldBindJSON(&superSecret); err != nil {
		log.Println("Body parsing to JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.BadRequest})
		return
	}
	if err := validate.Struct(&superSecret); err != nil {
		log.Println("Validation Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.ValidationError})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShortTimeout)
	defer cancel()
	success, err := orms.CreateSuperSecret(ctx, &superSecret)
	if err != nil || !success {
		log.Println(constants.InternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Successfully created secret"})
}

func ReadSuperSecret(c *gin.Context) {
	var superSecret *models.SuperSecret
	if isValidToken := checkInvalidToken(c); !isValidToken {
		log.Println(constants.Unauthorized)
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": constants.Unauthorized})
		return
	}
	secretId := c.Param("id")
	if secretId == "" {
		log.Println(constants.IdEmpty)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.IdEmpty})
		return
	}
	if !isValidUUID(secretId) {
		log.Println(constants.UUIDInvalid)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.UUIDInvalid})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShortTimeout)
	defer cancel()
	secretAvailable := orms.IsSecretAvailable(ctx, secretId)
	if !secretAvailable {
		log.Println("Secret Not Available")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": constants.NotFound})
		return
	}
	superSecret, err := orms.GetSecrect(ctx, secretId)
	if err != nil {
		log.Println(constants.InternalServerError, err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	if superSecret == nil {
		log.Println(constants.NotFound)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": constants.NotFound})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Successfully fetched secret", "data": superSecret})
}

func UpdatedSuperSecret(c *gin.Context) {
	var updatedSuperSecret models.SuperSecret
	if isTokenValid := checkInvalidToken(c); !isTokenValid {
		log.Println("Token Invalid")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": constants.Unauthorized})
		return
	}
	if err := c.ShouldBindJSON(&updatedSuperSecret); err != nil {
		log.Println(constants.BadRequest, err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.BadRequest})
		return
	}
	if err := validate.Struct(&updatedSuperSecret); err != nil {
		log.Println(constants.ValidationError, err)
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.ValidationError})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShortTimeout)
	defer cancel()
	updateSuccess, err := orms.UpdateSuperSecret(ctx, &updatedSuperSecret)
	if err != nil || !updateSuccess {
		log.Println(constants.InternalServerError, ":", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Successfully updated secret", "data": updatedSuperSecret})
}

func DeleteSuperSecret(c *gin.Context) {
	if tokenValid := checkInvalidToken(c); !tokenValid {
		log.Println("Token is not valid")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": constants.Unauthorized})
		return
	}
	secretId := c.Param("id")
	if secretId == "" {
		log.Println("Id empty")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.IdEmpty})
		return
	}
	if validUUID := isValidUUID(secretId); !validUUID {
		log.Println("UUID not valid")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.UUIDInvalid})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShortTimeout)
	defer cancel()
	availableSecret, err := orms.GetSecrect(ctx, secretId)
	if availableSecret == nil {
		log.Println("Secret Not Found")
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": constants.NotFound})
		return
	}
	if err != nil {
		log.Println("Error while fetching secret")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	deleteSuccess, err := orms.DeleteSuperSecret(ctx, availableSecret)
	if err != nil && !deleteSuccess {
		log.Println("Error occured while deleting secret:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	if !deleteSuccess && err == nil {
		log.Println("Rows not effected")
		c.JSON(http.StatusNotModified, gin.H{"success": false, "message": "Could not delete secret"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Delete action successful"})
}

func GetSuperSecretsForUser(c *gin.Context) {
	var secrets []models.SuperSecret
	tokenValid := checkInvalidToken(c)
	if !tokenValid {
		log.Println("Token is invalid")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": constants.Unauthorized})
		return
	}
	userId, err := strconv.ParseInt(c.Param("id"), 10, 0)
	if userId == 0 {
		log.Println("Id is zero")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": constants.IdCannotBeZero})
		return
	}
	if err != nil {
		log.Println("Could not parse id")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Id is invalid"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShortTimeout)
	defer cancel()
	secrets, err = orms.GetSecretsOfAUser(ctx, uint(userId))
	if err != nil {
		log.Println("Could not fetch secrets of a user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": constants.InternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Successfully fetched secrets", "data": secrets})
}
