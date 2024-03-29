package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// host=11.0.1.152,port=6380,db=3
func main() {
	startTimestamp := time.Now().Unix()
	defer func() {
		endTimestamp := time.Now().Unix()
		fmt.Printf("exec end, timestamp:%s, cost:%+v\n", time.Now(), endTimestamp-startTimestamp)
	}()
	fmt.Printf("exec start,timestamp:%s\r\n", time.Now())
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "11.0.1.152:6380", // 192.168.1.222:6379
		Password: "",
		DB:       3,
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
		fmt.Printf("key:%+v\r\n", key)
	}
}
