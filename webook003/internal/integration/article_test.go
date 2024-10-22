package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"goworkwebook/webook003/internal/integration/startup"
	"goworkwebook/webook003/internal/repository/dao"
	ijwt "goworkwebook/webook003/internal/web/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ArticleTestSuiteITF interface {
	SetupSuite()
	TearDownTest()

	TestEdit()
	//TestPublish()
	TestABC()
}

type ArticleHandlerSuite struct {
	suite.Suite
	db     *gorm.DB
	server *gin.Engine
}

// 运行全部
func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}

// SetupSuite函数用于初始化测试套件
func (s *ArticleHandlerSuite) SetupSuite() {
	//s.server = startup.InitWebServer()

	// 创建一个gin的默认实例
	server := gin.Default()
	// 使用中间件，设置用户信息
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user", ijwt.UserClaims{
			Uid: 123,
		})
	})

	// 初始化数据库
	s.db = startup.InitDB2()
	// 初始化文章处理器
	hdl := startup.InitArticleHandler()
	// 注册路由
	hdl.RegisterRoutes(server)
	// 将服务器赋值给s.server
	s.server = server
}

// TearDownTest函数用于在测试结束后清理数据库
func (s *ArticleHandlerSuite) TearDownTest() {
	// 清空articles表
	err := s.db.Exec("truncate table `articles`").Error
	// 断言没有错误
	assert.NoError(s.T(), err)
	// 清空published_articles表
	err = s.db.Exec("truncate table `published_articles`").Error
	// 断言没有错误
	assert.NoError(s.T(), err)
}

func (s *ArticleHandlerSuite) TestABC() {
	s.T().Log("hello 这是测试套件")
}

func (s *ArticleHandlerSuite) TestEdit() {
	t := s.T()
	testCase := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		// 前端传过来，肯定是一个 JSON
		art Article

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "编辑文章",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				// 定义一个dao.Article类型的变量art
				var art dao.Article
				// 在数据库中查找author_id为123的文章，并将结果赋值给art
				err := s.db.Where("author_id=?", 123).
					First(&art).Error
				// 断言err为nil，即查找成功
				assert.NoError(t, err)
				// 断言art的创建时间大于0
				assert.True(t, art.Ctime > 0)
				// 断言art的更新时间大于0
				assert.True(t, art.Utime > 0)
				// 断言art的id大于0
				assert.True(t, art.Id > 0)
				// 断言art的标题为"我的标题"
				assert.Equal(t, "我的标题", art.Title)
				// 断言art的内容为"我的内容"
				assert.Equal(t, "我的内容", art.Content)
				// 断言art的作者id为123
				assert.Equal(t, int64(123), art.AuthorId)
			},
			art: Article{
				//Id:      0,
				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				//Msg:  "ok",
				Data: 1,
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			//构造请求
			//执行
			//断言

			reqBody, err := json.Marshal(tc.art)
			// 准备Req和记录的 recorder
			req, err := http.NewRequest(http.MethodPost,
				"/articles/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			// 执行
			s.server.ServeHTTP(recorder, req)
			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}

			var res Result[int64]
			//var res web.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

// 定义一个泛型结构体Result，用于返回结果
type Result[T any] struct {
	Code int    `json:"code"` // 返回码
	Msg  string `json:"msg"`  // 返回信息
	Data T      `json:"data"` // 返回数据
}

// 定义一个文章结构体Article
type Article struct {
	Id      int64  `json:"id"`      // 文章ID
	Title   string `json:"title"`   // 文章标题
	Content string `json:"content"` // 文章内容
}
