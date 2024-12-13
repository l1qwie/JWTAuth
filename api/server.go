package api

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/l1qwie/JWTAuth/app"
	"github.com/l1qwie/JWTAuth/app/database"
	"github.com/l1qwie/JWTAuth/app/logs"
	"github.com/l1qwie/JWTAuth/app/types"
)

type server struct {
	router *gin.Engine
}

func code500() error {
	err := new(types.Err)
	err.Code = http.StatusBadRequest
	err.Msg = "invalid query parameters"
	return err
}

func isValidGUID(guid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(guid)
}

func newConnection() {
	var err error
	types.Conn, err = database.Connect()
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
	guid := ctx.Param("id")
	if guid != "" {
		if isValidGUID(guid) {
			// app logic
			body, err = app.NewAccessAndRefreshTokens(guid, ctx.ClientIP())
		} else {
			logs.ParameterIsRequired("GUID")
			err = code500()
		}
	} else {
		logs.ParameterIsRequired("id")
		err = code500()
	}
	response(ctx, body, err)
}

func (s *server) login() {
	path := "/login/:id"
	s.router.GET(path, loginAction)
	logs.StartPoint(path, "GET")
}

func refreshAction(ctx *gin.Context) {
	var err error
	var body []byte
	token := ctx.Request.Header["Refresh-Token"][0]
	if token != "" {
		// app logic
	} else {
		logs.ParameterIsRequired("Refresh-Token")
		err = code500()
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
}
