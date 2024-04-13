package api

import (
	"fmt"
	"net/http"

	"github.com/cs5224virgo/virgo/internal/datalayer"
	"github.com/gin-gonic/gin"
)

/*
rooms resp:

success bool
data:
	rooms: []
		code: string;
		description: string;
		lastActivity?: Date;
		lastMessagePreview?: string;
		users: Array<RoomUserPopulated>;
			user: User;
				username: string;
				displayName?: string;
			unread: number;
		createdAt: string;



*/

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
