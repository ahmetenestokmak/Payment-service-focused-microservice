package domain

import "context"

type User struct {
	ID       string `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Role     string `json:"role" db:"role"`
}
type UserSalt struct {
	ID       string `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Role     string `json:"role" db:"role"`
	Token string `json:"token" db:"token"`
}

type UserLogin struct {
	Email     string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type AuthRepository interface {
	Login(ctx context.Context,  email string, password string) (*UserSalt, error)
	Update(ctx context.Context, id string, email string, password string, role string) error
	Register(ctx context.Context,  email string, password string, role string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*UserSalt, error) 
}

type AuthUsecase interface {
	Login(ctx context.Context, email string, password string) (*UserSalt, error)
	Register(ctx context.Context, email string, password string, role string) (*User, error)
	Update(ctx context.Context, id string, email string, password string, role string) error
}
