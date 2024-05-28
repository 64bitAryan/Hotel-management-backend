package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minFirstNameLength = 2
	minLastNameLength  = 3
	minPasswordLength  = 7
)

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (u *UpdateUserParams) ToBISON() bson.M {
	m := bson.M{}
	if len(u.FirstName) > 0 {
		m["firstName"] = u.FirstName
	}
	if len(u.LastName) > 0 {
		m["lastName"] = u.LastName
	}
	return m
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (c CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(c.FirstName) < minFirstNameLength {
		errors["firstName"] = fmt.Sprintf("first name should be atleast %d character", minFirstNameLength)
	}
	if len(c.LastName) < minFirstNameLength {
		errors["lastName"] = fmt.Sprintf("last name should be atleast %d character", minLastNameLength)
	}
	if len(c.Password) < minPasswordLength {
		errors["password"] = fmt.Sprintf("password should be atleast %d character", minPasswordLength)
	}
	if isEmailValid(c.Email) {
		errors["email"] = fmt.Sprintf("invalid email")
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(e)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpv, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpv),
	}, nil
}

func IsValidPassword(EncryptedPassword, password string) error {

	err := bcrypt.CompareHashAndPassword([]byte(EncryptedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}
	return nil
}
