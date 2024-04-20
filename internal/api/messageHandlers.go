package api

import (
	"fmt"
	"net/http"

	"github.com/cs5224virgo/virgo/internal/datalayer"
	"github.com/cs5224virgo/virgo/logger"
	"github.com/gin-gonic/gin"
)

func (s *APIServer) handleGetMessages(c *gin.Context) {
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

	messages, err := s.DataLayer.GetAllMessagesForRoom(req.RoomCode)
	if err != nil {
		failureResponse(c, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	type resRoomStruct struct {
		Messages []datalayer.DetailedMessage `json:"messages"`
	}
	type resStruct struct {
		BaseResponse
		Data resRoomStruct `json:"data"`
	}
	var res resStruct
	res.Success = true
	res.Data.Messages = messages
	c.JSON(200, res)
}

func (s *APIServer) handleSummary(c *gin.Context) {
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

	// get messages from database
	messages, err := s.DataLayer.GetLastWeekMessagesForRoom(req.RoomCode)
	if err != nil {
		failureResponse(c, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	// combine all the messages together in one big string
	conversation := ""
	for _, msg := range messages {
		datetime := msg.CreatedAt.Format("[15:04, 02/01/2006]")
		conversation = conversation + fmt.Sprintf("%s %s: %s\n", datetime, *msg.User.DisplayName, msg.Content)
	}

	// send to gemini and get response
	response, err := s.AiClient.GetSummary(conversation)
	if err != nil {
		failureResponse(c, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	type resRoomStruct struct {
		Summary string `json:"summary"`
	}
	type resStruct struct {
		BaseResponse
		Data resRoomStruct `json:"data"`
	}
	var res resStruct
	res.Success = true
	res.Data.Summary = response
	c.JSON(200, res)
}
