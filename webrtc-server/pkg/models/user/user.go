package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync/atomic"
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
	Index    map[string]int
}

func (store *InMemoryUserStore) Clone() *InMemoryUserStore {
	newStore := new(InMemoryUserStore)
	if store != nil {
		*newStore = *store
		newStore.Revision += 1
		newStore.Users = make([]User, len(store.Users))
		newStore.Index = make(map[string]int)
		for i := range store.Users {
			u := store.Users[i]
			newStore.Users[i] = *u.Clone()
			newStore.Index[u.Id] = i
		}
	}
	return newStore
}

func (store *InMemoryUserStore) AddUser(user User) *InMemoryUserStore {
	newStore := store.Clone()

	// NOTE: each thread modififies the clone, not the same memory region
	newUsers := append(newStore.Users, user)
	newStore.Index[user.Id] = len(newUsers)
	newStore.Users = newUsers
	return newStore
}

type MemoryUserManager struct {
	store atomic.Pointer[InMemoryUserStore]
}

// Returns true means new user is created
func (memUserMngr *MemoryUserManager) doAddUser(user User) (*User, bool) {
	for {
		oldStore := memUserMngr.store.Load()
		if idx, hit := oldStore.Index[user.Id]; hit {
			return &oldStore.Users[idx], false
		}
		if memUserMngr.store.CompareAndSwap(oldStore, oldStore.AddUser(user)) {
			return &user, true
		}
	}
}

func (memUserMngr *MemoryUserManager) LoadOrCreateNewUserByGithubId(ctx context.Context, githubId string, newUser User) (User, bool, error) {

	if user, _ := memUserMngr.GetUserGithubId(ctx, githubId); user != nil {
		return *user, false, nil
	}

	hashedId := sha256.Sum256([]byte(fmt.Sprintf("Github-%s", githubId)))
	newUser.Id = string(hashedId[:])
	u, accepted := memUserMngr.doAddUser(newUser)
	return *u, accepted, nil
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
