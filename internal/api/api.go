package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type APIServer struct {
	DataLayer APIDataLayer

	router *gin.Engine
}

func NewAPIServer(datalayer APIDataLayer) *APIServer {
	sv := APIServer{
		DataLayer: datalayer,
	}
	sv.router = sv.initRoutes()
	return &sv
}

func (s *APIServer) initRoutes() *gin.Engine {
	router := gin.New()

	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// CORS babyyy
	router.Use(corsMiddleware())

	// Sessions middleware
	store := cookie.NewStore([]byte(viper.GetString("session_cookie_key")))
	store.Options(sessions.Options{
		HttpOnly: true,
		MaxAge:   604800, // a week
		Path:     "/",
	})
	router.Use(sessions.Sessions("virgo_session", store))

	// login detector
	// router.Use(loginDetector())

	// static
	// router.Use(static.Serve("/static", static.LocalFile("./static", true)))
	// router.StaticFile("/favicon.ico", "./static/favicon.ico")

	// ping
	router.GET("/ping", handlePing)

	v1 := router.Group("/v1")

	userRoutes := v1.Group("/users")
	userRoutes.POST("/checkAvailability", s.handleCheckAvailability)
	userRoutes.POST("/register", s.registerNewUser)

	return router
}

func (s *APIServer) Run() {
	port := viper.GetString("port")
	s.router.Run(":" + port)
}

func handlePing(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "pong"})
}
