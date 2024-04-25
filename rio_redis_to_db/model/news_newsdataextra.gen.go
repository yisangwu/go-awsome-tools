// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameNewsNewsdataextra = "news_newsdataextra"

// NewsNewsdataextra mapped from table <news_newsdataextra>
type NewsNewsdataextra struct {
	ID            int64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true" json:"id"`
	ContentType   int32     `gorm:"column:content_type;type:tinyint unsigned;not null;default:1;comment:内容类型 1:图文 2:视频 3:短视频" json:"content_type"`                                                // 内容类型 1:图文 2:视频 3:短视频
	CategoryID    int32     `gorm:"column:category_id;type:int;not null;comment:归属大类" json:"category_id"`                                                                                         // 归属大类
	CrawlerID     int32     `gorm:"column:crawler_id;type:int unsigned;not null;comment:爬虫数据表ID(勿删)" json:"crawler_id"`                                                                           // 爬虫数据表ID(勿删)
	Click         int32     `gorm:"column:click;type:int unsigned;comment:点击量" json:"click"`                                                                                                      // 点击量
	Like          int32     `gorm:"column:like;type:int unsigned;comment:点赞量" json:"like"`                                                                                                        // 点赞量
	Collect       int32     `gorm:"column:collect;type:int unsigned;comment:收藏量" json:"collect"`                                                                                                  // 收藏量
	Report        int32     `gorm:"column:report;type:int unsigned;comment:举报量" json:"report"`                                                                                                    // 举报量
	Comment       int32     `gorm:"column:comment;type:int unsigned;comment:总评论量（不区分用户）" json:"comment"`                                                                                          // 总评论量（不区分用户）
	CommentUsers  int32     `gorm:"column:comment_users;type:int unsigned;not null;comment:评论人数" json:"comment_users"`                                                                            // 评论人数
	Share         int32     `gorm:"column:share;type:int unsigned;comment:分享量" json:"share"`                                                                                                      // 分享量
	Play          int32     `gorm:"column:play;type:int unsigned;comment:观看量(短视频专用字段)" json:"play"`                                                                                               // 观看量(短视频专用字段)
	PlayEnd       int32     `gorm:"column:play_end;type:int unsigned;comment:完播量(短视频专用字段)" json:"play_end"`                                                                                       // 完播量(短视频专用字段)
	Push          int32     `gorm:"column:push;type:int unsigned;comment:已推未读数" json:"push"`                                                                                                      // 已推未读数
	AddRanLike    int32     `gorm:"column:add_ran_like;type:int unsigned;not null;comment:随机添加点赞量" json:"add_ran_like"`                                                                           // 随机添加点赞量
	PoorQuality   int32     `gorm:"column:poor_quality;type:int unsigned;comment:内容质量差量" json:"poor_quality"`                                                                                     // 内容质量差量
	IsHot         int32     `gorm:"column:is_hot;type:tinyint unsigned;comment:热点资讯: 0 :非热点 1:热点" json:"is_hot"`                                                                                  // 热点资讯: 0 :非热点 1:热点
	VideoDuration float64   `gorm:"column:video_duration;type:decimal(10,2);not null;default:0.00;comment:总时长、秒(短视频专用字段)" json:"video_duration"`                                                  // 总时长、秒(短视频专用字段)
	Status        int32     `gorm:"column:status;type:tinyint;comment:状态 0:正常 1:禁用 2:text/video_file_path为空 3:图文s3数据为空  4:图文html为空 5:video title为空 6:第一帧异常 7:分类异常 8:疑似重复数据 9:主图异常" json:"status"` // 状态 0:正常 1:禁用 2:text/video_file_path为空 3:图文s3数据为空  4:图文html为空 5:video title为空 6:第一帧异常 7:分类异常 8:疑似重复数据 9:主图异常
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"`                                                          // 创建时间
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"update_time"`                                                          // 更新时间
}

// TableName NewsNewsdataextra's table name
func (*NewsNewsdataextra) TableName() string {
	return TableNameNewsNewsdataextra
}
