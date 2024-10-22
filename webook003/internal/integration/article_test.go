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

func (s *ArticleHandlerSuite) TearDownTest() {
	err := s.db.Exec("truncate table `articles`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("truncate table `published_articles`").Error
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
		{},
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
