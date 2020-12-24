package user

import (
	"github.com/mhdiiilham/oauth2-auth-server-implementation/pkg/password"
)

type manager struct {
	Repo Repository
}

// NewManager create new repository
func NewManager(r Repository) *manager {
	return &manager{
		Repo: r,
	}
}

// Register new user
func (s *manager) Register(user User) string {
	hashPassword := password.Hash(user.Password)
	user.Password = hashPassword
	return s.Repo.Register(user)
}

// FindOne user
func (s *manager) FindOne(email string) (*User, error) {
	return s.Repo.FindOne(email)
}
