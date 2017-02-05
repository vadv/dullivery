package store

import (
	"fmt"
	"strings"
)

// без локов, только хардор

type Users struct {
	List []*User `json:"list"`
}

type User struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

func (s *Storage) AddUser(user User) error {
	if user.Email == "" {
		fmt.Errorf("пустой email")
	}
	if _, found := s.GetUser(user.Email); found {
		s.DelUser(user.Email)
	}
	s.Users.List = append(s.Users.List, &user)
	return s.save()
}

func (s *Storage) DelUser(email string) {
	lowerEmail := strings.ToLower(email)
	result := &Users{List: make([]*User, 0)}
	for _, user := range s.Users.List {
		if strings.ToLower(user.Email) != lowerEmail {
			result.List = append(result.List, user)
		}
	}
	s.Users = result
	s.save()
}

func (s *Storage) GetUser(email string) (User, bool) {
	lowerEmail := strings.ToLower(email)
	for _, user := range s.Users.List {
		if strings.ToLower(user.Email) == lowerEmail {
			return User{Email: user.Email, Hash: user.Hash}, true
		}
	}
	return User{}, false
}
