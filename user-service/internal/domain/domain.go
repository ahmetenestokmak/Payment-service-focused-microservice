package domain

import(
	"context"
)

type User struct {
	ID       	string `json:"id" db:"id"`
	FirstName  	string `json:"first_name" db:"first_name"`
	LastName 	string `json:"last_name" db:"last_name"`
}
type UserSalt struct {
	FirstName  	string `json:"first_name" db:"first_name"`
	LastName   	string `json:"last_name" db:"last_name"`
}

type UserRepository interface{
	Save(ctx context.Context, user *User) error
	GetProfile(ctx context.Context, id string) (*UserSalt, error)
}

type UserUsecase interface{
	Save(ctx context.Context, user *User) error
	GetProfile(ctx context.Context, id string) (*UserSalt, error)
} 