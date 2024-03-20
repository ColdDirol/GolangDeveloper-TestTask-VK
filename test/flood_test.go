package test

import (
	"GolangDeveloper-TestTask-VK/internal/config"
	"GolangDeveloper-TestTask-VK/internal/flood_control"
	"GolangDeveloper-TestTask-VK/internal/storage"
	"context"
	"testing"
	"time"
)

func TestFloodControl_Check(t *testing.T) {
	// Создаем экземпляр FloodControl с конфигурацией
	config := config.Config{
		IntervalSeconds: 10,
		MaxRequests:     3,
		Redis: config.Redis{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
	}

	fc := flood_control.NewFloodControl(config, storage.InitRedis(config.Redis))

	// userID для тестов
	userID := int64(123)

	// Случай 1: Проверяем, что при вызове метода Check меньше K раз за последние N секунд, возвращается true
	result, err := fc.Check(context.Background(), userID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != true {
		t.Error("Expected result to be true")
	}

	// Случай 2: Проверяем, что при вызове метода Check K или более раз за последние N секунд, возвращается false
	for i := 0; i < config.MaxRequests; i++ {
		_, err = fc.Check(context.Background(), userID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
	result, err = fc.Check(context.Background(), userID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != false {
		t.Error("Expected result to be false")
	}

	// Случай 3: Проверяем, что при вызове метода Check в разных интервалах времени количество вызовов учитывается корректно
	time.Sleep(time.Second * time.Duration(config.IntervalSeconds))
	result, err = fc.Check(context.Background(), userID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != true {
		t.Error("Expected result to be true")
	}
}
