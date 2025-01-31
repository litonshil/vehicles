package domain

import (
	"context"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"vehicles/types"
	"time"
)

type User struct {
	ID primitive.ObjectID `json:"id"`
	//Metadata Meta               `json:"meta"`
	IsAdmin bool `json:"is_admin"`
	Profile
}

type UserResp struct {
	ID          primitive.ObjectID `json:"id"`
	UserName    string             `json:"user_name"`
	FirstName   string             `json:"first_name"`
	LastName    string             `json:"last_name"`
	Email       string             `json:"email"`
	Phone       *string            `json:"phone"`
	ProfilePic  *string            `json:"profile_pic"`
	AppKey      string             `json:"app_key,omitempty"`
	RoleID      uint               `json:"role_id"`
	CreatedAt   time.Time          `json:"-"`
	UpdatedAt   time.Time          `json:"-"`
	LastLoginAt *time.Time         `json:"last_login_at"`
	FirstLogin  bool               `json:"first_login"`
}

func (u *User) Validate() error {
	return v.ValidateStruct(u,
		v.Field(&u.ID, v.Required),
	)
}

type Profile struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Meta struct {
	Method      string  `json:"method"`
	URI         string  `json:"uri"`
	ServiceName *string `json:"serviceName,omitempty"`
	AppKey      *string `json:"app-key,omitempty"`
	Profile
	Payload interface{} `json:"payload"`
}

type UserFilter struct {
	ID string `json:"id"`
}

type UserUseCase interface {
	CreateUser(ctx context.Context, req types.UserReq) error
	GetUsers(ctx context.Context, filter UserFilter) ([]User, error)
}

type UserRepo interface {
	CreateUser(req *User) error
	GetUsers(filter UserFilter) ([]User, error)
}
