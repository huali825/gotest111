package mongodb

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

// 测试 MongoDB 连接
func TestMongoDB(t *testing.T) {
	// 创建一个带有超时时间的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建一个命令监控器
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			fmt.Println(evt.Command)
		},
	}
	// 创建一个 MongoDB 客户端选项
	opts := options.Client().
		ApplyURI("mongodb://root:example@localhost:27017/").
		SetMonitor(monitor)

	// 连接到 MongoDB
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)

	// 操作 client
	col := client.Database("webook").
		Collection("articles")

	// 插入一条数据
	insertRes, err := col.InsertOne(ctx, Article{
		Id:       1,
		Title:    "我的标题",
		Content:  "我的内容",
		AuthorId: 123,
	})
	assert.NoError(t, err)
	oid := insertRes.InsertedID.(primitive.ObjectID)
	t.Log("插入ID", oid)

	// 查询数据
	// 定义一个filter，用于查询id为1的数据
	filter := bson.M{
		"id": 1,
	}
	// 在col集合中查找符合filter条件的数据
	findRes := col.FindOne(ctx, filter)
	// 如果没有找到数据，则输出提示信息
	if findRes.Err() == mongo.ErrNoDocuments {
		t.Log("没找到数据")
		// 否则，断言查询结果没有错误
	} else {
		assert.NoError(t, findRes.Err())
		// 定义一个Article类型的变量art
		var art Article
		// 将查询结果解码到art变量中
		err = findRes.Decode(&art)
		// 断言解码过程没有错误
		assert.NoError(t, err)
		// 输出art变量的值
		t.Log(art)
	}

	// 更新一条数据
	updateFilter := bson.D{bson.E{"id", 1}}
	set := bson.D{bson.E{Key: "$set", Value: bson.M{
		"title": "新的标题",
	}}}
	updateOneRes, err := col.UpdateOne(ctx, updateFilter, set)
	assert.NoError(t, err)
	t.Log("更新文档数量", updateOneRes.ModifiedCount)

	// 更新多条数据
	updateManyRes, err := col.UpdateMany(ctx, updateFilter,
		bson.D{bson.E{Key: "$set",
			Value: Article{Content: "新的内容"}}})
	assert.NoError(t, err)
	t.Log("更新文档数量", updateManyRes.ModifiedCount)

	// 删除数据
	deleteFilter := bson.D{bson.E{"id", 1}}
	delRes, err := col.DeleteMany(ctx, deleteFilter)
	assert.NoError(t, err)
	t.Log("删除文档数量", delRes.DeletedCount)
}

// 文章结构体
type Article struct {
	// omitempty : 忽略0值
	Id       int64  `bson:"id,omitempty"`        // 文章ID
	Title    string `bson:"title,omitempty"`     // 文章标题
	Content  string `bson:"content,omitempty"`   // 文章内容
	AuthorId int64  `bson:"author_id,omitempty"` // 作者ID
	Status   uint8  `bson:"status,omitempty"`    // 文章状态
	Ctime    int64  `bson:"ctime,omitempty"`     // 创建时间
	// 更新时间
	Utime int64 `bson:"utime,omitempty"` // 更新时间
}
