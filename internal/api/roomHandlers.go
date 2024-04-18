package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

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

	du := s.DataLayer.ToDetailedUser(*user)
	go s.WebSocketHub.AnnounceAddUserToRoom(du, room, false)

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

func (s *APIServer) handleLeaveRoom(c *gin.Context) {
	type reqStruct struct {
		RoomCode string `json:"roomCode"`
	}
	var req reqStruct
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		failureResponse(c, http.StatusBadRequest, "")
		return
	}
	user := getCurrentAuthUser(c)

	err = s.DataLayer.LeaveRoom(user.Username, req.RoomCode)
	if err != nil {
		failureResponse(c, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	// announce thru websocket
	du := s.DataLayer.ToDetailedUser(*user)
	s.WebSocketHub.AnnounceLeaveRoom(du, req.RoomCode)

	var res BaseResponse
	res.Success = true
	c.JSON(200, res)
}

func (s *APIServer) handleAddUserToRoom(c *gin.Context) {
	type reqStruct struct {
		Username string `json:"username"`
		RoomCode string `json:"roomCode"`
	}
	var req reqStruct
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		failureResponse(c, http.StatusBadRequest, "")
		return
	}

	regex, err := regexp.Compile(`\s+`)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return
	}
	// Replace all whitespace characters with nothing
	noWhitespace := regex.ReplaceAllString(req.Username, "")
	splittedUsernames := strings.Split(noWhitespace, ",")

	room, err := s.DataLayer.AddUsersToRoom(splittedUsernames, req.RoomCode)
	if err != nil {
		failureResponse(c, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	// announce thru websocket
	for _, user := range room.Users {
		isNewMember := false
		for _, username := range splittedUsernames {
			if username == user.User.Username {
				isNewMember = true
				break
			}
		}
		go s.WebSocketHub.AnnounceAddUserToRoom(user.User, room, isNewMember)
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

func (s *APIServer) handleJoinRoom(c *gin.Context) {
	type reqStruct struct {
		RoomCode string `json:"roomCode"`
	}
	var req reqStruct
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		failureResponse(c, http.StatusBadRequest, "")
		return
	}
	user := getCurrentAuthUser(c)

	room, err := s.DataLayer.JoinRoom(user.Username, req.RoomCode)
	if err != nil {
		failureResponse(c, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	// announce thru websocket
	du := s.DataLayer.ToDetailedUser(*user)
	go s.WebSocketHub.AnnounceAddUserToRoom(du, room, false)

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
