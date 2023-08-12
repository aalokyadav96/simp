package main
import (
	"context"
	"fmt"
	"log"
	"os"

	gredis "github.com/redis/go-redis/v9"
)

var conn *gredis.Client

func init() {
	conn = gredis.NewClient(&gredis.Options{
		Addr:         ":6379",
	})
}
var opt, _ = gredis.ParseURL(os.Getenv("rdx"))
var conn = gredis.NewClient(opt)


func init() {
	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
}

func rdxSet(key, value string) error {

	ctx := context.Background()

	_, err := conn.Set(ctx, key, value, 0).Result()
	if err != nil {
		return fmt.Errorf("error while doing SET command in gredis : %v", err)
	}

	return err

}

func rdxGet(key string) (string, error) {

	ctx := context.Background()

	value, err := conn.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("error while doing GET command in gredis : %v", err)
	}

	return value, err
}


func rdxDel(key string) (string, error) {

	ctx := context.Background()

	value, err := conn.Del(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("error while doing DEL command in gredis : %v", err)
	}

	return ""+string(value), err
}

func rdxHset(hash, key, value string) error {

	ctx := context.Background()

	_, err := conn.HSet(ctx, hash, key, value).Result()
	if err != nil {
		return fmt.Errorf("error while doing HSET command in gredis : %v", err)
	}

	return err
}

func rdxHget(hash, key string) (string, error) {

	ctx := context.Background()

	value, err := conn.HGet(ctx, hash, key).Result()
	if err != nil {
		return "error : ", err
	}

	return value, err

}

func rdxHdel(hash, key string) (string, error) {

	ctx := context.Background()

	value, err := conn.HDel(ctx, hash, key).Result()
	if err != nil {
		return string(value), fmt.Errorf("error while doing HGET command in gredis : %v", err)
	}

	return string(value), err

}

func rdxHgetall(hash string) map[string]string {

	ctx := context.Background()
	value, _ := conn.HGetAll(ctx, hash).Result()

	return value

}
