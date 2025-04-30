package service

import (
	"app/engine"
	"app/entity"
	"app/storage"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// InitUserStorage инициализирует userMap
func InitUserStorage() {
	st.userMap = make(map[uint64]entity.User)
}

func (e *Service) CreateUser(ctx *engine.Context) error { //POST
	decoder := json.NewDecoder(ctx.Request.Body)
	var user entity.User
	err := decoder.Decode(&user)
	if err != nil {
		return err
	}
	user = storage.CreateUser(user)
	coder, _ := json.Marshal(user)
	ctx.Response.Write(coder)
	return nil

}

func (e *Service) UpdateUser(ctx *engine.Context) error { //PUT
	decoder := json.NewDecoder(ctx.Request.Body)
	var user entity.User
	err := decoder.Decode(&user)
	if err != nil {
		return err
	}

	if err := storage.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (e *Service) DeleteUser(ctx *engine.Context) error { //DELETE
	path := strings.Split(ctx.Request.URL.Path, "/")

	uid, err := strconv.ParseUint(path[len(path)-1], 10, 64)
	if err != nil {
		return err
	}

	err = storage.DeleteUser(uid)
	if err != nil {
		return err
	}

	ctx.Response.WriteHeader(http.StatusOK)
	return nil
}

func (e *Service) GetUser(ctx *engine.Context) error { //GET
	userList, err := storage.GetAllUsers()
	if err != nil {
		return err
	}

	coder, _ := json.Marshal(userList)
	ctx.Response.Write(coder)

	return nil

}

func (e *Service) GetByUser(ctx *engine.Context) error { //GET
	uidStr := ctx.Request.URL.Query().Get("uid")
	fmt.Println(uidStr)
	if uidStr == "" {
		return errors.New("Параметр 'uid' не указан")
	}

	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		return err
	}
	fmt.Println(uid)
	user, err := storage.GetUserByID(uid)
	if err != nil {
		return err
	}

	coder, err := json.Marshal(user)
	if err != nil {
		return err
	}

	ctx.Response.Write(coder)
	return nil
}
