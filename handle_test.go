package main

import (
	"app/entity"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTestDB() (*gorm.DB, error) {
	// Настройка подключения к тестовой базе данных PostgreSQL
	dsn := "user=postgres password=postgre dbname=testdb"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Создаем таблицы
	err = db.AutoMigrate(&entity.User{}, &entity.Token{}, &entity.Vds{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestAuthHandler(t *testing.T) {
	// Настройка тестовой базы данных
	var err error
	db, err = setupTestDB()
	assert.NoError(t, err)

	// Очистка базы данных после тестов
	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM users")
	}()

	// Создаем mock данные пользователя
	mockUser := entity.User{
		Uid:    1,
		Login:  "testuser",
		Pass:   "password",
		Name:   "Test User",
		Access: true,
	}

	// Добавляем mock данные пользователя в базу данных
	err = db.Create(&mockUser).Error
	assert.NoError(t, err)

	// Создаем mock запрос
	authData := SAuth{
		Login: "testuser",
		Pass:  "password",
	}
	jsonData, err := json.Marshal(authData)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	AuthHandler(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockUser.Uid, uint64(response["uid"].(float64)))
	assert.Equal(t, mockUser.Name, response["name"])
	assert.Equal(t, mockUser.Access, response["access"])
	assert.NotEmpty(t, response["token"])
}

func TestGetUsersTableHandler(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Очистка базы данных после тестов
	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM users")
	}()

	// Создаем mock данные пользователей
	mockUsers := []entity.User{
		{Uid: 1, Login: "user1", Pass: "password", Name: "User One", Access: true},
		{Uid: 2, Login: "user2", Pass: "password", Name: "User Two", Access: true},
	}

	for _, user := range mockUsers {
		err = db.Create(&user).Error
		assert.NoError(t, err)
	}

	// Создаем mock запрос
	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	assert.NoError(t, err)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	getUsersTableHandler(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response []entity.User
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, mockUsers[0].Login, response[0].Login)
	assert.Equal(t, mockUsers[1].Login, response[1].Login)
}

func TestGetUsersTableHandlerFalse(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Очистка базы данных после тестов
	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM users")
	}()

	// Создаем mock данные пользователей
	mockUsers := []entity.User{
		{Uid: 1, Login: "user1", Pass: "password", Name: "User One", Access: true},
		{Uid: 2, Login: "user2", Pass: "password", Name: "User Two", Access: false},
	}

	for _, user := range mockUsers {
		err = db.Create(&user).Error
		assert.NoError(t, err)
	}

	// Создаем mock запрос
	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	assert.NoError(t, err)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	getUsersTableHandlerfalse(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Выводим тело ответа для отладки
	t.Logf("Response Body: %s", rr.Body.String())

	assert.Equal(t, "f", response)
}

func TestGetVdsTableHandler(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Очистка базы данных после тестов
	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM vds")
	}()

	// Создаем mock данные пользователей
	mockVds := []entity.Vds{
		{Vid: 1, Uid: 1, Login: "vds1", Pass: "password", Name: "One", Status: true},
		{Vid: 2, Uid: 2, Login: "vds2", Pass: "password", Name: "Two", Status: false},
	}

	for _, vds := range mockVds {
		err = db.Create(&vds).Error
		assert.NoError(t, err)
	}

	// Создаем mock запрос
	req, err := http.NewRequest(http.MethodGet, "/vds", nil)
	assert.NoError(t, err)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	getVdsTableHandler(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response []entity.Vds
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, mockVds[0].Login, response[0].Login)
	assert.Equal(t, mockVds[1].Login, response[1].Login)
}

func TestGetUserVdsTableHandler(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Очистка базы данных после тестов
	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM vds")
		_ = db.Exec("DELETE FROM users")
	}()

	// Создаем mock данные пользователей и VDS
	mockUser := entity.User{
		Uid:    1,
		Login:  "testuser",
		Pass:   "password",
		Name:   "Test User",
		Access: true,
	}

	err = db.Create(&mockUser).Error
	assert.NoError(t, err)

	mockVds := []entity.Vds{
		{Vid: 1, Uid: 1, Login: "vds1", Pass: "password", Name: "One", Status: true},
		{Vid: 2, Uid: 1, Login: "vds2", Pass: "password", Name: "Two", Status: false},
		{Vid: 3, Uid: 2, Login: "vds3", Pass: "password", Name: "Three", Status: true}, // VDS для другого пользователя
	}

	for _, vds := range mockVds {
		err = db.Create(&vds).Error
		assert.NoError(t, err)
	}

	// Создаем mock токен для пользователя
	mockToken := generateToken(mockUser.Uid)

	// Создаем mock запрос с токеном в заголовке
	req, err := http.NewRequest(http.MethodGet, "/getUserVdsTable", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", mockToken)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	getUserVdsTableHandler(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response []entity.Vds
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, mockVds[0].Login, response[0].Login)
	assert.Equal(t, mockVds[1].Login, response[1].Login)
}

