package users

import (
	"errors"
	"fmt"
	"net/mail"
)

var ErrNoResultFound = errors.New("no result found")

type User struct {
	FirtName string
	LastName string
	Email    mail.Address
}

type Manager struct {
	users []User
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) AddUser(firstName string, lastName string, email string) error {
	parsedAddress, err := mail.ParseAddress(email) //Parse from string to mail.Adress
	if err != nil {
		return fmt.Errorf("invalid email: %s", email)
	}
	if firstName == "" {
		return fmt.Errorf("invalid firstname: %s", firstName)
	}
	if lastName == "" {
		return fmt.Errorf("invalid lastname: %s", lastName)
	}

	existingUser, err := m.GetUserByName(firstName, lastName)
	if err != nil && !errors.Is(err, ErrNoResultFound) {
		return fmt.Errorf("error checking if user is already present: %v", err)
	}

	if existingUser != nil {
		return errors.New("user with this name already exist")
	}

	newUser := User{ //add user to struct
		FirtName: firstName,
		LastName: lastName,
		Email:    *parsedAddress,
	}

	m.users = append(m.users, newUser) //apend user to manager users slice
	return nil                         // return nil if we don't have mistakes
}

func (m *Manager) GetUserByName(first, last string) (*User, error) {
	for i, user := range m.users {
		if user.FirtName == first && user.LastName == last {
			result := m.users[i]
			return &result, nil
		}
	}

	return nil, ErrNoResultFound
}
