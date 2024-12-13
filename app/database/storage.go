package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Connection struct {
	db *sql.DB
}

// Create a new connection to the database according to the env-values
// You might get an error if the conection is failed
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

// Delete all of clients in the table Clients
func (c *Connection) DeleteCliets() {
	query := "DELETE FROM Clients"
	_, err := c.db.Exec(query)
	if err != nil {
		panic(err)
	}
}

// Get a refresh token from the database by a specific id
func (c *Connection) GetRefreshToken(ip string) ([]byte, error) {
	var token []byte
	query := "SELECT refreshtoken FROM Clients WHERE ip = $1"
	err := c.db.QueryRow(query, ip).Scan(&token)
	return token, err
}

func (c *Connection) SaveRefreshToken(token []byte, guid string, onceagain *bool) error {
	query := "UPDATE Clients SET refreshtoken = $1 WHERE guid = $2"
	_, err := c.db.Exec(query, token, guid)
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code != "23505" {
			*onceagain = false
		}
	} else {
		*onceagain = false
	}
	return err
}

// Create a new clients only for tests
func (c *Connection) CreateNewClient(email, ip string) int {
	var id int
	query := "INSERT INTO Clients (email, ip) VALUES ($1, $2)"
	_, err := c.db.Exec(query, email, ip)
	if err != nil {
		panic(err)
	}
	query = "SELECT id FROM Clients WHERE email = $1"
	err = c.db.QueryRow(query, email).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}

func (c *Connection) CheckID(guid string) (bool, error) {
	var count int
	var err error
	if guid != "" {
		query := "SELECT COUNT(*) FROM Clients WHERE guid = $1"
		err = c.db.QueryRow(query, guid).Scan(&count)
	}
	return count == 1, err
}

// Check an ip by a specific id. If the ip by the id exists, you'll get true, if vice versa - false
func (c *Connection) CheckIP(id int, ip string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM Clients WHERE id = $1 AND ip = $2"
	err := c.db.QueryRow(query, id, ip).Scan(&count)
	return count == 1, err
}

// Get an email from the database
func (c *Connection) SelectEmail(ip string) (string, error) {
	var email string
	query := "SELECT email FROM Clients WHERE ip = $1"
	err := c.db.QueryRow(query, ip).Scan(&email)
	return email, err
}

// Create a client for getting better user experience (only for a real using)
func (c *Connection) CreateMokData() {
	query := "INSERT INTO Clients (email, ip) VALUES ('example@example.com', '213.136.11.188')"
	_, err := c.db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (c *Connection) RewriteIP(oldIp, newIp string) error {
	query := "UPDATE Clients SET ip = %1 WHERE ip = %2"
	_, err := c.db.Exec(query, newIp, oldIp)
	return err
}

func (c *Connection) GetGUID(ip string) (string, error) {
	var guid string
	query := "SELECT guid FROM Clients WHERE ip = %1"
	err := c.db.QueryRow(query, ip).Scan(&guid)
	return guid, err
}
