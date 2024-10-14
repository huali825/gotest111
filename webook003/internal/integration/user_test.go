package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"goworkwebook/webook003/internal/integration/startup"
	"goworkwebook/webook003/internal/web"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestUserHandler_SendSMSCode(t *testing.T) {
	server := startup.InitWebServer()
	rdb := startup.InitRedis()

	testCases := []struct {
		name string
		// before
		before func(t *testing.T)
		after  func(t *testing.T)

		phone string

		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功的用例",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				//defer cancel()
				//key := "phone_code:login:15212345678"
				//code, err := rdb.Get(ctx, key).Result()
				//assert.NoError(t, err)
				//assert.True(t, len(code) > 0)
				//dur, err := rdb.TTL(ctx, key).Result()
				//assert.NoError(t, err)
				//assert.True(t, dur > time.Minute*9+time.Second+50)
				//err = rdb.Del(ctx, key).Err()
				//assert.NoError(t, err)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 你要清理数据
				// "phone_code:%s:%s"
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				cancel()
				assert.NoError(t, err)
				// 你的验证码是 6 位
				assert.True(t, len(val) == 6)
			},
			phone:    "15212345678",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},

		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				//defer cancel()
				//key := "phone_code:login:15212345678"
				//err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second*50).Err()
				//assert.NoError(t, err)
				// 这个手机号码，已经有一个验证码了
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:15212345678", "123456",
					time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				//defer cancel()
				//key := "phone_code:login:15212345678"
				//code, err := rdb.GetDel(ctx, key).Result()
				//assert.NoError(t, err)
				//assert.Equal(t, "123456", code)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 你要清理数据
				// "phone_code:%s:%s"
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				cancel()
				assert.NoError(t, err)
				// 你的验证码是 6 位,没有被覆盖，还是123456
				assert.Equal(t, "123456", val)
			},
			phone:    "15212345678",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 4,
				Msg:  "短信发送太频繁，请稍后再试",
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				err := rdb.Set(ctx, key, "123456", 0).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			phone:    "15212345678",
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "未输入手机号码",
			before: func(t *testing.T) {

			},
			after:    func(t *testing.T) {},
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 4,
				Msg:  "请输入手机号码",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			// 准备Req和记录的 recorder
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone": "%s"}`, tc.phone))))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)
			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}

			var res web.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
		})
	}

}
