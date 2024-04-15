package api

import (
	"fmt"
	"net/http"

	"github.com/cs5224virgo/virgo/internal/datalayer"
	"github.com/cs5224virgo/virgo/logger"
	"github.com/gin-gonic/gin"
)

func (s *APIServer) handleGetRooms(c *gin.Context) {
	user := getCurrentAuthUser(c)
	rooms, err := s.DataLayer.GetRoomsByUserID(user.ID)
	if err != nil {
		failureResponse(c, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	type resRoomStruct struct {
		Rooms []datalayer.DetailedRoom `json:"rooms"`
	}
	type resStruct struct {
		BaseResponse
		Data resRoomStruct `json:"data"`
	}
	var res resStruct
	res.Success = true
	res.Data.Rooms = rooms
	c.JSON(200, res)
}

func (s *APIServer) handleCreateRoom(c *gin.Context) {
	type reqStruct struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}
	var req reqStruct
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		failureResponse(c, http.StatusBadRequest, "")
		return
	}
	user := getCurrentAuthUser(c)

	room, err := s.DataLayer.CreateRoom(user.ID, req.Name, req.Description)
	if err != nil {
		failureResponse(c, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	type resRoomStruct struct {
		Room datalayer.DetailedRoom `json:"room"`
	}
	type resStruct struct {
		BaseResponse
		Data resRoomStruct `json:"data"`
	}
	var res resStruct
	res.Success = true
	res.Data.Room = room
	c.JSON(200, res)
}
