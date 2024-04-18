package api

import (
	"net/http"
	"strings"

	sqlc "github.com/cs5224virgo/virgo/db/generated"
	"github.com/cs5224virgo/virgo/internal/jwt"
	"github.com/cs5224virgo/virgo/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", viper.GetString("frontend_url"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

// authMiddleware : make sure restricted paths are protected
func (s *APIServer) authMiddleware(c *gin.Context) {
	ok := s.checkJWT(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

func (s *APIServer) checkJWT(c *gin.Context) bool {
	// get the raw jwt from cookie
	rawHeader := c.GetHeader("Authorization")
	tokenString, ok := strings.CutPrefix(rawHeader, "Bearer ")
	if !ok {
		logger.Warn("invalid auth header. assuming not logged in")
		return false
	}

	if tokenString == "" {
		logger.Warn("blank auth header. assuming not logged in")
		return false
	}

	// validate and stuff
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		logger.Warn("error parsing jwt:", err)
		return false
	}
	// logger.Info(claims)

	// we have the claims. does this user actually exist in db?
	user, err := s.DataLayer.GetUserByID(claims.UserID)
	if err != nil {
		logger.Error("error looking up user from jwt:", err)
		return false
	}

	// save the user and the claim to gin's storage
	c.Set("currentUser", user)
	c.Set("claims", claims)

	return true
}

// get current user from gin store. If not logged in, will return a nil pointer
func getCurrentAuthUser(c *gin.Context) *sqlc.User {
	if user, ok := c.Get("currentUser"); ok && user != nil {
		if user2, ok2 := user.(*sqlc.User); ok2 && user2 != nil {
			return user2
		}
	}
	return nil
}
