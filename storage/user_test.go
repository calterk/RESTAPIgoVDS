package storage

import (
	"app/entity"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	st = &UserStorage{
		userMap: make(map[uint64]entity.User),
		next:    1,
		mutx:    sync.RWMutex{},
	}
}

func TestCreateUser(t *testing.T) {
	setup()
	user := entity.User{
		Login: "testuser",
		Pass:  "password",
		Name:  "Test User",
	}
	createdUser := CreateUser(user)

	assert.Equal(t, uint64(1), createdUser.Uid)
	assert.Equal(t, user.Login, createdUser.Login)
	assert.Equal(t, user.Pass, createdUser.Pass)
	assert.Equal(t, user.Name, createdUser.Name)
}

func TestUpdateUser(t *testing.T) {
	setup()
	user := entity.User{
		Login: "testuser",
		Pass:  "password",
		Name:  "Test User",
	}
	createdUser := CreateUser(user)
	createdUser.Name = "Updated User"

	err := UpdateUser(createdUser)
	assert.NoError(t, err)

	updatedUser, err := GetUserByID(createdUser.Uid)
	assert.NoError(t, err)
	assert.Equal(t, "Updated User", updatedUser.Name)
}

func TestUpdateUserNotFound(t *testing.T) {
	setup()
	user := entity.User{
		Uid:   1,
		Login: "nonexistentuser",
		Pass:  "password",
		Name:  "Nonexistent User",
	}

	err := UpdateUser(user)
	assert.Error(t, err)
	assert.Equal(t, "User не найден", err.Error())
}

func TestDeleteUser(t *testing.T) {
	setup()
	user := entity.User{
		Login: "testuser",
		Pass:  "password",
		Name:  "Test User",
	}
	createdUser := CreateUser(user)

	err := DeleteUser(createdUser.Uid)
	assert.NoError(t, err)

	_, err = GetUserByID(createdUser.Uid)
	assert.Error(t, err)
	assert.Equal(t, "User не найден", err.Error())
}

func TestDeleteUserNotFound(t *testing.T) {
	setup()
	err := DeleteUser(1)
	assert.Error(t, err)
	assert.Equal(t, "User не найден", err.Error())
}

func TestGetAllUsers(t *testing.T) {
	setup()
	user1 := entity.User{
		Login: "testuser1",
		Pass:  "password1",
		Name:  "Test User 1",
	}
	user2 := entity.User{
		Login: "testuser2",
		Pass:  "password2",
		Name:  "Test User 2",
	}
	CreateUser(user1)
	CreateUser(user2)

	users, err := GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
}

func TestGetUserByID(t *testing.T) {
	setup()
	user := entity.User{
		Login: "testuser",
		Pass:  "password",
		Name:  "Test User",
	}
	createdUser := CreateUser(user)

	retrievedUser, err := GetUserByID(createdUser.Uid)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.Uid, retrievedUser.Uid)
	assert.Equal(t, createdUser.Login, retrievedUser.Login)
	assert.Equal(t, createdUser.Pass, retrievedUser.Pass)
	assert.Equal(t, createdUser.Name, retrievedUser.Name)
}

func TestGetUserByIDNotFound(t *testing.T) {
	setup()
	_, err := GetUserByID(1)
	assert.Error(t, err)
	assert.Equal(t, "User не найден", err.Error())
}
