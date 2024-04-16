package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	h       bool
	host    string
	port    int
	db      int
	key_pre string
)

func init() {
	// 定义命令行参数
	flag.BoolVar(&h, "h", false, "show help")
	flag.StringVar(&host, "host", "", "redis host")
	flag.IntVar(&port, "port", 3306, "redis port")
	flag.IntVar(&db, "db", 0, "redis db,default 0")
	flag.StringVar(&key_pre, "key_pre", "", "key prefix")

	// 自定义帮助信息
	flag.Usage = func() {

		fmt.Fprintf(os.Stderr, "redis key export to file\r\n")
		fmt.Fprintf(os.Stderr, "Usage: \r\n")
		fmt.Fprintf(os.Stderr, `  RedisKeyExport -host="localhost" -port=3306 -db=0 -key_pre="AAA_" >redis_key.txt
`)
		fmt.Fprintf(os.Stderr, "Options:\r\n")
		flag.PrintDefaults()
	}
	// 解析命令行参数
	flag.Parse()

}

func main() {

	// 如果没有提供任何命令行参数，则打印帮助信息
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	if host == "" {
		fmt.Println("Error! All args must have values")
		flag.Usage()
		return
	}

	startTimestamp := time.Now().UnixMilli()
	defer func() {
		endTimestamp := time.Now().UnixMilli()
		fmt.Printf("exec end, timestamp:%s, cost:%+vms\n", time.Now(), endTimestamp-startTimestamp)
	}()
	fmt.Printf("exec start,timestamp:%s\r\n", time.Now())
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})

	iter := client.Scan(ctx, 0, key_pre+"*", 0).Iterator()
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
