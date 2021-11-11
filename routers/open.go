package routers

import (
	redisClient "envelope_rain_group10/redisclient"
	"envelope_rain_group10/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
)

func LoadOpen(e *gin.Engine) {
	e.POST("/open", OpenHandler)
}

func String2Int(strArr []string) []int64 {
	res := make([]int64, len(strArr))

	for index, val := range strArr {
		res[index], _ = strconv.ParseInt(val, 10, 64)
	}
	return res
}

func OpenHandler(c *gin.Context) {
	uid, _ := c.GetPostForm("uid")
	envelope_id, _ := c.GetPostForm("envelope_id")

	logs.Printf("envelope %s opened by %s", envelope_id, uid)

	int_uid, _ := strconv.ParseInt(uid, 10, 64)
	int_envelope_id, _ := strconv.ParseInt(envelope_id, 10, 64)

	//用户红包列表中有无此envelope_id，redis需要返回一个用户红包的数组func GetEnvelopes(uid int64) []int64
	redPackerList, _ := redisClient.RedisClient.GetUserRedPackerList(int_uid)
	redPackerListInt := String2Int(redPackerList)

	flag := false
	//遍历数组中每个envelope_id，检查有没有和请求对应的envelope_id。
	for _, val := range redPackerListInt {
		if val == int_envelope_id {
			flag = true
		}
	}
	if flag == false {
		c.JSON(200, gin.H{
			"code": -2,
			"msg":  "It 's not your red packet!",
			"data": gin.H{
				"uid": uid,
			},
		})
		return
	}

	//检查红包是否打开，redis需要返回一个用户红包的数组func HasOpened(envelope_id int64) bool，
	opened, err := redisClient.RedisClient.RedPacketOpened(int_envelope_id)
	if err != nil {
		logs.Println(err)
	}
	//开过，返回提示
	if opened == true {
		c.JSON(200, gin.H{
			"code": -3,
			"msg":  "The red envelope has already been opened",
			"data": gin.H{
				"opened": opened,
			},
		})
		return
	}

	//没开过，用红包id查money并返回,redis需要提供func GetValueByUid(uid int64) int64，
	value, err := redisClient.RedisClient.GetRedPacketMoney(int_envelope_id)
	if err != nil {
		logs.Println(err)
	}

	//失效钱包列表缓存
	err = redisClient.RedisClient.MakeWalletCacheInvalid(int_uid)
	if err != nil {
		logs.Println(err)
	}

	//修改bitmap数组状态
	err = redisClient.RedisClient.OpenRedPacketInRedisBitMap(int_envelope_id)
	if err != nil {
		logs.Println(err)
	}

	//更新数据库opened状态 func UpdateState(envelope_id int64)
	sql.UpdateStateByEidAndUid(int_envelope_id, int_uid) //更新红包opened

	sql.UpdateState(int_envelope_id)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"value": value,
		},
	})
	return
}
