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
