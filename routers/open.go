package routers

import (
	"envelope_rain_group10/logger"
	"envelope_rain_group10/nsqclient"
	redisClient "envelope_rain_group10/redisclient"
	"fmt"
	"go.uber.org/zap"
	"strconv"

	"github.com/gin-gonic/gin"
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

	//值可以设为星号,也可以指定具体主机地址,可设置多个地址用逗号隔开,设为指定主机地址第三项才有效
	c.Header("Access-Control-Allow-Origin", "*")
	//允许请求头修改的类容
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	//允许使用cookie
	c.Header("Access-Control-Allow-Credentials", "true")

	uidString, _ := c.GetPostForm("uid")
	envelopeString, _ := c.GetPostForm("envelope_id")
	logger.Logger.Info("open envelope by", zap.String("uid", uidString))
	uid, _ := strconv.ParseInt(uidString, 10, 64)
	envelopeId, _ := strconv.ParseInt(envelopeString, 10, 64)

	//用户红包列表中有无此envelope_id，redis需要返回一个用户红包的数组func GetEnvelopes(uid int64) []int64
	redPackerList, _ := redisClient.RedisClient.GetUserRedPackerList(uid)
	redPackerListInt := String2Int(redPackerList)

	flag := false
	//遍历数组中每个envelope_id，检查有没有和请求对应的envelope_id。
	for _, val := range redPackerListInt {
		if val == envelopeId {
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
	opened, err := redisClient.RedisClient.RedPacketOpened(envelopeId)
	if err != nil {
		logger.Logger.Error(err.Error())
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
	value, err := redisClient.RedisClient.GetRedPacketMoney(envelopeId)
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	//失效钱包列表缓存
	err = redisClient.RedisClient.MakeWalletCacheInvalid(uid)
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	//修改bitmap数组状态
	err = redisClient.RedisClient.OpenRedPacketInRedisBitMap(envelopeId)
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	//更新数据库opened状态
	nsqclient.ProduceMessage("UpdateStateByEidAndUid", fmt.Sprintf("%d,%d", envelopeId, uid))

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"value": value,
		},
	})
	return
}
