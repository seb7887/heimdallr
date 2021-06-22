package storage

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

type Client struct {
	clientId  string
	publicKey string
}

type Repository interface {
	Health(ctx context.Context) error
	CreateClient(ctx context.Context, client *Client) error
}

type redisRepository struct {
	pool *redis.Pool
}

func NewRepository(pool *redis.Pool) Repository {
	return redisRepository{
		pool: pool,
	}
}

func (r redisRepository) Health(ctx context.Context) error {
	// get a connection
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}

	// test the connection
	_, err = conn.Do("PING")
	if err != nil {
		return err
	}

	return nil
}

func (r redisRepository) CreateClient(ctx context.Context, client *Client) error {
	// TODO
	return nil
}
