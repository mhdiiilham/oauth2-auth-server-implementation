package user

// Writer interface
type Writer interface {
	Register(user User) string
}

// Reader interface
type Reader interface {
	FindOne(email string) (*User, error)
}

// Repository interface
type Repository interface {
	Writer
	Reader
}

// Manager interface
type Manager interface {
	Repository
}
