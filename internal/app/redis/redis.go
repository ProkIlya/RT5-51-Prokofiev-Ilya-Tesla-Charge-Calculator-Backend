package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

func New(host string, port int) (*Client, error) {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: host + ":" + strconv.Itoa(port),
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) SetJWTBlacklist(ctx context.Context, token string, expiresIn time.Duration) error {
	return c.client.Set(ctx, "jwt_blacklist:"+token, "1", expiresIn).Err()
}

func (c *Client) IsJWTBlacklisted(ctx context.Context, token string) (bool, error) {
	result, err := c.client.Get(ctx, "jwt_blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result == "1", nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
