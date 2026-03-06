package user

import (
	"context"
	"sync/atomic"

	"github.com/google/uuid"
)

// `User`s are immutable, no one should modify a `User` returned from a `UserManager`.
type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	GithubId    string `json:"github_id"`
}

type UserManager interface {
	// returns (theUser, created, error)
	LoadOrCreateNewUserByGithubId(ctx context.Context, githubId string, newUser User) (User, bool, error)

	GetUserById(ctx context.Context, userId string) (*User, error)
}

type InMemoryUserStore struct {
	Users           []User
	IndexById       map[string]int
	IndexByGithubId map[string]int
}

func (store *InMemoryUserStore) Clone() *InMemoryUserStore {
	// Each new Store would be private to current goroutine, so
	// there will NEVER be a concurrent write!
	newStore := new(InMemoryUserStore)
	newStore.IndexById = make(map[string]int)
	newStore.IndexByGithubId = make(map[string]int)
	if store != nil {
		*newStore = *store
		newStore.Users = make([]User, len(store.Users))
		for i := range store.Users {
			u := store.Users[i]
			newStore.Users[i] = u
			newStore.IndexById[u.Id] = i
			if ghId := u.GithubId; ghId != "" {
				newStore.IndexByGithubId[u.GithubId] = i
			}
		}
	}
	return newStore
}

func (store *InMemoryUserStore) AddUser(user User) *InMemoryUserStore {
	newStore := store.Clone()

	// NOTE: after Clone(), each thread modififies the clone of its own, not the same memory region

	newStore.IndexById[user.Id] = len(newStore.Users)
	newStore.Users = append(newStore.Users, user)
	return newStore
}

type MemoryUserManager struct {
	store atomic.Pointer[InMemoryUserStore]
}

// Returns true means new user is created
func (memUserMngr *MemoryUserManager) doAddUserByGithubId(user User) (*User, bool) {
	for {
		oldStore := memUserMngr.store.Load()
		if oldStore != nil {
			if idx, hit := oldStore.IndexByGithubId[user.GithubId]; hit {
				return &oldStore.Users[idx], false
			}
		}
		if memUserMngr.store.CompareAndSwap(oldStore, oldStore.AddUser(user)) {
			return &user, true
		}
	}
}

func (memUserMngr *MemoryUserManager) LoadOrCreateNewUserByGithubId(ctx context.Context, githubId string, newUser User) (User, bool, error) {

	if id := newUser.Id; id == "" {
		newUser.Id = uuid.NewString()
	}

	u, accepted := memUserMngr.doAddUserByGithubId(newUser)
	return *u, accepted, nil
}

func (memUserMngr *MemoryUserManager) GetUserById(ctx context.Context, userId string) (*User, error) {
	if store := memUserMngr.store.Load(); store != nil {
		if idx, hit := store.IndexById[userId]; hit {
			return &store.Users[idx], nil
		}
	}
	return nil, nil
}
