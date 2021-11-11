package main

import (
	"envelope_rain_group10/allocation"
	"envelope_rain_group10/logger"
	"envelope_rain_group10/ratelimit"
	redisClient "envelope_rain_group10/redisclient"
	"envelope_rain_group10/routers"
	"envelope_rain_group10/sql"
	"envelope_rain_group10/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	logger.InitLogger()

	//创建数据库连接
	db, err := sql.InitDB()
	if err != nil {
		logger.Logger.Error("数据库链接失败")
	}
	defer db.Close()

	//读取配置文件（zyz）
	//初始化六个变量
	//每个用户最多可抢到的次数MaxTimes、抢到的概率Probability、总金额TotalMoney、总个数TotalNum、每个红包的金额范围[MaxMoney, MinMoney]
	//除probability类型为（Float64）外，其余变量均为（Int64）,通过utils.MaxTimes调用这些变量
	utils.InitConfigs("./config.json")
	redisClient.InitRedisClient()


	//算法生成红包的id和value的对应表
	//初始化redis中envelop_id 和 value的对应表
	//redis需要提供函数func InitEnvelopeValue(values []int)
	initTag, _ := redisClient.RedisClient.ShouldInit()
	//如果为tru表明setnx成功，那么本次应该初始化，否则表明之前有客户端启动的时候初始化过了，本次不应该初始化
	if initTag{
		a := allocation.NewAllocation(int(utils.TotalMoney), int(utils.TotalNum), int(utils.MinMoney), int(utils.MaxMoney))
		//fmt.Printf("%#v\n", a)
		//这地方是不是要改回来
		values := a.AllocateMoney(int(utils.TotalNum))
		s := make([]interface{}, len(values))
		for i, v := range values {
			s[i] = v
		}
		redisClient.RedisClient.InitRedPacket(s)
	}

	redisClient.RedisClient.InitCurrentRedPacketID()
	r := gin.Default()

	//添加限流中间件,每秒流量从配置中获取
	r.Use(ratelimit.RateLimitMiddleware(time.Second, utils.QpsLimit, utils.QpsLimit))
	routers.LoadSnatch(r)
	routers.LoadOpen(r)
	routers.LoadWalletList(r)
	r.Run()
}
