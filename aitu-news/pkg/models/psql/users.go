package psql

import (
	"aitu.edu.kz/aitu-news/pkg/models"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, lastname, email, password string) (int, error) {
	stmt := `INSERT INTO users(name, lastname, email, password, role_id) VALUES ($1, $2, $3, $4, 2) RETURNING id`
	hashedPass, err := HashPassword(password)
	if err != nil {
		return 0, err
	}
	insertedId := 0
	err = m.DB.QueryRow(stmt, name, lastname, email, hashedPass).Scan(&insertedId)
	if err != nil {
		return 0, err
	}
	return int(insertedId), nil
}

func (m *UserModel) Get(email, password string) (*models.User, error) {
	stmt := `SELECT * FROM users WHERE email=$1`
	u := &models.User{}
	err := m.DB.QueryRow(stmt, email).Scan(&u.ID, &u.Name, &u.Lastname, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	if CheckPasswordHash(password, u.Password) {
		return u, nil
	}
	return nil, err
}

func (m *UserModel) GetById(id int) (*models.User, error) {
	stmt := `SELECT * FROM users WHERE id=$1`
	u := &models.User{}
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Lastname, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
