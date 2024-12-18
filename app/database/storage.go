package database

import (
	"database/sql"
	"fmt"
	"os"

	errh "github.com/l1qwie/JWTAuth/app/errorhandler"
	"github.com/lib/pq"
)

var Conn *Connection

type Connection struct {
	db *sql.DB
}

func Connect() (*Connection, error) {
	con := new(Connection)

	conninf := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("host_db"),
		os.Getenv("port_db"),
		os.Getenv("user_db"),
		os.Getenv("password_db"),
		os.Getenv("dbname_db"),
		os.Getenv("sslmode_db"))

	db, err := sql.Open("postgres", conninf)
	if err == nil {
		err = db.Ping()
	}
	if err == nil {
		con.db = db
	}
	return con, err
}

func (c *Connection) DeleteUsers() {
	query := "DELETE FROM Users"
	_, err := c.db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (c *Connection) GetRefreshToken(ip string) ([]byte, error) {
	var token []byte
	var err error
	if ip != "" {
		query := "SELECT refreshtoken FROM Users WHERE ip = $1"
		err = c.db.QueryRow(query, ip).Scan(&token)
	} else {
		err = errh.Code01()
	}
	return token, err
}

func (c *Connection) SaveRefreshToken(token []byte, guid string, onceagain *bool) error {
	var err error
	if token != nil && guid != "" && onceagain != nil {
		query := "UPDATE Users SET refreshtoken = $1 WHERE guid = $2"
		_, err = c.db.Exec(query, token, guid)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code != "23505" {
				*onceagain = false
			}
		} else {
			*onceagain = false
		}
	} else {
		err = errh.Code01()
	}
	return err
}

func (c *Connection) CheckGUID(guid string) (bool, error) {
	var count int
	var err error
	if guid != "" {
		query := "SELECT COUNT(*) FROM Users WHERE guid = $1"
		err = c.db.QueryRow(query, guid).Scan(&count)
	} else {
		err = errh.Code01()
	}
	return count == 1, err
}

func (c *Connection) SelectEmail(ip string) (string, error) {
	var email string
	var err error
	if ip != "" {
		query := "SELECT email FROM Users WHERE ip = $1"
		err = c.db.QueryRow(query, ip).Scan(&email)
	} else {
		err = errh.Code01()
	}
	return email, err
}

func (c *Connection) CreateMokData(guid, ip string) error {
	var err error
	if guid != "" && ip != "" {
		query := "INSERT INTO Users (guid, email, ip) VALUES ($1, 'example@example.com', $2)"
		_, err = c.db.Exec(query, guid, ip)
	} else {
		err = errh.Code01()
	}
	return err
}

func (c *Connection) RewriteIP(oldIp, newIp string) error {
	var err error
	if oldIp != "" && newIp != "" {
		query := "UPDATE Users SET ip = $1 WHERE ip = $2"
		_, err = c.db.Exec(query, newIp, oldIp)
	} else {
		err = errh.Code01()
	}
	return err
}

func (c *Connection) GetGUID(ip string) (string, error) {
	var guid string
	var err error
	if ip != "" {
		query := "SELECT guid FROM Users WHERE ip = $1"
		err = c.db.QueryRow(query, ip).Scan(&guid)
	} else {
		err = errh.Code01()
	}
	return guid, err
}

func (c *Connection) IsThereRefreshToken(guid, ip string) (bool, error) {
	var res int
	var err error
	if guid != "" && ip != "" {
		query := "SELECT COUNT(*) FROM Users WHERE guid = $1 AND ip = $2 AND refreshtoken IS NOT NULL"
		err = c.db.QueryRow(query, guid, ip).Scan(&res)
	} else {
		err = errh.Code01()
	}
	return res != 0, err
}
