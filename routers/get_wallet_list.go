package routers

import (
	redisClient "envelope_rain_group10/redisclient"
	"envelope_rain_group10/sql"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
)

func LoadWalletList(e *gin.Engine) {
	e.POST("/get_wallet_list", WalletListHandler)
}

func WalletListHandler(c *gin.Context) {

	uid, _ := c.GetPostForm("uid")
	logs.Printf("query %s 's wallets", uid)

	int_uid, _ := strconv.ParseInt(uid, 10, 64)

	//现在主要是针对redis查询
	//需要redis提供用户红包列表(id信息)func GetEnvelopeIdsByUid(uid int64) []int64
	//walletList, err := redisClient.RedisClient.GetUserWalletInRedis(int_uid)
	//if err != nil {
	//	logs.Println(err)
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
		logs.Println(err)
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

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"amount":        amount,
			"envelope_list": myArray,
		},
	})
}
