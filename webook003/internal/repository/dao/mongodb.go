package dao

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MongoDBArticleDAO 结构体，用于存储 MongoDB 数据库中的文章信息
type MongoDBArticleDAO struct {
	node    *snowflake.Node   // 雪花算法生成的节点
	col     *mongo.Collection // 存储文章信息的集合
	liveCol *mongo.Collection // 存储实时文章信息的集合
}

// Insert 在MongoDBArticleDAO结构体中定义Insert方法，用于向数据库中插入一条Article记录
func (m *MongoDBArticleDAO) Insert(ctx context.Context, art IsDaoArticle) (int64, error) {
	// 获取当前时间戳
	now := time.Now().UnixMilli()
	// 将当前时间戳赋值给Article的创建时间和更新时间
	art.Ctime = now
	art.Utime = now
	// 生成一个唯一的ID  雪花算法生成的节点
	art.Id = m.node.Generate().Int64()
	// 将Article记录插入到数据库中
	_, err := m.col.InsertOne(ctx, &art)
	// 返回生成的ID和可能出现的错误
	return art.Id, err
}

func (m *MongoDBArticleDAO) UpdateById(ctx context.Context, art IsDaoArticle) error {
	now := time.Now().UnixMilli()
	filter := bson.D{bson.E{"id", art.Id},
		bson.E{"author_id", art.AuthorId}}
	set := bson.D{bson.E{"$set", bson.M{
		"title":   art.Title,
		"content": art.Content,
		"status":  art.Status,
		"utime":   now,
	}}}
	res, err := m.col.UpdateOne(ctx, filter, set)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		// 创作者不对，说明有人在瞎搞
		return errors.New("ID 不对或者创作者不对")
	}
	return nil
}

func (m *MongoDBArticleDAO) Sync(ctx context.Context, art IsDaoArticle) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	now := time.Now().UnixMilli()
	art.Utime = now
	// liveCol 是 INSERT or Update 语义
	filter := bson.D{bson.E{"id", art.Id},
		bson.E{"author_id", art.AuthorId}}
	set := bson.D{bson.E{"$set", art},
		bson.E{"$setOnInsert",
			bson.D{bson.E{"ctime", now}}}}
	_, err = m.liveCol.UpdateOne(ctx,
		filter, set,
		options.Update().SetUpsert(true))
	return id, err
}

func (m *MongoDBArticleDAO) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	filter := bson.D{bson.E{Key: "id", Value: id},
		bson.E{Key: "author_id", Value: uid}}
	sets := bson.D{bson.E{Key: "$set",
		Value: bson.D{bson.E{Key: "status", Value: status}}}}
	res, err := m.col.UpdateOne(ctx, filter, sets)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return errors.New("ID 不对或者创作者不对")
	}
	_, err = m.liveCol.UpdateOne(ctx, filter, sets)
	return err
}

func NewMongoDBArticleDAO(mdb *mongo.Database, node *snowflake.Node) *MongoDBArticleDAO {
	return &MongoDBArticleDAO{
		node:    node,
		liveCol: mdb.Collection("published_articles"),
		col:     mdb.Collection("articles"),
	}
}
