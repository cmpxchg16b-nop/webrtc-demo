package user

import (
	"context"
	"sync/atomic"

	"github.com/google/uuid"
)

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	GithubId    string `json:"github_id"`
}

func (u *User) Clone() *User {
	newUser := new(User)
	*newUser = *u
	return newUser
}

type UserManager interface {
	// returns (theUser, created, error)
	LoadOrCreateNewUserByGithubId(ctx context.Context, githubId string, newUser User) (User, bool, error)

	GetUserById(ctx context.Context, userId string) (*User, error)
}

type InMemoryUserStore struct {
	Revision int
	Users    []User
}

func (store *InMemoryUserStore) Clone() *InMemoryUserStore {
	newStore := new(InMemoryUserStore)
	if store != nil {
		*newStore = *store
		newStore.Revision += 1
		newStore.Users = make([]User, len(store.Users))
		for i := range store.Users {
			newStore.Users[i] = *store.Users[i].Clone()
		}
	}
	return newStore
}

func (store *InMemoryUserStore) AddUser(user User) *InMemoryUserStore {
	newStore := store.Clone()
	newStore.Users = append(newStore.Users, user)
	return newStore
}

type MemoryUserManager struct {
	store atomic.Pointer[InMemoryUserStore]
}

func (memUserMngr *MemoryUserManager) doAddUser(user User) {
	for {
		oldStore := memUserMngr.store.Load()
		if memUserMngr.store.CompareAndSwap(oldStore, oldStore.AddUser(user)) {
			break
		}
	}
}

func (memUserMngr *MemoryUserManager) doCreateUser(user User) User {
	user.Id = uuid.NewString()
	memUserMngr.doAddUser(user)
	return user
}

func (memUserMngr *MemoryUserManager) LoadOrCreateNewUserByGithubId(ctx context.Context, githubId string, newUser User) (User, bool, error) {

	if user, _ := memUserMngr.GetUserGithubId(ctx, githubId); user != nil {
		return *user, false, nil
	}

	return memUserMngr.doCreateUser(newUser), true, nil
}

func (memUserMngr *MemoryUserManager) GetUserById(ctx context.Context, userId string) (*User, error) {
	if store := memUserMngr.store.Load(); store != nil {
		for _, u := range store.Users {
			if u.Id == userId {
				return &u, nil
			}
		}
	}
	return nil, nil
}

func (memUserMngr *MemoryUserManager) GetUserGithubId(ctx context.Context, githubId string) (*User, error) {
	if store := memUserMngr.store.Load(); store != nil {
		for _, u := range store.Users {
			if u.GithubId == githubId {
				return &u, nil
			}
		}
	}
	return nil, nil
}
