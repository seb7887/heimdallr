package storage

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/seb7887/heimdallr/config"
	log "github.com/sirupsen/logrus"
)

func NewConn(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				log.Fatalf("%s", err.Error())
			}
			return c, nil
		},
	}
}

func InitializeRepository() Repository {
	addr := config.GetConfig().RedisHost
	pool := NewConn(addr)

	return NewRepository(pool)
}
