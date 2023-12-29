// package middlewares

// import (
// 	"server/utils"

// 	"github.com/gin-gonic/gin"
// )

// func IsAuthorized() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		cookie, err := c.Cookie("token")

// 		if err != nil {
// 			c.JSON(401, gin.H{"error": "unauthorized"})
// 			c.Abort()
// 			return
// 		}

// 		claims, err := utils.ParseToken(cookie)

//			if err != nil {
//				c.JSON(401, gin.H{"error": "unauthorized"})
//				c.Abort()
//				return
//			}
//			c.Set("user_id", claims.UserID)
//			// c.Set("role", claims.Role) 		решено убрать до лучших времен
//			c.Next()
//		}
//	}
package middlwares

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header from the request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		// Extract the JWT token from the Authorization header
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		err := godotenv.Load()
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "could not load environment variables"})
			return
		}
		// Parse the JWT token and extract the user ID
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// TODO: Replace with your own secret key
			return []byte(os.Getenv("JWT_KEY")), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
			return
		}

		// Save the user ID in the request context
		c.Set("user_id", uint(userID))

		// Call the next middleware function
		c.Next()
	}
}
