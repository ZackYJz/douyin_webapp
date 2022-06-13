package redis

import (
	"fmt"
	"go_webapp/pkg/settings"

	"github.com/go-redis/redis"
)

var (
	rdb *redis.Client
	Nil = redis.Nil
)

//type SliceCmd = redis.SliceCmd
//type StringStringMapCmd = redis.StringStringMapCmd

// Init 初始化连接
func Init(cfg *settings.RedisSettingS) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.Db,       // use default DB
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = rdb.Close()
}
