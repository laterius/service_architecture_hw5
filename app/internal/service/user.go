package service

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/repo"
	"github.com/laterius/service_architecture_hw3/app/pkg/nullable"
	"github.com/laterius/service_architecture_hw3/app/pkg/types"
)

type UserData struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type User struct {
	Id int64 `json:"id"`
	UserData
}

func (u *User) FromDomain(d *domain.User) *User {
	u.Id = int64(d.Id)
	u.Username = d.Username
	u.FirstName = d.FirstName
	u.LastName = d.LastName
	u.Email = d.Email
	u.Phone = d.Phone

	return u
}

type UserCreate UserData

func (u UserCreate) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(2, 32)),
		validation.Field(&u.FirstName, validation.Required, validation.Length(1, 32)),
		validation.Field(&u.LastName, validation.Required, validation.Length(1, 32)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Phone, validation.Required, is.E164),
	)
}

func (u UserCreate) ToDomain() *domain.User {
	return &domain.User{
		Id:        0,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
	}
}

type UserUpdate UserData

func (u UserUpdate) ToDomain() *domain.User {
	return &domain.User{
		Id:        0,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
	}
}

type UserPartialUpdate struct {
	Username  nullable.String `json:"username"`
	FirstName nullable.String `json:"firstName"`
	LastName  nullable.String `json:"lastName"`
	Email     nullable.String `json:"email"`
	Phone     nullable.String `json:"phone"`
}

func (pu UserPartialUpdate) ToDomain() *domain.UserPartialData {
	d := types.NewKv()
	if pu.Username.Set {
		d.Set("username", pu.Username.Value)
	}
	if pu.FirstName.Set {
		d.Set("firstName", pu.FirstName.Value)
	}
	if pu.LastName.Set {
		d.Set("lastName", pu.LastName.Value)
	}
	if pu.Email.Set {
		d.Set("email", pu.Email.Value)
	}
	if pu.Phone.Set {
		d.Set("phone", pu.Phone.Value)
	}

	return d
}

type UserReader interface {
	Get(domain.UserId) (*domain.User, error)
}

type UserCreator interface {
	Create(*UserCreate) (*domain.User, error)
}

type UserUpdater interface {
	Update(domain.UserId, *UserUpdate) (*domain.User, error)
}

type UserPartialUpdater interface {
	PartialUpdate(domain.UserId, *domain.UserPartialData) (*domain.User, error)
}

type UserDeleter interface {
	Delete(domain.UserId) error
}

type UserService interface {
	UserReader
	UserCreator
	UserUpdater
	UserPartialUpdater
	UserDeleter
}

type userService struct {
	reader         repo.UserReader
	observer       repo.UserObserver
	creator        repo.UserCreator
	updater        repo.UserUpdater
	partialUpdater repo.UserPartialUpdater
	deleter        repo.UserDeleter
}

func NewUserService(repo repo.UserRepo) UserService {
	return &userService{
		reader:         repo,
		observer:       repo,
		creator:        repo,
		updater:        repo,
		partialUpdater: repo,
		deleter:        repo,
	}
}

func (s *userService) Get(id domain.UserId) (*domain.User, error) {
	err := id.Validate()
	if err != nil {
		return nil, err
	}

	user, err := s.reader.Get(id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, domain.ErrUserNotFound
		}
	}
	return user, err
}

func (s *userService) Create(req *UserCreate) (*domain.User, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	return s.creator.Create(req.ToDomain())
}

func (s *userService) Update(id domain.UserId, req *UserUpdate) (*domain.User, error) {
	exists, err := s.observer.Exists(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return s.updater.Update(id, req.ToDomain())
}

func (s *userService) PartialUpdate(id domain.UserId, data *domain.UserPartialData) (*domain.User, error) {
	return s.partialUpdater.PartialUpdate(id, data)
}

func (s *userService) Delete(id domain.UserId) error {
	return s.deleter.Delete(id)
}
