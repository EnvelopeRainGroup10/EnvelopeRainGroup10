package routers

import (
	"encoding/json"
	"envelope_rain_group10/logger"
	redisClient "envelope_rain_group10/redisclient"
	"envelope_rain_group10/sql"
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

	uid, _ := c.GetPostForm("uid")
	logger.Logger.Info("query wallet list",zap.String("uid",uid))

	int_uid, _ := strconv.ParseInt(uid, 10, 64)
	walletList, err2 := redisClient.RedisClient.GetUserWalletInRedis(int_uid)
	if err2!=nil{
		logger.Logger.Error("读取用户钱包列表缓存时出现未知错误，此错误不影响程序运行，请及时检查")
	}
	//存在钱包列表缓存
	if walletList!=""{
		c.String(http.StatusOK,walletList)
		return
	}

	//现在主要是针对redis查询
	//需要redis提供用户红包列表(id信息)func GetEnvelopeIdsByUid(uid int64) []int64
	//walletList, err := redisClient.RedisClient.GetUserWalletInRedis(int_uid)
	//if err != nil {
	//	logger.Logger.Error(err.Error())
	//}

	//if walletList != "" {
	//	c.JSON(200, gin.H{
	//		"code": 0,
	//		"msg":  "success",
	//		"data": gin.H{
	//			"amount":        amount,
	//			"envelope_list": myArray,
	//		},
	//	})
	//	return
	//}

	redPackageList, err := redisClient.RedisClient.GetUserRedPackerList(int_uid)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	redPackageListInt := String2Int(redPackageList)

	envelopeTimeList, err := redisClient.RedisClient.GetUserRedPackerTimeList(int_uid)

	envelopeTimeListInt := String2Int(envelopeTimeList)

	var  envelopes []*sql.Envelope = make([]*sql.Envelope, len(redPackageListInt))
	for i, _ := range redPackageListInt {
		opened, _ := redisClient.RedisClient.RedPacketOpened(redPackageListInt[i])
		value, _ := redisClient.RedisClient.GetRedPacketMoney(redPackageListInt[i])
		envelopes[i] = &sql.Envelope{ID: redPackageListInt[i], Opened: opened, Value: value, SnatchTime: envelopeTimeListInt[i]}
	}
	//需要redis提供用户红包列表(time信息)func GetEnvelopeTimesByUid(uid int64) []int64
	//需要redis提供根据envelope_id查询红包钱数，func GetValueByEnvelopeId(envelope_id int64) int64
	//剩下的排序和构造json之前已经有了

	//envelopes, _ := sql.GetAllEnvelopesByUID(int_uid)
	//先按照时间排序
	sort.SliceStable(envelopes, func(i, j int) bool {
		return envelopes[i].SnatchTime < envelopes[j].SnatchTime
	})

	var amount int64 = 0
	var myArray []map[string]interface{}

	for i := 0; i < len(envelopes); i++ {
		var curEnvelope map[string]interface{} = make(map[string]interface{})
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

	wallet_list_byte, err2 := json.Marshal(myArray)
	if err2!=nil{
		logger.Logger.Error("json转换错误，放弃将用户钱包列表缓存进数据库，此报错不影响程序运行，请及时检查")
	}else{
		//json转换成功才将钱包列表缓存进数据库
		wallet_json_string:=`{"code":0,"msg":"success","data":{"amount":`+strconv.FormatInt(amount,10)+`,"envelope_list":`+string(wallet_list_byte)+`}}`
		err2 = redisClient.RedisClient.AddUserWalletToRedis(int_uid, wallet_json_string, 300)
		if err2!=nil{
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
