package test

import (
	"context"
	"fmt"
	"github.com/OmnitechX/snowflake"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestRedisIdGenerator(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "111.13.233.34:7001",
		Password: "Bfsw@123",
		DB:       0,
	})
	ctx := context.Background()
	g := snowflake.NewRedisGenerator(ctx, client, "test")
	for i := 0; i < 10; i++ {
		fmt.Println(g.Next())
	}

	time.Sleep(2 * time.Minute)
}
