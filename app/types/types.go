package types

import (
	"fmt"
)

type Tokens struct {
	Access  string `json:"access_jwt_key"`
	Refresh string `json:"refresh_jwt_key"`
}

type Err struct {
	Code int
	Msg  string
}

func (e *Err) Error() string {
	return fmt.Sprintf("[ERROR:%d] %s", e.Code, e.Msg)
}
