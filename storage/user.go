package storage

import (
	"app/entity"
	"errors"
	"sync"
)

type UserStorage struct {
	userMap map[uint64]entity.User
	next    uint64
	mutx    sync.RWMutex
}

var st *UserStorage

func init() {
	st = new(UserStorage)
}
func CreateUser(user entity.User) entity.User {
	st.mutx.Lock()
	defer st.mutx.Unlock()
	user.Uid = st.next
	st.next++
	st.userMap[user.Uid] = user
	return user
}

func UpdateUser(user entity.User) error {
	st.mutx.Lock()
	defer st.mutx.Unlock()
	if _, exists := st.userMap[user.Uid]; !exists {
		return errors.New("User не найден")
	}
	st.userMap[user.Uid] = user
	return nil
}

func DeleteUser(uid uint64) error {
	st.mutx.Lock()
	defer st.mutx.Unlock()
	if _, exists := st.userMap[uid]; !exists {
		return errors.New("User не найден")
	}
	delete(st.userMap, uid)
	return nil
}

func GetAllUsers() ([]entity.User, error) {
	st.mutx.RLock()
	defer st.mutx.RUnlock()
	userList := make([]entity.User, len(st.userMap))
	var inc int
	for _, user := range st.userMap {
		userList[inc] = user
		inc++
	}
	return userList, nil
}

func GetUserByID(uid uint64) (*entity.User, error) {
	st.mutx.RLock()
	defer st.mutx.RUnlock()
	user, exists := st.userMap[uid]
	if !exists {
		return nil, errors.New("User не найден")
	}
	return &user, nil
}
