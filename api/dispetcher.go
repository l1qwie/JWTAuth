package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/l1qwie/JWTAuth/app"
	"github.com/l1qwie/JWTAuth/app/database"
	"github.com/l1qwie/JWTAuth/app/logs"
	"github.com/l1qwie/JWTAuth/app/types"
)

type server struct {
	router *gin.Engine
}

func code400() error {
	err := new(types.Err)
	err.Code = http.StatusBadRequest
	err.Msg = "invalid query parameters"
	return err
}

func newConnection() {
	var err error
	database.Conn, err = database.Connect()
	if err != nil {
		panic(err)
	}
}

func newServer() *server {
	s := new(server)
	s.router = gin.Default()
	return s
}

func configurations() *server {
	newConnection()
	return newServer()
}

func response(ctx *gin.Context, msg []byte, err error) {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	} else {
		ctx.JSON(http.StatusOK, msg)
	}
}

func loginAction(ctx *gin.Context) {
	var err error
	var body []byte
	guid := ctx.Query("id")
	justcall := ctx.Query("justacall")
	if guid != "" {
		if justcall != "true" {
			body, err = app.NewAccessAndRefreshTokens(guid, ctx.ClientIP())
		}
	} else {
		logs.ParameterIsRequired("id")
		err = code400()
	}
	response(ctx, body, err)
}

func (s *server) login() {
	path := "/login"
	s.router.GET(path, loginAction)
	logs.StartPoint(path, "GET")
}

func refreshAction(ctx *gin.Context) {
	var err error
	var body []byte
	justcall := ctx.Query("justacall")
	token := ctx.Request.Header["Refresh-Token"][0]
	if token != "" {
		if justcall != "true" {
			body, err = app.RefreshAction(ctx.ClientIP(), token)
		}
	} else {
		logs.ParameterIsRequired("Refresh-Token")
		err = code400()
	}
	response(ctx, body, err)
}

func (s *server) refresh() {
	path := "/refresh"
	s.router.PATCH(path, refreshAction)
	logs.StartPoint(path, "PATCH")
}

func StartAPI() {
	s := configurations()
	s.login()
	s.refresh()

	s.router.Run(":8080")
}
