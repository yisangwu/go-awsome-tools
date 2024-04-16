package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type KeyMemStruct struct{
	key string
	mem int
}

func main(){

	var keyMemSlice []KeyMemStruct
	total_bytes := 0
	redis_db_file := "redis_6380_3_keys.txt"


	f_handle, err := os.Open(redis_db_file)
	if err!= nil{
		fmt.Printf("parse redis key file failed, file:%s, err:%+v", redis_db_file, err)
		return
	}
	defer f_handle.Close()

	fileScanner := bufio.NewScanner(f_handle)
	for fileScanner.Scan(){
		line := fileScanner.Text()
		if line == "" {
			continue
		}
		// RRRR_REC_S_U_PD_AMOUNT_123, memory usage RRRR_REC_S_U_PD_AMOUNT_123: 72, ttl RRRR_REC_S_U_PD_AMOUNT_123: 44h33m23s
		str_arr := strings.Split(string(line), ",")
		if len(str_arr) == 0 {
			continue
		}
		key := strings.TrimSpace(str_arr[0])
		// RRRR_REC_S_U_PD_AMOUNT_123: 72
		mem_arr := strings.Split(strings.TrimSpace(string(str_arr[1])), ":")
		if len(mem_arr) == 0 {
			continue
		}
		key_mem, _ := strconv.Atoi(strings.TrimSpace(mem_arr[1]))
		// slice
		keyMemSlice = append(keyMemSlice, KeyMemStruct{key:key, mem: key_mem})

		// total
		total_bytes += key_mem
	}

	// 输出总大小
	fmt.Printf("total key: %d, total memory usage:%dMB\r\n", len(keyMemSlice) ,total_bytes/1024/1024)
	// 使用sort包的Slice函数排序
	sort.Slice(keyMemSlice, func(i, j int) bool {
		return keyMemSlice[i].mem > keyMemSlice[j].mem
	})
	// 输出top 10
	fmt.Printf("top max memory usage: \r\n")
	for _, item := range keyMemSlice[:10] {
		fmt.Printf("   key: %s, memory usage:%+vMB\r\n", item.key, item.mem / 1024 /1024)
	}

	// 缩减, 用户标签 label 分类权重表, 保留权重排名前100
	// Z_TXT_VIDEO_LABEL_WEIGHT_
	// Z_SHORT_VIDEO_LABEL_WEIGHT_

	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "11.0.1.152:6380", 
		Password: "",
		DB: 3,
	})
	for _, item:= range keyMemSlice {

		find_txt := strings.Contains(item.key, "Z_TXT_VIDEO_LABEL_WEIGHT_")
		find_video := strings.Contains(item.key, "Z_SHORT_VIDEO_LABEL_WEIGHT_")

		if !find_txt && !find_video {
			continue
		}
		// zcard
		count_zcard, err_zcard := client.ZCard(ctx, item.key).Result()
		if err_zcard != redis.Nil {
			continue
		}
		if count_zcard < 100 {
			continue
		}
		result, err_rem_txt := client.ZRemRangeByRank(ctx, item.key, 0, -101).Result()
		if err_rem_txt != redis.Nil {
			fmt.Printf("ZRemRangeByRank failed, key:%s, err:%+v\r\n", item.key, err_rem_txt)
		}
		fmt.Printf("ZRemRangeByRank key succ, key:%s, ret:%+v\r\n", item.key, result)
	}
}