func TestHandleVdsButtonClick(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Очистка базы данных после тестов
	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM vds")
		_ = db.Exec("DELETE FROM users")
	}()

	// Создаем mock данные пользователя и VDS
	mockUser := entity.User{
		Uid:    1,
		Login:  "testuser",
		Pass:   "password",
		Name:   "Test User",
		Access: true,
	}

	err = db.Create(&mockUser).Error
	assert.NoError(t, err)

	mockVds := entity.Vds{
		Vid: 1, Uid: 1, Login: "vds1", Pass: "password", Name: "One", Status: true,
	}

	err = db.Create(&mockVds).Error
	assert.NoError(t, err)

	// Создаем mock токен для пользователя
	mockToken := generateToken(mockUser.Uid)

	// Создаем mock запрос с токеном в заголовке и VDS ID в теле
	vdsId := mockVds.Vid
	reqBody, err := json.Marshal(map[string]interface{}{
		"vdsId": vdsId,
	})
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/handleVdsButtonClick", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Authorization", mockToken)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	handleVdsButtonClick(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
}

func TestDeleteUserHandler(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Очистка базы данных после тестов
	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM users")
	}()

	// Создаем mock данные пользователя
	mockUser := entity.User{
		Uid:    1,
		Login:  "testuser",
		Pass:   "password",
		Name:   "Test User",
		Access: true,
	}

	err = db.Create(&mockUser).Error
	assert.NoError(t, err)

	// Создаем mock запрос с ID пользователя в теле
	reqBody, err := json.Marshal(map[string]interface{}{
		"uid": mockUser.Uid,
	})
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/deleteUser", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	deleteUserHandler(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User deleted successfully", response["message"])

	// Проверяем, что пользователь был удален из базы данных
	var count int64
	db.Model(&entity.User{}).Where("uid = ?", mockUser.Uid).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteVdsHandlerAdmin(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Очистка базы данных после тестов
	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM vds")
		_ = db.Exec("DELETE FROM users")
	}()

	// Создаем mock данные пользователя-администратора
	mockAdmin := entity.User{
		Uid:    1,
		Login:  "admin",
		Pass:   "password",
		Name:   "Admin User",
		Access: true,
	}

	err = db.Create(&mockAdmin).Error
	assert.NoError(t, err)

	// Создаем mock данные VDS
	mockVds := entity.Vds{
		Vid: 1, Uid: 2, Login: "vds1", Pass: "password", Name: "One", Status: true,
	}

	err = db.Create(&mockVds).Error
	assert.NoError(t, err)

	// Создаем mock запрос с ID VDS в теле
	reqBody, err := json.Marshal(map[string]interface{}{
		"vid": mockVds.Vid,
	})
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/deleteVds", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	deleteVdsHandlerAdmin(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Vds deleted successfully", response["message"])

	// Проверяем, что VDS была удалена из базы данных
	var count int64
	db.Model(&entity.Vds{}).Where("vid = ?", mockVds.Vid).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestTurnVdsHandler(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM vds")
		_ = db.Exec("DELETE FROM users")
	}()

	// Создаем mock данные пользователя и VDS
	mockUser := entity.User{
		Uid:    1,
		Login:  "testuser",
		Pass:   "password",
		Name:   "Test User",
		Access: true,
	}

	err = db.Create(&mockUser).Error
	assert.NoError(t, err)

	mockVds := entity.Vds{
		Vid: 1, Uid: 1, Login: "vds1", Pass: "password", Name: "One", Status: true,
	}

	err = db.Create(&mockVds).Error
	assert.NoError(t, err)

	// Создаем mock токен для пользователя
	mockToken := generateToken(mockUser.Uid)

	// Создаем mock запрос с токеном в заголовке и VDS ID в теле
	vdsId := mockVds.Vid
	reqBody, err := json.Marshal(map[string]interface{}{
		"vid": vdsId,
	})
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/turnVds", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Authorization", mockToken)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	turnVdsHandler(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "VDS status updated successfully", response["message"])

	// Проверяем, что статус VDS был обновлен
	var updatedVds entity.Vds
	db.First(&updatedVds, mockVds.Vid)
	assert.Equal(t, !mockVds.Status, updatedVds.Status)
}

