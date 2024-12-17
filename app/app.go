package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/l1qwie/JWTAuth/app/database"
	errh "github.com/l1qwie/JWTAuth/app/errorhandler"
	"github.com/l1qwie/JWTAuth/app/types"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

const (
	fromemail string = "cogratulationservice@gmail.com"
	subject   string = "!!!WARNING!!!"
)

func isValidGUID(guid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(guid)
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
			encodedString = base64.StdEncoding.EncodeToString(append(randB, []byte("::"+ip)...))
			if bcryptHash, err = bcrypt.GenerateFromPassword(randB, bcrypt.DefaultCost); err == nil {
				err = database.Conn.SaveRefreshToken(bcryptHash, guid, &onceagain)
			}
		}
	}
	return encodedString, err
}

func sendEmail(ip string) error {
	message := fmt.Sprintf("Someone is trying to login into your account from a diffrent device! Their IP is %s. If this is you just ignore the message.", ip)
	to, err := database.Conn.SelectEmail(ip)
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
	ok, err := database.Conn.CheckGUID(guid)
	if !ok && err == nil {
		err = errh.Code10()
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
	if isValidGUID(guid) {
		if net.ParseIP(userIP) != nil {
			tokens := new(types.Tokens)
			if err = isThereTheGUID(guid); err == nil {
				body, err = newBothTokens(tokens, userIP, guid)
			}
		} else {
			err = errh.Code11()
		}
	} else {
		err = errh.Code12()
	}
	return body, err
}

func isIPv4(token []byte) (net.IP, []byte, bool) {
	str := string(token)
	parts := strings.Split(str, "::")
	if len(parts) > 0 {
		ipStr := parts[len(parts)-1]
		if ip := net.ParseIP(ipStr); ip != nil && ip.To4() != nil {
			return ip, []byte(strings.Join(parts[:len(parts)-1], "::")), true
		}
	}
	return nil, nil, false
}

func isIPv6(token []byte) (net.IP, []byte, bool) {
	str := string(token)
	parts := strings.Split(str, "::")
	if len(parts) > 0 {
		ipStr := parts[len(parts)-1]
		if ip := net.ParseIP(ipStr); ip != nil && ip.To16() != nil {
			return ip, []byte(strings.Join(parts[:len(parts)-1], "::")), true
		}
	}
	return nil, nil, false
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
			err = errh.Code13()
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
				if err = database.Conn.RewriteIP(ip.String(), clientIP); err == nil {
					trueIp = clientIP
					err = sendEmail(trueIp)
				}
			} else {
				trueIp = ip.String()
			}
			if err == nil {
				if bcryptHash, err = database.Conn.GetRefreshToken(trueIp); err == nil {
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
	if net.ParseIP(ip) != nil {
		if reftoken != "" {
			tokens := new(types.Tokens)
			if err = checkRefreshToken(reftoken, ip); err == nil {
				if guid, err = database.Conn.GetGUID(ip); err == nil {
					body, err = newBothTokens(tokens, ip, guid)
				}
			}
		} else {
			err = errh.Code14()
		}
	} else {
		err = errh.Code11()
	}
	return body, err
}
