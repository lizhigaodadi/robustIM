package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
	"time"
)

var rdb *redis.Client

func Init() {
	/*TODO: Init Redis Client*/
	//endpoints := config.GetRedisEndpoint()
	//if len(endpoints) == 0 {
	//	panic("Init Redis Client Error")
	//}
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.5.38:6379",
		Password: "",
		DB:       0,
	})

	/*Test Client Connect To Redis Server*/
	pong, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatalf("Connnect To Redis Server Failed: %v\n", err)
	}
	log.Printf("Success Connect To Redis: %s", pong)

	initLuaScript(context.TODO())
}

func SetString(ctx context.Context, key string, val interface{}, expireTime time.Duration) error {
	cmd := rdb.Set(ctx, key, val, expireTime)
	if cmd == nil {
		return errors.New("Set Sds To Redis Failed\n")
	}
	return cmd.Err()
}

func SetBytes(ctx context.Context, key string, val []byte, expireTime time.Duration) error {
	cmd := rdb.Set(ctx, key, val, expireTime)
	if cmd == nil {
		return errors.New("Set Bytes To Redis Failed\n")
	}
	return cmd.Err()
}

func GetBytes(ctx context.Context, key string) ([]byte, error) {
	cmd := rdb.Conn().Get(ctx, key)
	if cmd == nil {
		return nil, errors.New("Set Sds To Redis Failed\n")
	}

	return cmd.Bytes()
}

func SAdd(ctx context.Context, key string, val interface{}) error {
	cmd := rdb.SAdd(ctx, key, val)
	if cmd == nil {
		return errors.New("Redis SAdd Failed\n")
	}
	return cmd.Err()
}

func SRem(ctx context.Context, key string, members ...interface{}) error {
	cmd := rdb.Conn().SRem(ctx, key, members)

	if cmd == nil {
		return errors.New("Redis SRem Failed\n")
	}
	return cmd.Err()
}

func GetString(ctx context.Context, key string) (string, error) {
	cmd := rdb.Conn().Get(ctx, key)
	if cmd == nil {
		return "", errors.New("Redis SRem Cmd is Nil\n")
	}
	return cmd.String(), nil
}
func SMemberStringSlice(ctx context.Context, key string) ([]string, error) {
	members := rdb.Conn().SMembers(ctx, key)
	if members == nil {
		return nil, errors.New("SMember Get is Nil")
	}

	return members.Result()
}

func ExecuteLuaScript(ctx context.Context, luaScriptName string, keys []string, args ...interface{}) (int, error) {
	script, ok := luaScriptTable[luaScriptName]
	if !ok {
		return -1, errors.New(fmt.Sprintf("Script Name: %sNot Exists", luaScriptName))
	}
	cmd := rdb.EvalSha(ctx, script.Sha, keys, args)
	if cmd == nil {
		return -1, errors.New("Execute Lua Failed\n")
	}

	return cmd.Int()
}
