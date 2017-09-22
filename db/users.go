package db

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User represents the db schema of a user
type User struct {
	ID        int64
	UUID      string
	Email     string
	Nickname  string
	About     string
	Address   StringSlice
	ZIP       string
	City      string
	Country   string
	Activated bool
	AuthToken StringSlice
}

// LoadUserByUUID loads a user by UUID from the database
func (context *APIContext) LoadUserByUUID(uuid string) (User, error) {
	user := User{}
	if len(uuid) == 0 {
		return user, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, uuid, nickname, about, email, address, zip, city, country, activated FROM users WHERE uuid = $1", uuid).
		Scan(&user.ID, &user.UUID, &user.Nickname, &user.About, &user.Email, &user.Address, &user.ZIP, &user.City, &user.Country, &user.Activated)
	return user, err
}

// GetUserByUUID returns a user by UUID from the cache
func (context *APIContext) GetUserByUUID(uuid string) (User, error) {
	user := User{}
	usersCache, err := usersCache.Value(uuid, context)
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
	err := context.QueryRow("SELECT id, uuid, nickname, about, email, address, zip, city, country, activated, authtoken, password FROM users WHERE nickname = $1", name).
		Scan(&user.ID, &user.UUID, &user.Nickname, &user.About, &user.Email, &user.Address, &user.ZIP, &user.City, &user.Country, &user.Activated, &user.AuthToken, &hashedPassword)
	if err != nil {
		return User{}, errors.New("Invalid username or password")
	}

	//FIXME: cryptpepper
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+"cryptpepper"))
	if err != nil {
		return User{}, errors.New("Invalid username or password")
	}

	return user, nil
}

// GetUserByEmail loads a user by email from the database
func (context *APIContext) GetUserByEmail(email string) (User, error) {
	user := User{}
	err := context.QueryRow("SELECT id, uuid, nickname, about, email, address, zip, city, country, activated, authtoken FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.UUID, &user.Nickname, &user.About, &user.Email, &user.Address, &user.ZIP, &user.City, &user.Country, &user.Activated, &user.AuthToken)
	if err != nil {
		return User{}, errors.New("Invalid email address")
	}

	return user, nil
}

// GetUserByAccessToken loads a user by accesstoken from the database
func (context *APIContext) GetUserByAccessToken(token string) (interface{}, error) {
	user := User{}
	err := context.QueryRow("SELECT id, uuid, nickname, about, email, address, zip, city, country, activated, authtoken FROM users WHERE $1 = ANY(authtoken)", token).
		Scan(&user.ID, &user.UUID, &user.Nickname, &user.About, &user.Email, &user.Address, &user.ZIP, &user.City, &user.Country, &user.Activated, &user.AuthToken)

	return user, err
}

// LoadAllUsers loads all users from the database
func (context *APIContext) LoadAllUsers() ([]User, error) {
	users := []User{}

	rows, err := context.Query("SELECT id, uuid, nickname, about, email, address, zip, city, country, activated FROM users")
	if err != nil {
		return users, err
	}

	defer rows.Close()
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.UUID, &user.Nickname, &user.About, &user.Email, &user.Address, &user.ZIP, &user.City, &user.Country, &user.Activated)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, err
}

// Update a user in the database
func (user *User) Update(context *APIContext) error {
	_, err := context.Exec("UPDATE users SET about = $1, email = $2, address = $3, authtoken = $4 WHERE id = $5",
		user.About, user.Email, user.Address, user.AuthToken, user.ID)
	if err != nil {
		panic(err)
	}

	usersCache.Delete(user.UUID)
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
	usersCache.Delete(user.UUID)
	return err
}

// Save a user to the database
func (user *User) Save(context *APIContext) error {
	uuid, err := UUID()
	if err != nil {
		return err
	}

	user.UUID, _ = UUID()
	user.AuthToken = StringSlice{uuid}
	err = context.QueryRow("INSERT INTO users (uuid, nickname, password, about, address, zip, city, country, email, authtoken) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		user.UUID, user.Nickname, uuid, user.About, user.Address, user.ZIP, user.City, user.Country, user.Email, user.AuthToken).Scan(&user.ID)
	usersCache.Delete(user.UUID)
	return err
}
