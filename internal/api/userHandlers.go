package api

import (
	"errors"
	"fmt"
	"net/http"

	sqlc "github.com/cs5224virgo/virgo/db/generated"
	"github.com/cs5224virgo/virgo/internal/datalayer"
	"github.com/cs5224virgo/virgo/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (s *APIServer) handleCheckAvailability(c *gin.Context) {
	type reqStruct struct {
		Username string `json:"username"`
	}
	var req reqStruct
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		failureResponse(c, http.StatusBadRequest, "")
		return
	}

	type resStruct struct {
		BaseResponse
		IsAvailable bool `json:"isAvailable"`
	}
	var res resStruct
	res.IsAvailable = true
	res.Success = true
	if req.Username != "" {
		isAvailable, err := s.DataLayer.IsUsernameAvailable(req.Username)
		if err != nil {
			logger.Error(err)
			failureResponse(c, http.StatusInternalServerError, "")
			return
		}
		res.IsAvailable = isAvailable
	}

	c.JSON(http.StatusOK, res)
}

func (s *APIServer) registerNewUser(c *gin.Context) {
	type reqStruct struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName"`
	}
	var req reqStruct
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		failureResponse(c, http.StatusBadRequest, "")
		return
	}

	dbparams := sqlc.CreateUserParams{
		Username: req.Username,
	}
	pepperedPassword := viper.GetString("password_pepper") + req.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pepperedPassword), bcrypt.MinCost)
	if err != nil {
		logger.Error("unable to hash user password:", err)
		failureResponse(c, http.StatusInternalServerError, "")
		return
	}
	dbparams.Password = string(hashedPassword)
	dbparams.DisplayName.Valid = false
	if req.DisplayName != "" {
		dbparams.DisplayName.Valid = true
		dbparams.DisplayName.String = req.DisplayName
	}
	err = s.DataLayer.CreateUser(dbparams)
	if err != nil {
		logger.Error("unable to create user:", err)
		failureResponse(c, http.StatusInternalServerError, "")
		return
	}
	c.JSON(http.StatusCreated, BaseResponse{
		Success: true,
	})
}

func (s *APIServer) userLogin(c *gin.Context) {
	type reqStruct struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req reqStruct
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		failureResponse(c, http.StatusBadRequest, "")
		return
	}

	type userDetailsStruct struct {
		Username    string  `json:"username"`
		DisplayName *string `json:"displayName"`
	}
	type userDataStruct struct {
		Message     string            `json:"message"`
		UserDetails userDetailsStruct `json:"userDetails"`
	}
	type resStruct struct {
		BaseResponse
		Authorization string         `json:"authorization"`
		Data          userDataStruct `json:"data"`
	}
	var res resStruct

	pepperedPassword := viper.GetString("password_pepper") + req.Password
	user, token, err := s.DataLayer.AuthenticateUser(req.Username, pepperedPassword)
	if err != nil {
		if errors.Is(err, datalayer.ErrLoginFailed) {
			failureResponse(c, http.StatusUnauthorized, "Incorrect login")
			return
		}
		failureResponse(c, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	res.Success = true
	res.Authorization = token
	res.Data.UserDetails.Username = user.Username
	if user.DisplayName.Valid {
		res.Data.UserDetails.DisplayName = &user.DisplayName.String
	}
	c.JSON(http.StatusOK, res)
}

func failureResponse(c *gin.Context, code int, message string) {
	c.JSON(code, BaseResponse{
		Success: false,
		Message: message,
	})
	c.Abort()
}
