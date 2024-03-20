package main

import (
	"GolangDeveloper-TestTask-VK/internal/config"
	"GolangDeveloper-TestTask-VK/internal/flood_control"
	"GolangDeveloper-TestTask-VK/internal/logger"
	"GolangDeveloper-TestTask-VK/internal/storage"
	"context"
	"log/slog"
	"math"
	"os"
	"strconv"
	"time"
)

func main() {
	config := config.InitConfig("config.json")
	log := logger.SetupLogger(config.ENV)

	// Инициализация Redis клиента
	client := storage.InitRedis(config.Redis)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		//fmt.Println("Error connecting to Redis", err)
		log.Error("Error connecting to Redis", logger.Err(err))
		os.Exit(1)
	}

	ctx := context.Background()
	fc := flood_control.NewFloodControl(*config, client)

	var count int64 = 1

	for {
		userID := count
		for j := 0; j < 10; j++ {
			ok, err := fc.Check(ctx, userID)
			if err != nil {
				//fmt.Printf("Error checking flood control: %v\n", err)
				log.Error("Error checking flood control", logger.Err(err))
				os.Exit(1)
			}
			if ok {
				//fmt.Printf("Request allowed for user with %d ID\n", count)
				log.Info("Request allowed", slog.String("userID", strconv.FormatInt(count, 10)))
			} else {
				//fmt.Printf("Request denied for user with %d ID\n", count)
				log.Info("Request denied", slog.String("userID", strconv.FormatInt(count, 10)))
			}
			time.Sleep(1 * time.Second)
		}

		if count == math.MaxInt64 {
			break
		}

		count++
		time.Sleep(1 * time.Second)
	}
}
