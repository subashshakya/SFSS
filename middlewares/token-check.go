package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/subashshakya/SFSS/constants"
	"github.com/subashshakya/SFSS/utils"
)

func CheckInvalidToken() gin.HandlerFunc {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return func(c *gin.Context) {
		if isTokenValid := checkInvalidToken(c); !isTokenValid {
			log.Println("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": constants.Unauthorized})
			c.Abort()
			return
		}
		c.Next()
	}
}

func checkInvalidToken(c *gin.Context) bool {
	if err := utils.TokenValid(c); err != nil {
		return false
	}
	return true
}
