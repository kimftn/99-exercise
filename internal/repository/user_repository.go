package repository

import (
	"property-api/internal/domain"
	"property-api/internal/store"
)

type UserRepository interface {
	Create(user domain.User) domain.User
	GetAll() []domain.User
	GetByID(id int) (domain.User, bool)
}

type InMemoryUserRepository struct {
	store *store.Store
}

func NewInMemoryUserRepository(store *store.Store) *InMemoryUserRepository {
	return &InMemoryUserRepository{store: store}
}

func (r *InMemoryUserRepository) Create(user domain.User) domain.User {
	user.ID = len(r.store.Users) + 1
	r.store.Users = append(r.store.Users, user)
	return user
}

func (r *InMemoryUserRepository) GetAll() []domain.User {
	return r.store.Users
}

func (r *InMemoryUserRepository) GetByID(id int) (domain.User, bool) {
	for _, user := range r.store.Users {
		if user.ID == id {
			return user, true
		}
	}

	return domain.User{}, false
}
