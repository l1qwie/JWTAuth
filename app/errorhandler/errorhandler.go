package errorhandler

import "github.com/l1qwie/JWTAuth/app/types"

func writeMsgAndCode(code int, msg string, err *types.Err) {
	err.Code = code
	err.Msg = msg
}

func Code14() error {
	err := new(types.Err)
	writeMsgAndCode(14, "a refresh token is required", err)
	return err
}

func Code13() error {
	err := new(types.Err)
	writeMsgAndCode(13, "an ip in a refresh token is required", err)
	return err
}

func Code12() error {
	err := new(types.Err)
	writeMsgAndCode(12, "invalid guid", err)
	return err
}

func Code11() error {
	err := new(types.Err)
	writeMsgAndCode(11, "invalid ip", err)
	return err
}

func Code10() error {
	err := new(types.Err)
	writeMsgAndCode(10, "the guid is unknown", err)
	return err
}

func Code01() error {
	err := new(types.Err)
	writeMsgAndCode(01, "invalid input data", err)
	return err
}