func TestDeleteVdsHandler(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM vds")
		_ = db.Exec("DELETE FROM users")
	}()

	// Создаем mock данные пользователя и VDS
	mockUser := entity.User{
		Uid:    1,
		Login:  "testuser",
		Pass:   "password",
		Name:   "Test User",
		Access: true,
	}

	err = db.Create(&mockUser).Error
	assert.NoError(t, err)

	mockVds := entity.Vds{
		Vid: 1, Uid: 1, Login: "vds1", Pass: "password", Name: "One", Status: true,
	}

	err = db.Create(&mockVds).Error
	assert.NoError(t, err)

	// Создаем mock токен для пользователя
	mockToken := generateToken(mockUser.Uid)

	// Создаем mock запрос с токеном в заголовке и VDS ID в теле
	vdsId := mockVds.Vid
	reqBody, err := json.Marshal(map[string]interface{}{
		"vid": vdsId,
	})
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/deleteVds", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Authorization", mockToken)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	deleteVdsHandler(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "VDS deleted successfully", response["message"])

	// Проверяем, что VDS была удалена из базы данных
	var count int64
	db.Model(&entity.Vds{}).Where("vid = ?", mockVds.Vid).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestAddUserHandler(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM users")
	}()

	// Создаем mock данные пользователя
	mockUser := entity.User{
		Uid:    1,
		Login:  "testuser",
		Pass:   "password",
		Name:   "Test User",
		Access: true,
	}

	// Создаем mock запрос с данными пользователя в теле
	reqBody, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/addUser", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	addUserHandler(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User added successfully", response["message"])

	// Проверяем, что пользователь был добавлен в базу данных
	var count int64
	db.Model(&entity.User{}).Where("uid = ?", mockUser.Uid).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestGenerateToken(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	defer func() {
		_ = db.Exec("DELETE FROM tokens")
		_ = db.Exec("DELETE FROM users")
	}()

	// Создаем mock данные пользователя
	mockUser := entity.User{
		Uid:    1,
		Login:  "testuser",
		Pass:   "password",
		Name:   "Test User",
		Access: true,
	}

	err = db.Create(&mockUser).Error
	assert.NoError(t, err)

	// Генерируем токен для пользователя
	tokenString := generateToken(mockUser.Uid)

	// Проверяем, что токен был добавлен в базу данных
	var token entity.Token
	result := db.Where("uuid = ?", tokenString).First(&token)
	assert.NoError(t, result.Error)
	assert.Equal(t, mockUser.Uid, token.Uid)

	// Проверяем, что срок действия токена установлен корректно
	expirationTime := time.Now().Add(24 * time.Hour)
	assert.True(t, token.Exp.After(time.Now()))
	assert.True(t, expirationTime.After(token.Exp))
}

func TestHandle(t *testing.T) {
	// Создаем mock запрос
	reqBody := []byte(`{"key": "value"}`)
	req, err := http.NewRequest(http.MethodPost, "/handle", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	// Создаем mock ответ
	rr := httptest.NewRecorder()

	// Вызываем тестируемую функцию
	handle(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusOK, rr.Code)
}
