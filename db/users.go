package db

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User represents the db schema of a user
type User struct {
	ID        int64
	Email     string
	Nickname  string
	About     string
	Activated bool
	AuthToken StringSlice
}

// LoadUserByID loads a user by ID from the database
func (context *APIContext) LoadUserByID(id int64) (User, error) {
	user := User{}
	if id < 1 {
		return user, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, nickname, about, email, activated FROM users WHERE id = $1", id).Scan(&user.ID, &user.Nickname, &user.About, &user.Email, &user.Activated)
	return user, err
}

// GetUserByID returns a user by ID from the cache
func (context *APIContext) GetUserByID(id int64) (User, error) {
	user := User{}
	usersCache, err := usersCache.Value(id, context)
	if err != nil {
		return user, err
	}

	user = *usersCache.Data().(*User)
	return user, nil
}

// GetUserByNameAndPassword loads a user by name & password from the database
func (context *APIContext) GetUserByNameAndPassword(name, password string) (User, error) {
	user := User{}
	hashedPassword := ""
	err := context.QueryRow("SELECT id, nickname, about, email, activated, authtoken, password FROM users WHERE username = $1", name).Scan(&user.ID, &user.Nickname, &user.About, &user.Email, &user.Activated, &user.AuthToken, &hashedPassword)
	if err != nil {
		return User{}, errors.New("Invalid username or password")
	}

	//FIXME: cryptpepper
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+"cryptpepper"))
	if err != nil {
		return User{}, errors.New("Invalid username or password")
	}

	return user, err
}

// GetUserByAccessToken loads a user by accesstoken from the database
func (context *APIContext) GetUserByAccessToken(token string) (interface{}, error) {
	user := User{}
	err := context.QueryRow("SELECT id, nickname, about, email, activated, authtoken FROM users WHERE $1 = ANY(authtoken)", token).Scan(&user.ID, &user.Nickname, &user.About, &user.Email, &user.Activated, &user.AuthToken)

	return user, err
}

// LoadAllUsers loads all users from the database
func (context *APIContext) LoadAllUsers() ([]User, error) {
	users := []User{}

	rows, err := context.Query("SELECT id, username, about, email, activated FROM users")
	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.Nickname, &user.About, &user.Email, &user.Activated)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, err
}

// Update a user in the database
func (user *User) Update(context *APIContext) error {
	_, err := context.Exec("UPDATE users SET about = $1, email = $2, authtoken = $3 WHERE id = $4", user.About, user.Email, user.AuthToken, user.ID)
	if err != nil {
		panic(err)
	}

	usersCache.Delete(user.ID)
	return err
}

// UpdatePassword sets a new user password in the database
func (user *User) UpdatePassword(context *APIContext, password string) error {
	//FIXME: cryptpepper
	hash, err := bcrypt.GenerateFromPassword([]byte(password+"cryptpepper"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = context.Exec("UPDATE users SET password = $1, activated = true WHERE id = $2", string(hash), user.ID)
	usersCache.Delete(user.ID)
	return err
}

// Save a user to the database
func (user *User) Save(context *APIContext) error {
	uuid, err := UUID()
	if err != nil {
		return err
	}

	user.AuthToken = StringSlice{uuid}
	err = context.QueryRow("INSERT INTO users (nickname, password, about, email, authtoken) VALUES ($5, $4, $2, $1, $3) RETURNING id", user.Email, user.About, user.AuthToken, uuid, user.Nickname).Scan(&user.ID)
	usersCache.Delete(user.ID)
	return err
}
