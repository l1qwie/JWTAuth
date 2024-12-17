package main

import (
	"github.com/l1qwie/JWTAuth/api"
	"github.com/l1qwie/JWTAuth/app/logs"
)

func main() {
	logs.SetDebug()
	api.StartAPI()
}
