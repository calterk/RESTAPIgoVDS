package main

import (
	"app/entity"
	"net/http"
	"testing"
)

func TestWebServ(t *testing.T) {
	// Запуск веб-сервера в горутине
	go WebServ()

	// Попытка подключиться к веб-серверу
	resp, err := http.Get("http://localhost:8090")
	if err != nil {
		t.Errorf("Ошибка при подключении к серверу: %v", err)
	}
	defer resp.Body.Close()

	// Проверка, что сервер возвращает статус OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус OK, получен: %v", resp.Status)
	}
}

func TestInitDB(t *testing.T) {
	// Вызываем функцию initDB
	db, err := initDB()
	if err != nil {
		t.Errorf("Ошибка при инициализации базы данных: %v", err)
	}

	// Проверка, что соединение с базой данных открыто
	if db == nil {
		t.Error("Соединение с базой данных не открыто")
	}

	// Проверка, что миграции выполнены успешно
	if !db.Migrator().HasTable(&entity.User{}) {
		t.Error("Таблица User не создана")
	}
	if !db.Migrator().HasTable(&entity.Vds{}) {
		t.Error("Таблица Vds не создана")
	}
	if !db.Migrator().HasTable(&entity.Token{}) {
		t.Error("Таблица Token не создана")
	}

	// Закрытие соединения с базой данных
	sqlDB, err := db.DB()
	if err != nil {
		t.Errorf("Ошибка соединения с базой данных: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Errorf("Ошибка соединения с базой данных: %v", err)
	}
}
