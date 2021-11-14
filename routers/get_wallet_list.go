package routers

import (
	"encoding/json"
	"envelope_rain_group10/logger"
	"envelope_rain_group10/model"
	redisClient "envelope_rain_group10/redisclient"
	"go.uber.org/zap"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LoadWalletList(e *gin.Engine) {
	e.POST("/get_wallet_list", WalletListHandler)
}

func WalletListHandler(c *gin.Context) {

	//值可以设为星号,也可以指定具体主机地址,可设置多个地址用逗号隔开,设为指定主机地址第三项才有效
	c.Header("Access-Control-Allow-Origin", "*")
	//允许请求头修改的类容
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	//允许使用cookie
	c.Header("Access-Control-Allow-Credentials", "true")

	uidString, _ := c.GetPostForm("uid")
	logger.Logger.Info("query wallet list", zap.String("uid", uidString))
	uid, _ := strconv.ParseInt(uidString, 10, 64)
	walletList, err2 := redisClient.RedisClient.GetUserWalletInRedis(uid)
	//存在钱包列表缓存
	if walletList != "" {
		c.String(http.StatusOK, walletList)
		return
	}

	exist, err2 := redisClient.RedisClient.ExistUser(uid)
	if err2==nil{
		//如果用户不存在，则直接返回
		if exist!=true{
			c.JSON(200,gin.H{
				"code": -1,
				"msg":  "user not exist",
				"data": gin.H{
					"amount":        0,
					"envelope_list": gin.H{},
				},
			})
			return
		}

	}


	redPackageList, err := redisClient.RedisClient.GetUserRedPackerList(uid)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	redPackageListInt := String2Int(redPackageList)

	envelopeTimeList, err := redisClient.RedisClient.GetUserRedPackerTimeList(uid)

	envelopeTimeListInt := String2Int(envelopeTimeList)

	var envelopes = make([]*model.Envelope, len(redPackageListInt))
	for i := range redPackageListInt {
		opened, _ := redisClient.RedisClient.RedPacketOpened(redPackageListInt[i])
		value, _ := redisClient.RedisClient.GetRedPacketMoney(redPackageListInt[i])
		envelopes[i] = &model.Envelope{ID: redPackageListInt[i], Opened: opened, Value: value, SnatchTime: envelopeTimeListInt[i]}
	}

	//按照时间排序
	sort.SliceStable(envelopes, func(i, j int) bool {
		return envelopes[i].SnatchTime < envelopes[j].SnatchTime
	})

	var amount int64 = 0
	var myArray []map[string]interface{}

	for i := 0; i < len(envelopes); i++ {
		var curEnvelope = make(map[string]interface{})
		if envelopes[i].Opened == false {
			curEnvelope["envelope_id"] = envelopes[i].ID
			curEnvelope["opened"] = false
			curEnvelope["snatch_time"] = envelopes[i].SnatchTime
		} else {
			curEnvelope["envelope_id"] = envelopes[i].ID
			curEnvelope["value"] = envelopes[i].Value
			curEnvelope["opened"] = true
			curEnvelope["snatch_time"] = envelopes[i].SnatchTime
			amount = amount + envelopes[i].Value
		}
		myArray = append(myArray, curEnvelope)
	}

	walletListByte, err2 := json.Marshal(myArray)
	if err2 != nil {
		logger.Logger.Error("json转换错误，放弃将用户钱包列表缓存进数据库，此报错不影响程序运行，请及时检查")
	} else {
		//json转换成功才将钱包列表缓存进数据库
		walletJsonString := `{"code":0,"msg":"success","data":{"amount":` + strconv.FormatInt(amount, 10) + `,"envelope_list":` + string(walletListByte) + `}}`
		err2 = redisClient.RedisClient.AddUserWalletToRedis(uid, walletJsonString, 300)
		if err2 != nil {
			logger.Logger.Error("将用户钱包列表缓存进数据库时出错，此报错不影响程序运行，请及时检查")
		}
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":        amount,
			"envelope_list": myArray,
		},
	})
}
