package api

import (
	"net/http"

	"github.com/cs5224virgo/virgo/internal/jwt"
	"github.com/cs5224virgo/virgo/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type APIServer struct {
	DataLayer    APIDataLayer
	WebSocketHub WebSocketHub

	router *gin.Engine
}

func NewAPIServer(datalayer APIDataLayer, hub WebSocketHub) *APIServer {
	sv := APIServer{
		DataLayer:    datalayer,
		WebSocketHub: hub,
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
	// store := cookie.NewStore([]byte(viper.GetString("session_cookie_key")))
	// store.Options(sessions.Options{
	// 	HttpOnly: true,
	// 	MaxAge:   604800, // a week
	// 	Path:     "/",
	// })
	// router.Use(sessions.Sessions("virgo_session", store))

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
	userRoutes.POST("/login", s.userLogin)
	userRoutes.GET("/wstoken", s.authMiddleware, s.getUserWsToken)

	roomRoutes := v1.Group("/rooms", s.authMiddleware)
	roomRoutes.GET("/", s.handleGetRooms)
	roomRoutes.POST("/new", s.handleCreateRoom)

	v1.GET("/ws", s.handleWebSocket)

	return router
}

func (s *APIServer) Run() {
	port := viper.GetString("port")
	s.router.Run(":" + port)
}

func handlePing(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "pong"})
}

func (s *APIServer) handleWebSocket(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		logger.Error("empty token string")
		failureResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		logger.Error(err)
		failureResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	user, err := s.DataLayer.GetUserByID(claims.UserID)
	if err != nil {
		logger.Error("error looking up user from jwt:", err)
		failureResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	s.WebSocketHub.ServeWs(c, user.Username)
}
