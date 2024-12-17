package tests

import (
	"os"
	"time"

	"github.com/l1qwie/JWTAuth/api"
)

func PutEnvVal() {
	if err := os.Setenv("host_db", "localhost"); err != nil {
		panic(err)
	}
	if err := os.Setenv("port_db", "3333"); err != nil {
		panic(err)
	}
	if err := os.Setenv("user_db", "postgres"); err != nil {
		panic(err)
	}
	if err := os.Setenv("password_db", "postgres"); err != nil {
		panic(err)
	}
	if err := os.Setenv("dbname_db", "postgres"); err != nil {
		panic(err)
	}
	if err := os.Setenv("sslmode_db", "disable"); err != nil {
		panic(err)
	}
}

func StartAPI() {
	go api.StartAPI()
	time.Sleep(time.Second)
}
