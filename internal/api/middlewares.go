package api

import (
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

// if an user is logged in (indicated by the JWT in the cookie), then we check the JWT
// and extract the user for all handlers to use
func loginDetector() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok := checkJWT(c)
		if !ok {
			logger.Info("No JWT detected. A Guest!")
		} else {
			// u := getCurrentAuthUser(c)
			// logger.Info("Browsing as user", u.Email)
		}

		// c.Next()
	}
}

func checkJWT(c *gin.Context) bool {
	// get the raw jwt from cookie
	tokenString, err := c.Cookie("auth")
	if err != nil {
		logger.Info("error getting jwt from cookies:", err, ". Assuming not logged in.")
		return false
	}

	if tokenString == "" {
		logger.Warn("blank cookie. assuming not logged in")
		return false
	}

	// validate and stuff
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		logger.Warn("error parsing jwt:", err)
		return false
	}
	logger.Info(claims)

	// we have the claims. does this user actually exist in db?
	// user := gorm.User{}
	// err = user.PopulateByID(claims.UserID)
	// if err != nil {
	// 	logger.Error("error looking up user from jwt:", err)
	// 	return false
	// }

	// save the user and the claim to gin's storage
	// c.Set("currentUser", &user)
	// c.Set("claims", claims)

	return true
}

// get current user from gin store. If not logged in, will return a nil pointer
// func getCurrentAuthUser(c *gin.Context) *gorm.User {
// 	if user, ok := c.Get("currentUser"); ok && user != nil {
// 		if user2, ok2 := user.(*gorm.User); ok2 && user2 != nil {
// 			return user2
// 		}
// 	}
// 	return nil
// }
