package queue

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Successfully connected to Redis")
	return &RedisClient{Client: client}
}

func (rc *RedisClient) Close() error {
	return rc.Client.Close()
}

func (rc *RedisClient) PushJob(queueName, jobData string) error {
	ctx := context.Background()
	return rc.Client.LPush(ctx, queueName, jobData).Err()
}

func (rc *RedisClient) Subscribe(channel string) *redis.PubSub {
	ctx := context.Background()
	return rc.Client.Subscribe(ctx, channel)
}
