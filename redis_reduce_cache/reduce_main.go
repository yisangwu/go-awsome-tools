package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// initViper
func initViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName("conf")
	v.SetConfigType("yaml")
	v.AddConfigPath("./")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return v
}

var v = initViper()

type RedisConfStruct struct {
	Host                   string `json:"host"`
	Port                   int    `json:"port"`
	Db                     int    `json:"db"`
	KeyPrefLabelTxt        string `json:"key_label_txt_pref"`
	KeyPrefLabelShortVideo string `json:"key_label_short_pref"`
}

// GetRedisConf 获取Redis配置
func GetRedisConf() RedisConfStruct {
	return RedisConfStruct{
		Host:                   strings.TrimSpace(v.GetString("redis.host")),
		Port:                   v.Get("redis.port").(int),
		Db:                     v.Get("redis.db").(int),
		KeyPrefLabelTxt:        strings.TrimSpace(v.GetString("key_pref.key_pref_label_txt")),
		KeyPrefLabelShortVideo: strings.TrimSpace(v.GetString("key_pref.key_pref_label_short_video")),
	}
}

var MathUserId = v.Get("math_userid").(int)

func main() {
	startTimestamp := time.Now().Unix()
	defer func() {
		endTimestamp := time.Now().Unix()
		fmt.Printf("exec end, timestamp:%s, cost:%+v\n", time.Now(), endTimestamp-startTimestamp)
	}()
	fmt.Printf("exec start,timestamp:%s\r\n", time.Now())

	// redis
	redisConf := GetRedisConf()
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password: "",
		DB:       redisConf.Db,
	})
	prefix := "*WEIGHT*"
	iter := client.Scan(ctx, 0, prefix, 0).Iterator()
	if iter == nil {
		fmt.Printf("client scan not found\r\n")
		return
	}
	for iter.Next(ctx) {
		key := iter.Val()
		if key == "" {
			continue
		}
		if MathUserId != 0 && !strings.Contains(key, strconv.Itoa(MathUserId)) {
			continue
		}
		find_txt := strings.Contains(key, redisConf.KeyPrefLabelTxt)
		find_video := strings.Contains(key, redisConf.KeyPrefLabelShortVideo)
		if !find_txt && !find_video {
			continue
		}
		// zcard
		count_zcard, err_zcard := client.ZCard(ctx, key).Result()
		if err_zcard != nil && err_zcard != redis.Nil {
			fmt.Printf("111111, key:%+v, err:%+v, err_t:%T\r\n",key, err_zcard, err_zcard)
			continue
		}
		if count_zcard <= 100 {
			continue
		}
		result, errRem := client.ZRemRangeByRank(ctx, key, 0, -101).Result()
		if errRem != nil && errRem != redis.Nil {
			fmt.Printf("ZRemRangeByRank failed, key:%s, err:%+v\r\n", key, errRem)
			continue
		}
		fmt.Printf("ZRemRangeByRank key succ, key:%s, ret:%+v\r\n", key, result)
	}
}
