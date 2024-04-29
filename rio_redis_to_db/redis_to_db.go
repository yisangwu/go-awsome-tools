package main

import (
	"context"
	"fmt"
	"go-awsome-tools/rio_redis_to_db/model"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
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
	Host                         string `json:"host"`
	Port                         int    `json:"port"`
	Db                           int    `json:"db"`
	KeySourceTimestampTxt        string `json:"key_source_timestamp_txt"`
	KeyLikeScoreTxt              string `json:"key_like_score_txt"`
	KeySourceTimestampShortVideo string `json:"key_source_timestamp_short_video"`
	KeyLikeScoreShortVideo       string `json:"key_like_score_short_video"`
	KeySourceTimestampVideo      string `json:"key_source_timestamp_video"`
	KeyLikeScoreVideo            string `json:"key_like_score_video"`
}

// GetRedisConf 获取Redis配置
func GetRedisConf() RedisConfStruct {
	return RedisConfStruct{
		Host:                         strings.TrimSpace(v.GetString("redis.host")),
		Port:                         v.Get("redis.port").(int),
		Db:                           v.Get("redis.db").(int),
		KeySourceTimestampTxt:        strings.TrimSpace(v.GetString("redis.key_source_timestamp_txt")),
		KeyLikeScoreTxt:              strings.TrimSpace(v.GetString("redis.key_like_score_txt")),
		KeySourceTimestampShortVideo: strings.TrimSpace(v.GetString("redis.key_source_timestamp_short_video")),
		KeyLikeScoreShortVideo:       strings.TrimSpace(v.GetString("redis.key_like_score_short_video")),
		KeySourceTimestampVideo:      strings.TrimSpace(v.GetString("redis.key_source_timestamp_video")),
		KeyLikeScoreVideo:            strings.TrimSpace(v.GetString("redis.key_like_score_video")),
	}
}

type MysqlConfStruct struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	User           string `json:"user"`
	Password       string `json:"password"`
	Db             string `json:"db"`
	TableNewsLove  string `json:"table_news_love"`
	TableNewsExtra string `json:"table_news_extra"`
}

// GetMysqlConf 获取Redis配置
func GetMysqlConf() MysqlConfStruct {

	return MysqlConfStruct{
		Host:           strings.TrimSpace(v.GetString("mysql.host")),
		Port:           v.Get("mysql.port").(int),
		User:           strings.TrimSpace(v.GetString("mysql.user")),
		Password:       strings.TrimSpace(v.GetString("mysql.password")),
		Db:             strings.TrimSpace(v.GetString("mysql.db")),
		TableNewsLove:  strings.TrimSpace(v.GetString("mysql.table_news_love")),
		TableNewsExtra: strings.TrimSpace(v.GetString("mysql.table_news_extra")),
	}
}

// GenTableModel 从model 生成 orm
func GenTableModel() {
	mysqlConf := GetMysqlConf()
	dsn := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConf.User, mysqlConf.Password, mysqlConf.Host, mysqlConf.Port, mysqlConf.Db,
	)
	fmt.Printf("check dsn:%s\r\n", dsn)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error! Connect database failed, err:%+v\r\n", err)
		os.Exit(1)
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:           "models",
		OutFile:           "go",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldWithTypeTag:  true,
		FieldWithIndexTag: false,
	})

	g.UseDB(conn)
	for _, tabel := range []string{mysqlConf.TableNewsLove, mysqlConf.TableNewsExtra} {
		g.GenerateModel(tabel)
	}
	g.Execute()
}

func main() {

	startTimestamp := time.Now().Unix()
	defer func() {
		endTimestamp := time.Now().Unix()
		fmt.Printf("exec end, timestamp:%s, cost:%+v\n", time.Now(), endTimestamp-startTimestamp)
	}()
	fmt.Printf("exec start,timestamp:%s\r\n", time.Now())

	// mysql
	mysqlConf := GetMysqlConf()
	dsn := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConf.User, mysqlConf.Password, mysqlConf.Host, mysqlConf.Port, mysqlConf.Db,
	)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error! Connect database failed, err:%+v\r\n", err)
		os.Exit(1)
	}
	// redis
	redisConf := GetRedisConf()
	ctx := context.Background()
	rds := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password: "",
		DB:       redisConf.Db,
	})

	// languages := []redis.Z{
	//     {Score: 123, Member: 89293},
	// }
	// fmt.Println(rds.ZAdd(ctx, redisConf.KeySourceTimestampTxt, languages...).Result())
	// fmt.Println(rds.ZScore(ctx, redisConf.KeySourceTimestampTxt, fmt.Sprintf("%+v",89293)).Result())

	// 迭代extra表记录，查询redis，获取 like_score，source_timestamp
	// 迭代获取数据
	pageSize := 1000 // 每页数据条数
	page := 1
	for {
		var newsExtra []model.NewsNewsdataextra
		// 分页查询
		offset := (page - 1) * pageSize
		if err := conn.Limit(pageSize).Offset(offset).Find(&newsExtra).Error; err != nil {
			fmt.Printf("Error! Mysql find data failed, err:%+v\r\n", err)
			os.Exit(1)
		}
		if len(newsExtra) == 0 {
			break // 如果当前页没有数据了，退出循环
		}

		// 处理当前页的数据
		for _, items := range newsExtra {

			batchNewsLove := []*model.NewsLove{}
			//  fmt.Printf("ID: %+v, ContentType: %+v, CreateTime: %+v\n", items.ID, items.ContentType, items.CreateTime)
			key_timestamp := ""
			key_like_score := ""

			if items.ID != 89293 {
				continue
			}

			switch items.ContentType {
			case 1: // 图文
				key_timestamp = redisConf.KeySourceTimestampTxt
				key_like_score = redisConf.KeyLikeScoreTxt
			case 2: // 视频
				key_timestamp = redisConf.KeySourceTimestampVideo
				key_like_score = redisConf.KeyLikeScoreVideo
			case 3: // 短视频
				key_timestamp = redisConf.KeySourceTimestampShortVideo
				key_like_score = redisConf.KeyLikeScoreShortVideo
			default:
				fmt.Printf("unknown content_type, id:%+v, content_type:%+v\r\n", items.ID, items.ContentType)

			}
			if key_timestamp == "" || key_like_score == "" {
				continue
			}

			// 取like_score
			likeSocre, _ := rds.ZScore(ctx, key_like_score, fmt.Sprintf("%+v", items.ID)).Result()
			// 取source_timestamp
			sourceTimstamp, _ := rds.ZScore(ctx, key_timestamp, fmt.Sprintf("%+v", items.ID)).Result()

			batchNewsLove = append(batchNewsLove, &model.NewsLove{
				NewsID:          int32(items.ID),
				ContentType:     items.ContentType,
				CategoryID:      items.CategoryID,
				LoveScore:       int32(likeSocre),
				SourceTimestamp: int32(sourceTimstamp),
			})
			result := conn.Create(batchNewsLove)
			if result.Error != nil {
				fmt.Printf("mysql conn.Createunknown content_type, err:%+v\r\n", result.Error)
				continue
			}
			fmt.Printf("mysql Create succ, RowsAffected:%+v\r\n", result.RowsAffected)
		}
		// 进入下一页
		page++
	}

}

// SELECT * FROM `news_newsdataextra LIMIT 1000 OFFSET 343000
// /data-local/go-awsome-tools/rio_redis_to_db/redis_to_db.go:156 SLOW SQL >= 200ms
// Affected:1
// Rush-web-frontend
