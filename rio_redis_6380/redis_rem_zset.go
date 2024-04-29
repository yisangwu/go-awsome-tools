package main

// 1) "RRRR_Z_TXT_VIDEO_CATEGORY_WEIGHT_1888721_0"
// 2) "RRRR_Z_SHORT_VIDEO_LABEL_WEIGHT_1888721_3"
// 3) "RRRR_Z_SHORT_VIDEO_CATEGORY_WEIGHT_1888721_3"
// 4) "RRRR_Z_TXT_VIDEO_LABEL_WEIGHT_1888721_0"

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

func main() {
	key_pre := "RRRR_Z_"
	ctx := context.Background()
	rds := redis.NewClient(&redis.Options{
		Addr:     "11.0.1.152:6380",
		Password: "",
		DB:       3,
	})
	iter := rds.Scan(ctx, 0, key_pre+"*", 0).Iterator()
	if iter == nil {
		fmt.Printf("client scan not found\r\n")
		return
	}
	for iter.Next(ctx) {
		key := iter.Val()
		if key == "" {
			continue
		}
		find_txt := strings.Contains(key, "Z_TXT_VIDEO_LABEL_WEIGHT_")
		find_video := strings.Contains(key, "Z_SHORT_VIDEO_LABEL_WEIGHT_")

		if !find_txt && !find_video {
			continue
		}
		// if !strings.Contains(key, "1888721") {
		// 	continue
		// }
		// zcard
		count_zcard, err_zcard := rds.ZCard(ctx, key).Result()
		if err_zcard != nil {
			fmt.Printf("ZCard failed, key:%s, err:%+v\r\n", key, err_zcard)
			continue
		}
		if count_zcard < 100 {
			continue
		}
		result, err_rem_txt := rds.ZRemRangeByRank(ctx, key, 0, -101).Result()
		if err_rem_txt != nil {
			fmt.Printf("ZRemRangeByRank failed, key:%s, err:%+v\r\n", key, err_rem_txt)
			continue
		}
		fmt.Printf("ZRemRangeByRank key succ, key:%s, result:%+v\r\n", key, result)
	}
}
