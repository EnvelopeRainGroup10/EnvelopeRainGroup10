package routers

import (
	"envelope_rain_group10/allocation"
	redisClient "envelope_rain_group10/redisclient"
	"envelope_rain_group10/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter_Open(t *testing.T) {
	tests := []struct {
		name   string
		param  string
		expect string
	}{
		{"test1", `{"uid":"777","envelope_id":"1"}`, `{"code":0,"data":{"value":30},"msg":"success"}`},
		{"test2", `{"uid":"777""envelope_id":"1",}`, `{"code":-3,"data":{"opened":"true"},"msg":"The red envelope has already been opened"}`},
		{"test3", `{"uid":"888","envelope_id":"1"}`, `{"code":-2,"data":{"uid":777},"msg":"It 's not your red packet!"}`},
	}

	r := SetupRouter()

	redisClient.InitRedisClient()
	utils.InitConfigs("./config-test.json")
	redisClient.InitRedisClient()

	//算法生成红包的id和value的对应表
	//初始化redis中envelop_id 和 value的对应表
	//redis需要提供函数func InitEnvelopeValue(values []int)
	a := allocation.NewAllocation(int(utils.TotalMoney), int(utils.TotalNum), int(utils.MinMoney), int(utils.MaxMoney))
	//fmt.Printf("%#v\n", a)
	values := a.AllocateMoney(1000000)

	s := make([]interface{}, len(values))
	for i, v := range values {
		s[i] = v
	}
	redisClient.RedisClient.InitRedPacket(s)
	redisClient.RedisClient.InitCurrentRedPacketID()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock一个HTTP请求
			req := httptest.NewRequest(
				"POST",                      // 请求方法
				"/open",                     // 请求URL
				strings.NewReader(tt.param), // 请求参数
			)

			// mock一个响应记录器
			w := httptest.NewRecorder()

			// 让server端处理mock请求并记录返回的响应内容
			r.ServeHTTP(w, req)

			// 校验状态码是否符合预期
			assert.Equal(t, http.StatusOK, w.Code)

			// 解析并检验响应内容是否复合预期
			assert.Equal(t, tt.expect, w.Body.String())
		})
	}

}
