package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	Revision  int
	Users     []User
	IndexById map[string]int
}

func (store *InMemoryUserStore) Clone() *InMemoryUserStore {
	// Each new Store would be private to current goroutine, so
	// there will NEVER be a concurrent write!
	newStore := new(InMemoryUserStore)
	if store != nil {
		*newStore = *store
		newStore.Revision += 1
		newStore.Users = make([]User, len(store.Users))
		newStore.IndexById = make(map[string]int)
		for i := range store.Users {
			u := store.Users[i]
			newStore.Users[i] = *u.Clone()
			newStore.IndexById[u.Id] = i
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
func (memUserMngr *MemoryUserManager) doAddUser(user User) (*User, bool) {
	for {
		oldStore := memUserMngr.store.Load()
		if idx, hit := oldStore.IndexById[user.Id]; hit {
			return &oldStore.Users[idx], false
		}
		if memUserMngr.store.CompareAndSwap(oldStore, oldStore.AddUser(user)) {
			return &user, true
		}
	}
}

func (memUserMngr *MemoryUserManager) getIdFromGithubId(githubId string) string {
	// deterministic ID generation: if two goroutines are trying to register two users that has same github id,
	// they would be treated as the same user, one attempty of them will success, another one will failed, the failed goroutine
	// will fetch the already created user instead of override the previously created one.
	hashBin := sha256.Sum256([]byte(fmt.Sprintf("Github-%s", githubId)))
	return hex.EncodeToString(hashBin[:])
}

func (memUserMngr *MemoryUserManager) LoadOrCreateNewUserByGithubId(ctx context.Context, githubId string, newUser User) (User, bool, error) {

	newUser.Id = memUserMngr.getIdFromGithubId(githubId)
	if user, _ := memUserMngr.GetUserById(ctx, newUser.Id); user != nil {
		return *user, false, nil
	}

	// deterministic ID generation: if two goroutines are trying to register two users that has same github id,
	// they would be treated as the same user, one attempty of them will success, another one will failed, the failed goroutine
	// will fetch the already created user instead of override the previously created one.

	u, accepted := memUserMngr.doAddUser(newUser)
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
