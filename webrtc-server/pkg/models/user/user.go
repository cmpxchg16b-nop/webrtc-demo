package user

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/google/uuid"
)

// ErrUsernameDuplicated is returned when attempting to create a user with a username that already exists.
var ErrUsernameDuplicated = errors.New("username already exists")

// `User`s are immutable, no one should modify a `User` returned from a `UserManager`.
type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`

	// could be a data URL
	AvatarURL string `json:"avatar_url"`
	GithubId  string `json:"github_id"`
	DN42ASN   string `json:"dn42_asn"`
	GoogleId  string `json:"google_id"`
	IsBot     bool   `json:"is_bot"`
}

type UserCreationPayload struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type UserManager interface {
	// returns (theUser, created, error)
	LoadOrCreateNewUserByGithubId(ctx context.Context, githubId string, newUser User) (*User, bool, error)

	LoadOrCreateNewUserByDN42ASN(ctx context.Context, dn42ASN string, newUser User) (*User, bool, error)

	GetUserById(ctx context.Context, userId string) (*User, error)

	GetUserByUsername(ctx context.Context, username string) (*User, error)

	CreateUser(ctx context.Context, payload UserCreationPayload, isBot bool) (*User, error)
}

type InMemoryUserStore struct {
	Users           []User
	IndexById       map[string]int
	IndexByGithubId map[string]int
	IndexByUsername map[string]int
	IndexByDN42ASN  map[string]int
}

// Only call this after indexes were cloned, make sure you are modifying a memory region that
// is private to current goroutine
func (newStore *InMemoryUserStore) updateIndex(u *User, i int) {
	newStore.IndexById[u.Id] = i
	newStore.IndexByGithubId[u.GithubId] = i
	newStore.IndexByUsername[u.Username] = i
	newStore.IndexByDN42ASN[u.DN42ASN] = i
}

func (store *InMemoryUserStore) Clone() *InMemoryUserStore {
	// Each new Store would be private to current goroutine, so
	// there will NEVER be a concurrent write!
	newStore := new(InMemoryUserStore)
	newStore.IndexById = make(map[string]int)
	newStore.IndexByGithubId = make(map[string]int)
	newStore.IndexByUsername = make(map[string]int)
	newStore.IndexByDN42ASN = make(map[string]int)
	if store != nil {
		*newStore = *store
		newStore.Users = make([]User, len(store.Users))
		for i := range store.Users {
			u := store.Users[i]
			newStore.Users[i] = u
			newStore.updateIndex(&u, i)
		}
	}
	return newStore
}

func (store *InMemoryUserStore) AddUser(user User) (*InMemoryUserStore, error) {
	newStore := store.Clone()

	// NOTE: after Clone(), each thread modififies the clone of its own, not the same memory region

	if _, exists := newStore.IndexByUsername[user.Username]; exists {
		return nil, ErrUsernameDuplicated
	}

	numId := len(newStore.Users)
	newStore.updateIndex(&user, numId)
	newStore.Users = append(newStore.Users, user)
	return newStore, nil
}

type MemoryUserManager struct {
	store atomic.Pointer[InMemoryUserStore]
}

// Returns true means new user is created
func (memUserMngr *MemoryUserManager) doAddUserByGithubId(user User) (*User, bool, error) {
	for {
		oldStore := memUserMngr.store.Load()
		if oldStore != nil {
			if idx, hit := oldStore.IndexByGithubId[user.GithubId]; hit {
				return &oldStore.Users[idx], false, nil
			}
		}
		newStore, err := oldStore.AddUser(user)
		if err != nil {
			return nil, false, err
		}
		if memUserMngr.store.CompareAndSwap(oldStore, newStore) {
			return &user, true, nil
		}
	}
}

// Returns true means new user is created
func (memUserMngr *MemoryUserManager) doAddUserByDN42ASN(user User) (*User, bool, error) {
	for {
		oldStore := memUserMngr.store.Load()
		if oldStore != nil {
			if idx, hit := oldStore.IndexByDN42ASN[user.DN42ASN]; hit {
				return &oldStore.Users[idx], false, nil
			}
		}
		newStore, err := oldStore.AddUser(user)
		if err != nil {
			return nil, false, err
		}
		if memUserMngr.store.CompareAndSwap(oldStore, newStore) {
			return &user, true, nil
		}
	}
}

func (memUserMngr *MemoryUserManager) LoadOrCreateNewUserByGithubId(ctx context.Context, githubId string, newUser User) (*User, bool, error) {

	if id := newUser.Id; id == "" {
		newUser.Id = uuid.NewString()
	}

	u, accepted, err := memUserMngr.doAddUserByGithubId(newUser)
	if err != nil {
		return nil, false, err
	}
	return u, accepted, nil
}

func (memUserMngr *MemoryUserManager) LoadOrCreateNewUserByDN42ASN(ctx context.Context, dn42ASN string, newUser User) (*User, bool, error) {

	if id := newUser.Id; id == "" {
		newUser.Id = uuid.NewString()
	}

	u, accepted, err := memUserMngr.doAddUserByDN42ASN(newUser)
	if err != nil {
		return nil, false, err
	}
	return u, accepted, nil
}

func (memUserMngr *MemoryUserManager) GetUserById(ctx context.Context, userId string) (*User, error) {
	if store := memUserMngr.store.Load(); store != nil {
		if idx, hit := store.IndexById[userId]; hit {
			return &store.Users[idx], nil
		}
	}
	return nil, nil
}

func (memUserMngr *MemoryUserManager) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	if store := memUserMngr.store.Load(); store != nil {
		if idx, hit := store.IndexByUsername[username]; hit {
			return &store.Users[idx], nil
		}
	}
	return nil, nil
}

func (memUserMngr *MemoryUserManager) CreateUser(ctx context.Context, payload UserCreationPayload, isBot bool) (*User, error) {
	user := User{
		Id:          uuid.NewString(),
		Username:    payload.Username,
		DisplayName: payload.DisplayName,
		AvatarURL:   payload.AvatarURL,
		IsBot:       isBot,
	}

	for {
		oldStore := memUserMngr.store.Load()
		if oldStore != nil {
			if _, hit := oldStore.IndexByUsername[user.Username]; hit {
				return nil, ErrUsernameDuplicated
			}
		}
		newStore, err := oldStore.AddUser(user)
		if err != nil {
			return nil, err
		}
		if memUserMngr.store.CompareAndSwap(oldStore, newStore) {
			return &user, nil
		}
	}
}
