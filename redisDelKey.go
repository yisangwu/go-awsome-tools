package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

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

	fileName := "more_three_days_key.txt"
	f_handle, err := os.Open(fileName)
	if err != nil {
		fmt.Println("打开错误", err)
		return
	}
	defer f_handle.Close()

	fileScanner := bufio.NewScanner(f_handle)

	for fileScanner.Scan(){
		line :=fileScanner.Text()
		if line == "" {
			break
		}
		str_arr := strings.Split(string(line), `,`)
		if len(str_arr)==0 {
			fmt.Printf("strings.Split empry, :%s\r\n", time.Now())
			break
		}
		key := strings.TrimSpace(string(str_arr[0]))
		client.Del(ctx, key)
		fmt.Printf("redis delete key::%s\r\n", key)
	}
	f_handle.Close()
}
