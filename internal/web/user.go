package web

import "github.com/gin-gonic/gin"

type UserHandler struct {

}

func NewUserHandler() *UserHandler{
	 return &UserHandler{}
}


func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.
}