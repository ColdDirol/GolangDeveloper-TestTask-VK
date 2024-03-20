package flood_control

import (
	"GolangDeveloper-TestTask-VK/internal/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"time"
)

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

type FloodControlImpl struct {
	config  config.Config
	client  *redis.Client
	storage map[int64][]time.Time
	mutex   sync.Mutex
}

func NewFloodControl(config config.Config, client *redis.Client) FloodControl {
	return &FloodControlImpl{
		config: config,
		client: client,
	}
}

func (fc *FloodControlImpl) Check(ctx context.Context, userID int64) (bool, error) {
	currentTime := time.Now()
	intervalStart := currentTime.Add(-time.Duration(fc.config.IntervalSeconds) * time.Second)

	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	// Проверяем количество вызовов за последние N секунд из Redis
	// .ZCount(ctx, key, min, max)
	count, err := fc.client.ZCount(ctx, "requests:"+strconv.FormatInt(userID, 10), strconv.FormatInt(intervalStart.Unix(), 10), strconv.FormatInt(currentTime.Unix(), 10)).Result()
	fmt.Println(count)
	if err != nil {
		return false, err
	}

	if count >= int64(fc.config.MaxRequests) {
		return false, nil
	}

	// Добавляем текущий вызов в Redis
	_, err = fc.client.ZAdd(
		ctx, "requests:"+strconv.FormatInt(userID, 10), redis.Z{
			Score:  float64(currentTime.Unix()),
			Member: currentTime.Unix(),
		}).Result()
	if err != nil {
		return false, err
	}

	return true, nil
}
