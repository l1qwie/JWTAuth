package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/l1qwie/JWTAuth/app/types"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

const (
	fromemail string = "cogratulationservice@gmail.com"
	subject   string = "!!!WARNING!!!"
)

func code500(msg string) error {
	err := new(types.Err)
	err.Code = http.StatusBadRequest
	err.Msg = msg
	return err
}

func newAccessToken(ip string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"expired_at": time.Now().Add(30 * time.Minute).Unix(),
		"ip":         ip,
	})
	key := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(key)
}

func newRefreshToken(guid, ip string) (string, error) {
	var encodedString string
	var bcryptHash []byte
	var err error
	onceagain := true
	for onceagain {
		randB := make([]byte, 64)
		if _, err = rand.Read(randB); err == nil {
			encodedString = base64.StdEncoding.EncodeToString(append(randB, []byte(ip)...))
			if bcryptHash, err = bcrypt.GenerateFromPassword(randB, bcrypt.DefaultCost); err == nil {
				err = types.Conn.SaveRefreshToken(bcryptHash, guid, &onceagain)
			}
		}
	}
	return encodedString, err
}

func sendEmail(ip string) error {
	message := fmt.Sprintf("Someone is trying to login into your account from a diffrent device! Their IP is %s. If this is you just ignore the message.", ip)
	to, err := types.Conn.SelectEmail(ip)
	if err == nil {
		m := gomail.NewMessage()
		m.SetHeader("From", fromemail)
		m.SetHeader("To", to)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", message)

		d := gomail.NewDialer("smtp.gmail.com", 587, fromemail, "ycuw acml gnor qcir")
		err = d.DialAndSend(m)
	}
	return err
}

func isThereTheGUID(guid string) error {
	ok, err := types.Conn.CheckID(guid)
	if !ok && err == nil {
		err = code500("the guid doesn't exist")
	}
	return err
}

func newBothTokens(tokens *types.Tokens, userIP, guid string) ([]byte, error) {
	var err error
	var body []byte
	if tokens.Access, err = newAccessToken(userIP); err == nil {
		if tokens.Refresh, err = newRefreshToken(guid, userIP); err == nil {
			body, err = json.Marshal(tokens)
		}
	}
	return body, err
}

func NewAccessAndRefreshTokens(guid, userIP string) ([]byte, error) {
	var err error
	var body []byte
	tokens := new(types.Tokens)
	if err = isThereTheGUID(guid); err == nil {
		body, err = newBothTokens(tokens, userIP, guid)
	}
	return body, err
}

func isIPv4(token []byte) (net.IP, []byte, bool) {
	var res bool
	ip := net.IP(token[len(token)-4:])
	if ip.To4() != nil {
		res = true
	}
	return ip, token[:len(token)-4], res
}

func isIPv6(token []byte) (net.IP, []byte, bool) {
	var res bool
	ip := net.IP(token[len(token)-16:])
	if ip.To4() != nil {
		res = true
	}
	return ip, token[:len(token)-16], res
}

func findIP(token []byte) (net.IP, []byte, error) {
	var (
		err       error
		ip        net.IP
		ok        bool
		origtoken []byte
	)
	if ip, origtoken, ok = isIPv4(token); !ok {
		if ip, origtoken, ok = isIPv6(token); !ok {
			err = code500("there isn't an IP in the Refresh-Token")
		}
	}
	return ip, origtoken, err
}

func checkRefreshToken(refreshToken, clientIP string) error {
	var token, originaltoken, bcryptHash []byte
	var err error
	var ip net.IP
	var trueIp string
	if token, err = base64.StdEncoding.DecodeString(refreshToken); err == nil {
		if ip, originaltoken, err = findIP(token); err == nil {
			if !ip.Equal(net.IP(clientIP)) {
				if err = types.Conn.RewriteIP(ip.String(), clientIP); err == nil {
					trueIp = clientIP
					err = sendEmail(trueIp)
				}
			} else {
				trueIp = ip.String()
			}
			if err == nil {
				if bcryptHash, err = types.Conn.GetRefreshToken(trueIp); err == nil {
					err = bcrypt.CompareHashAndPassword(bcryptHash, originaltoken)
				}
			}
		}
	}
	return err
}

func RefreshAction(ip, reftoken string) ([]byte, error) {
	var err error
	var body []byte
	var guid string
	tokens := new(types.Tokens)
	if err = checkRefreshToken(reftoken, ip); err == nil {
		if guid, err = types.Conn.GetGUID(ip); err == nil {
			body, err = newBothTokens(tokens, ip, guid)
		}
	}
	return body, err
}
