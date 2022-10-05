package repo

import "github.com/laterius/service_architecture_hw3/app/internal/domain"

type UserReader interface {
	Get(domain.UserId) (*domain.User, error)
}

type UserObserver interface {
	Exists(domain.UserId) (bool, error)
}

type UserCreator interface {
	Create(*domain.User) (*domain.User, error)
}

type UserUpdater interface {
	Update(domain.UserId, *domain.User) (*domain.User, error)
}

type UserPartialUpdater interface {
	PartialUpdate(domain.UserId, *domain.UserPartialData) (*domain.User, error)
}

type UserDeleter interface {
	Delete(domain.UserId) error
}

type UserRepo interface {
	UserReader
	UserObserver
	UserCreator
	UserUpdater
	UserPartialUpdater
	UserDeleter
}
