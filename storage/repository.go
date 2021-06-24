package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
)

const (
	blacklist string = "blacklist"
)

type Client struct {
	ClientId  string
	PublicKey string
}

type Repository interface {
	Health(ctx context.Context) error
	CreateClient(ctx context.Context, client *Client) error
	GetClientKey(ctx context.Context, clientId string) ([]byte, error)
	DeleteClient(ctx context.Context, clientId string) error
	UpsertBlacklist(ctx context.Context, newBlacklist []string) error
	GetBlacklist(ctx context.Context) ([]string, error)
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
	// get a connection
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}

	// set a new key/value pair
	_, err = conn.Do("SET", client.ClientId, client.PublicKey)
	return err
}

func (r redisRepository) GetClientKey(ctx context.Context, clientId string) ([]byte, error) {
	// get a connection
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	data, err := redis.Bytes(conn.Do("GET", clientId))
	if err != nil {
		if err == redis.ErrNil {
			return nil, fmt.Errorf("Client does not exist")
		}
		return nil, err
	}

	return data, nil
}

func (r redisRepository) DeleteClient(ctx context.Context, clientId string) error {
	// get a connection
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}

	// delete entry
	_, err = conn.Do("DEL", clientId)
	if err == redis.ErrNil {
		return fmt.Errorf("Client does not exist")
	}
	return err
}

func (r redisRepository) UpsertBlacklist(ctx context.Context, newBlacklist []string) error {
	// get a connection
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return err
	}

	data := strings.Join(newBlacklist[:], ",")
	_, err = conn.Do("SET", blacklist, data)
	return err
}

func (r redisRepository) GetBlacklist(ctx context.Context) ([]string, error) {
	// get a connection
	conn, err := r.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}

	data, err := redis.Bytes(conn.Do("GET", blacklist))
	if err != nil {
		if err == redis.ErrNil {
			return []string{}, nil
		}
		return nil, err
	}
	blacklistStr := string(data)

	return strings.Split(blacklistStr, ","), nil
}
