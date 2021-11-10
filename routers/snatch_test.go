package routers

import (
	"envelope_rain_group10/allocation"
	redisClient "envelope_rain_group10/redisclient"
	"envelope_rain_group10/sql"
	"envelope_rain_group10/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/snatch", SnatchHandler)
	router.POST("/open", OpenHandler)
	router.POST("/get_wallet_list", WalletListHandler)
	return router
}

func TestRouters_Snatch(t *testing.T)  {
	tests := []struct {
		name   string
		param  string
		expect string
	}{
		{"test1", `{"uid":"777"}`, `{"code":0,"data":{"cur_count":1,"envelope_id":1,"max_count":5},"msg":"success"}`},
		{"test2", `{"uid":"777"}`, `{"code":0,"data":{"cur_count":2,"envelope_id":2,"max_count":5},"msg":"success"}`},
		{"test3", `{"uid":"888"}`, `{"code":0,"data":{"cur_count":1,"envelope_id":3,"max_count":5},"msg":"success"}`},
		{"test4", `{"uid":"777"}`, `{"code":0,"data":{"cur_count":3,"envelope_id":4,"max_count":5},"msg":"success"}`},
		{"test5", `{"uid":"777"}`, `{"code":0,"data":{"cur_count":4,"envelope_id":5,"max_count":5},"msg":"success"}`},
		{"test6", `{"uid":"777"}`, `{"code":0,"data":{"cur_count":5,"envelope_id":6,"max_count":5},"msg":"success"}`},
		{"test7", `{"uid":"777"}`, `{"code":-,"data":{"cur_count":-1,"max_count":5},"msg":"fail! snatch too much","data"}`},

	}

	r := SetupRouter()
	db, err := sql.InitDB()
	if err != nil {
		log.Println("database connection failure")
	}
	defer db.Close()
	redisClient.InitRedisClient()
	utils.InitConfigs("./config-test.json")
	redisClient.InitRedisClient()

	//算法生成红包的id和value的对应表
	//初始化redis中envelop_id 和 value的对应表
	//redis需要提供函数func InitEnvelopeValue(values []int)
	a := allocation.NewAllocation(int(utils.TotalMoney), int(utils.TotalNum) , int(utils.MinMoney), int(utils.MaxMoney))
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
				"/snatch",                    // 请求URL
				strings.NewReader(tt.param), // 请求参数
			)

			// mock一个响应记录器
			w := httptest.NewRecorder()

			// 让server端处理mock请求并记录返回的响应内容
			r.ServeHTTP(w, req)

			// 校验状态码是否符合预期
			assert.Equal(t, http.StatusOK, w.Code)

			// 解析并检验响应内容是否复合预期
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, w.Body.String())
		})
	}
	
}
