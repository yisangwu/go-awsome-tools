package bak

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// host=11.0.1.152,port=6380,db=3
func main() {

	// addr :="192.168.1.222:6379"
	// db := 1

	addr :="127.0.0.1:6380"
	db := 3

	startTimestamp := time.Now().Unix()
	defer func() {
		endTimestamp := time.Now().Unix()
		fmt.Printf("exec end, timestamp:%s, cost:%+v\n", time.Now(), endTimestamp-startTimestamp)
	}()
	fmt.Printf("exec start,timestamp:%s\r\n", time.Now())
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})
	prefix := "RRRR_"
	iter := client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	if iter == nil {
		fmt.Printf("client scan not found\r\n")
		return
	}
	for iter.Next(ctx) {
		key := iter.Val()
		if key == "" {
			continue
		}
		key_mem := client.MemoryUsage(ctx, key)
		key_ttl := client.TTL(ctx, key)
		fmt.Printf("%+v, %+v, %+v\r\n", key, key_mem, key_ttl)
	}
}
