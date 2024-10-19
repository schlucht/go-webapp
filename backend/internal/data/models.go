package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User:  User{},
		Token: Token{},
	}
}

type Models struct {
	User  User
	Token Token
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpatedAt  time.Time `json:"updated_at"`
	Token     Token     `json:"token"`
}

func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var urs []*User
	query := `select id, email, first_name, last_name, password, created_at, updated_at from users order by last_name`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var us User
		err := rows.Scan(
			&us.ID,
			&us.Email,
			&us.FirstName,
			&us.LastName,
			&us.Password,
			&us.CreatedAt,
			&us.UpatedAt,
		)
		if err != nil {
			return nil, err
		}
		urs = append(urs, &us)
	}
	return urs, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, created_at, updated_at from users where email = $1`
	var us User
	row := db.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&us.ID,
		&us.Email,
		&us.FirstName,
		&us.LastName,
		&us.Password,
		&us.CreatedAt,
		&us.UpatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &us, nil
}

func (u *User) GetById(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, created_at, updated_at from users where id = $1`
	var us User
	row := db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&us.ID,
		&us.Email,
		&us.FirstName,
		&us.LastName,
		&us.Password,
		&us.CreatedAt,
		&us.UpatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &us, nil
}

func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		updated_at = $4
		where id = $5
	`
	_, err := db.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		delete from users		
		where id = $1
	`
	_, err := db.ExecContext(ctx, stmt,
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	pw := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pw, 12)
	if err != nil {
		return 0, err
	}
	stmt := `
		insert into users (
			email,
			first_name,
			last_name,
			password
		) values (
		 $1, $2, $3, $4
		) returning id
	`
	var id int
	res := db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
	)
	err = res.Scan(
		&id,
	)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	pw := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pw, 12)
	if err != nil {
		return err
	}
	stmt := `
		password = $1
		where id = $5
	`
	_, err = db.ExecContext(ctx, stmt,
		hashedPassword,
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
func (u *User) ResetPassword2(password string) error {
	pw := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pw, 12)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	err = u.Update()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type Token struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	TokenHash []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpatedAt  time.Time `json:"ubdated_at"`
	Expiry    time.Time `json:"expiry"`
}

func (t *Token) GetByToken(plainText string) (*Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, user_id, email, token, token_hash, created_at, updated_at, expiry
				from tokens where token = $1
				`
	var token Token
	row := db.QueryRowContext(ctx, query, plainText)
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.Email,
		&token.Token,
		&token.TokenHash,
		&token.CreatedAt,
		&token.UpatedAt,
		&token.Expiry,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
func (t *Token) GetUserForToken(token Token) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, first_name, last_name, email, password, created_at, updated_at
				from users where id = $1
				`
	var user User
	row := db.QueryRowContext(ctx, query, token.UserID)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (t *Token) GenerateToken(userID int, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Token = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Token))
	token.TokenHash = hash[:]

	return token, nil
}

func (t *Token) AuthenticateToken(r *http.Request) (*User, error) {
	authoridationHeader := r.Header.Get("Authorization")
	if authoridationHeader == "" {
		return nil, errors.New("no authorization header recevied")
	}
	headerParts := strings.Split(authoridationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no valid authorization header recevied")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("token wrong size")
	}

	tkn, err := t.GetByToken(token)
	if err != nil {
		return nil, errors.New("no matching token now")
	}
	if tkn.Expiry.Before(time.Now()) {
		return nil, errors.New("expired token")
	}

	user, err := t.GetUserForToken(*tkn)
	if err != nil {
		return nil, errors.New("no user found")
	}
	return user, nil
}

func (t *Token) Insert(token Token, u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from tokens where user_id = $1`
	_, err := db.ExecContext(ctx, stmt, token.UserID)
	if err != nil {
		return err
	}
	token.Email = u.Email

	stmt = `insert into tokens(user_id, email, token, token_hash, created_at, updated_at, expiry)
			values ($1,$2,$3,$4,$5,$6,$7)
			`
	_, err = db.ExecContext(ctx, stmt,
		u.ID,
		token.Email,
		token.Token,
		token.TokenHash,
		time.Now(),
		time.Now(),
		token.Expiry,
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) DeleteByToken(plainText string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from tokens where token = $1`
	_, err := db.ExecContext(ctx, stmt, plainText)
	if err != nil {
		return err
	}
	return nil
}

func (t *Token) ValidToken(plainText string) (bool, error) {
	token, err := t.GetByToken(plainText)
	if err != nil {
		return false, errors.New("no matching token found")
	}
	_, err = t.GetUserForToken(*token)
	if err != nil {
		return false, errors.New("no matching user found")
	}
	if token.Expiry.Before(time.Now()) {
		return false, errors.New("expired token")
	}

	return true, nil
}
