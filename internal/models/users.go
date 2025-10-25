package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UsersModel struct {
	DB *sql.DB
}

func (model *UsersModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = model.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (model *UsersModel) Authenticate(email, password string) (int, error) {
	findEmailStmt := `SELECT id, hashed_password FROM users WHERE email = ?`

	row := model.DB.QueryRow(findEmailStmt, email)

	var usrId int
	var hashedPassword []byte

	err := row.Scan(&usrId, &hashedPassword)
	if err != nil {
		fmt.Print("err no recs")
		return 0, ErrNoRecords
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return 0, ErrInvalidCredentials
	}

	return usrId, nil
}

func (model *UsersModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT 1 FROM users where id = ?)`

	err := model.DB.QueryRow(stmt, id).Scan(&exists)
	if err != nil {
		return false, ErrNoRecords
	}

	return exists, nil
}

func (model *UsersModel) Get(id int) (User, error) {
	stmt := `SELECT id, name, email, created FROM users WHERE id = ?`

	info := User{}
	err := model.DB.QueryRow(stmt, id).Scan(&info.ID, &info.Name, &info.Email, &info.Created)
	if err != nil {
		return User{}, ErrNoRecords
	}

	return info, nil

}
