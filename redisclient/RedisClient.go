package redisClient

import (
	"bytes"
	"envelope_rain_group10/logger"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

//红包列表在redis中的名字
var redPacketListKeyName="redPacketListKeyName"
//当前可获取的红包id在redis中key的名字
var currentRedPacketId="currentRedPacketId"
//用于记录红包是否开启的bitmap的名字
var redPacketOpenedBitMapKeyName="redPacketOpenedBitMap"
//用于记录用户是否存在的bitMap的key的名字
var userExistedBitMapKeyName="userExistedBitMapKeyName"

type redisClient struct {
	//实际连接的客户端
	 rdb *redis.Client
	 //用于区分的前缀
	 keyPre string
	 //当前红包id
	 currentRedPacketId int64
	 //最大红包个数
	 maxPacketNum int64
	 //每个人最大抢红包数
	 maxGetNum int64
	 //redis中红包列表的key的名字
	redPacketListKeyName string
	 //redis中当前可获取红包id的key的名字
	currentRedPacketIdKeyName string
}


//---------------------------------------初始化链接部分----------------------------------------
//初始化客户端链接
//addr:redisip及端口号如 127.0.0.1:6379
//db哪个redis库，一般输入0即可
//poolSize数据库链接池大小，建议1000
//maxPacketNum：总共发多少红包
//maxGetNum：每个人最多可以抢几个红包
//keyPre：redis中容易的key的前缀，建议每次部署使用不同前缀
func NewRedisClient(addr,password string,db,poolSize,maxPacketNum,maxGetNum int64,keyPre string)( *redisClient,error)  {

	rdb:=redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       int(db),
		PoolSize: int(poolSize),
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		logger.Logger.Error("链接redis失败")
		return nil,err
	}

	return &redisClient{
		rdb: rdb,
		keyPre: keyPre,
		//后面再初始化，现在默认为0
		currentRedPacketId: 0,
		maxPacketNum: maxPacketNum,
		maxGetNum: maxGetNum,
		redPacketListKeyName: splicingString(keyPre,redPacketListKeyName),
		currentRedPacketIdKeyName: splicingString(keyPre,currentRedPacketId),
	},nil

}
//关闭客户端连接
func (rc *redisClient)CloseClient()error{
	err := rc.rdb.Close()
	return err
}

//初始化红包值列表到Redis中，重置当前可抢红包id为0
//返回值表示是否初始化成功
//本方法只应该调用一次，常规情况无需调用
func (rc *redisClient)InitRedPacket(redPacketList []interface{}) (bool,error){

	//把红包金额列表压入缓存中
	//这个要不要存到数据库里面？
	_,err:= rc.rdb.RPush(rc.redPacketListKeyName,redPacketList...).Result()

	if err!=nil{
		logger.Logger.Error("缓存红包列表失败")
		return false,err
	}


	//如果插入成功了，那么就把当前可获取的红包id也给加到redis中
	_, err = rc.rdb.Set(rc.currentRedPacketIdKeyName, 0, 0).Result()

	if err!=nil{
		logger.Logger.Error("初始化可获取红包id失败")
		return false,err
	}

	//初始化红包开启列表
	redPacketOpenedBitMapInited := rc.initRedPacketOpenedBitMap()
	userExistedbitMapInited := rc.initUserExistedBitMap()

	if redPacketOpenedBitMapInited==true && userExistedbitMapInited==true{
		return true,nil
	}else {
		return false,nil
	}



}

//初始化RedisClient中保存的currrntRedPacketId的值
//从Redis中读取这个值
func (rc *redisClient)InitCurrentRedPacketID(){
	id := rc.getCurrentRedPacketID()
	if id!=-1{
		rc.currentRedPacketId=id
	}
}
//-------------------------------------------------------------------------------------------

//获取红包金额
//如果出现错误返回-1
func (rc *redisClient)GetRedPacketMoney(redPacketId int64)(int64,error){
	result, err := rc.rdb.LIndex(rc.redPacketListKeyName, redPacketId - 1).Result()
	if err!=nil{
		logger.Logger.Error("未能成功获取红包金额")
		return -1,err
	}else {
		parseInt, err := strconv.ParseInt(result, 10, 64)
		if err!=nil{
			logger.Logger.Error("红包金额转换为int时失败")
			return -1,err
		}
		return parseInt,nil
	}

}

//这个方法会通过incr竞争redPacket，目的是为了防止同一个user的并发
//如果竞争到的值大于最大能抢的红包数那么代表这次抢夺失败，必须通过decr还原
//否则参与后续竞争，返回值为本次竞争在抢第几个红包
func(rc *redisClient)GetCountWithNextRedPacketByUserId (userId int64)(int64,error)  {
	buffer := bytes.Buffer{}
	buffer.WriteString(rc.keyPre)
	//这个或许可以抽取成配置文件
	buffer.WriteString("countNumOfPacket:")
	buffer.WriteString(strconv.FormatInt(userId,10))
	key:=buffer.String()

	//先做一个get判断，如果已经大于或者等于了，就没必要进去尝试抢了
	res, err := rc.rdb.Get(key).Result()
	if err!=nil{
		//如果是key不存在报错那么不管他，否则报错
		if err!=redis.Nil{
			logger.Logger.Error("用户获取他抢了多少个红包时出错")
		}
	}else {
		count,err:=strconv.ParseInt(res,10,64)
		if err==nil{
			if count>=rc.maxGetNum{
				return -1,nil
			}
		}
	}


	result, err := rc.rdb.Incr(key).Result()
	if err!=nil{
		logger.Logger.Error("用户竞争红包时出现错误")
		return -1,err
	}else {
		if result>rc.maxGetNum{
			_, _ = rc.rdb.Decr(key).Result()
			return -1,nil
		}else {
			return result,nil
		}
	}
}

//将缓存中用户已抢到的红包计数值减1
func (rc *redisClient)ReduceUserGetRedPacketCount(userId int64)error{
	buffer := bytes.Buffer{}
	buffer.WriteString(rc.keyPre)
	//这个或许可以抽取成配置文件
	buffer.WriteString("countNumOfPacket:")
	buffer.WriteString(strconv.FormatInt(userId,10))
	key:=buffer.String()
	_, err := rc.rdb.Decr(key).Result()
	if err!=nil{
		logger.Logger.Error("减少用户已抢到红包计数时出错")
		return err
	}
	return nil

}

func (rc *redisClient)GetRedPacket()(int64,error)  {
	//这个里面由于rc的currentRedPacketId没有做并发保护，确实可能会出现同步问题
	//但是这个问题并不会影响线程获取红包id，短时间内放入一些请求后后面还是能够成功拦截。
	//做同步反而浪费性能
	//这样可以红包发完后请求可以全部拦截掉，不用请求redis
	//先判断一次本地结构体内部存储的当前红包id是否大于可用红包数
	if rc.currentRedPacketId<(rc.maxPacketNum){
		//使用incr来尝试获取一个红包

		result, err := rc.rdb.Incr(rc.currentRedPacketIdKeyName).Result()
		if err!=nil{
			logger.Logger.Error("将红包id自增时发生错误")
			return -1,err
		}
		//抢到的红包id大于红包个数，抢夺无效
		//此时将currentRedPacketId设置为抢回来的值，那么以后就不会再访问redis
		if result>rc.maxPacketNum{
			//没抢到的话就把redis里面的红包计数减回去
			rc.rdb.Decr(rc.currentRedPacketIdKeyName).Result()
			rc.currentRedPacketId=rc.maxPacketNum
			return -1,nil
		}else {
			rc.currentRedPacketId=result
			return result,nil
		}
	}else{
		return -1,nil
	}
}


//向用户的红包列表总添加一项
//缓存未命中从数据库加载到缓存这种事情的都是在外面而不是里面做的
func(rc *redisClient)AddToUserRedPacketList(userId,redPacketId int64)error{
	buffer := bytes.Buffer{}
	buffer.WriteString(rc.keyPre)
	buffer.WriteString("redPacketList:")
	buffer.WriteString(strconv.FormatInt(userId,10))
	key:=buffer.String()
	_, err := rc.rdb.RPush(key, redPacketId).Result()
	if err!=nil{
		logger.Logger.Error("将红包Id插入红包列表的缓存时出错")
		return err
	}
	return nil
}

//向用户获取红包的时间列表中加一项
func(rc *redisClient)AddToUserRedPacketTimeList(userId,timeStamp int64)error{
	buffer := bytes.Buffer{}
	buffer.WriteString(rc.keyPre)
	buffer.WriteString("redPacketTimeList:")
	buffer.WriteString(strconv.FormatInt(userId,10))
	key:=buffer.String()
	_, err := rc.rdb.RPush(key, timeStamp).Result()
	if err!=nil{
		logger.Logger.Error("将红包时间插入红包时间列表的缓存时出错")
		return err
	}
	return nil
}

//读取用户的红包列表
//返回值是id的list,注意返回的是[]string
func(rc *redisClient)GetUserRedPackerList(userId int64)([]string,error){
	buffer := bytes.Buffer{}
	buffer.WriteString(rc.keyPre)
	buffer.WriteString("redPacketList:")
	buffer.WriteString(strconv.FormatInt(userId,10))
	key:=buffer.String()

	result, err := rc.rdb.LRange(key, 0, -1).Result()
	if err!=nil{
		
		logger.Logger.Error("获取用户红包列表时出错")
		return nil, err
	}

	return result,nil

}

//读取用户的红包时间列表
//返回值是id的list,注意返回的是[]string
func(rc *redisClient)GetUserRedPackerTimeList(userId int64)([]string,error){
	buffer := bytes.Buffer{}
	buffer.WriteString(rc.keyPre)
	buffer.WriteString("redPacketTimeList:")
	buffer.WriteString(strconv.FormatInt(userId,10))
	key:=buffer.String()

	result, err := rc.rdb.LRange(key, 0, -1).Result()
	if err!=nil{
		logger.Logger.Error("获取用户红包时间列表时出错")
		return nil, err
	}

	return result,nil

}

//查询用户是否存在
func (rc *redisClient)ExistUser(userId int64)(bool,error)  {
	key:=splicingString(rc.keyPre,userExistedBitMapKeyName)
	result, err := rc.rdb.GetBit(key, userId).Result()
	if err!=nil{
		logger.Logger.Error("获取用户是否存在信息时失败")
		return false,err
	}
	if result==1{
		return true,nil
	}else {
		return false,nil
	}
}

func (rc *redisClient)CreateUserInRedis(userId int64)error{
	key:=splicingString(rc.keyPre,userExistedBitMapKeyName)
	_, err := rc.rdb.SetBit(key, userId, 1).Result()
	if err!=nil{
		logger.Logger.Error("在缓存中创建用户时失败")
		return err
	}
	return nil
}

//使钱包缓存失效
func (rc *redisClient)MakeWalletCacheInvalid(userId int64)error{
	buffer := bytes.Buffer{}
	buffer.WriteString(rc.keyPre)
	buffer.WriteString("userWalletCacahe:")
	buffer.WriteString(strconv.FormatInt(userId,10))
	key:=buffer.String()
	_, err := rc.rdb.Del(key).Result()
	if err!=nil{
		logger.Logger.Error("删除钱包缓存出错")
		return err
	}
	return  nil
}
//将用户钱包列表加入缓存,过期时间的单位为秒
func (rc *redisClient)AddUserWalletToRedis(userId int64,walletJson string,expirationTime int)error  {
	buffer := bytes.Buffer{}
	buffer.WriteString(rc.keyPre)
	buffer.WriteString("userWalletCacahe:")
	buffer.WriteString(strconv.FormatInt(userId,10))
	key:=buffer.String()
	_, err := rc.rdb.Set(key, walletJson, time.Duration(expirationTime)*time.Second).Result()
	if err!=nil{
		logger.Logger.Error("将用户钱包列表缓存进redis时出错")
		return err
	}
	return nil

}

//查看redis缓存中是否存在用户钱包列表的缓存
//如果有则返回缓存值
//如果没有则返回空串""
func (rc *redisClient)GetUserWalletInRedis(userId int64)(string,error){
	buffer := bytes.Buffer{}
	buffer.WriteString(rc.keyPre)
	buffer.WriteString("userWalletCacahe:")
	buffer.WriteString(strconv.FormatInt(userId,10))
	key:=buffer.String()
	result, err := rc.rdb.Get(key).Result()
	if err!=nil{
		if err==redis.Nil{
			return "",nil
		}else {
			return "",err
		}
	}
	return result,nil

}


//通过bitmap查看红包是否拆开
//返回boolean
func (rc *redisClient)RedPacketOpened(redPacketId int64)(bool,error)  {
	key := splicingString(rc.keyPre, redPacketOpenedBitMapKeyName)
	result, err := rc.rdb.GetBit(key, redPacketId).Result()
	if err!=nil{
		logger.Logger.Error("查询红包是否拆开时出错")
		return false, err
	}
	if result==1{
		return true,nil
	}else {
		return false,nil
	}

}

//在redis中将红包标记为打开
func (rc *redisClient)OpenRedPacketInRedisBitMap(redPacketId int64) error {
	key := splicingString(rc.keyPre, redPacketOpenedBitMapKeyName)
	_, err := rc.rdb.SetBit(key, redPacketId,1).Result()
	if err!=nil{
		logger.Logger.Error("将红包标记为拆开时出错")
		return err
	}
	return nil
}



//--------------------------------------下方为工具方法，请勿调用----------------------------------------------------------
//将[]int 转为 []interface{}
func IntSliceConvertToInterfaceSlice(listInt []int)(res []interface{})  {
	res=[]interface{}{}
	for _,v:= range listInt {
		res=append(res,v)
	}
	return res
}

//效率较高的字符串拼接方式
func splicingString(s1 ,s2 string)string{
	buffer := bytes.Buffer{}
	buffer.WriteString(s1)
	buffer.WriteString(s2)
	return buffer.String()
	
}

//从redis中获取当前可获取的红包id
//仅在InitCurrentRedPacketID中调用
//如果获取失败返回-1
func (rc *redisClient)getCurrentRedPacketID()int64{

	result, err := rc.rdb.Get(rc.currentRedPacketIdKeyName).Result()
	if err!=nil{
		if err==redis.Nil{
			//当前redis里面没有当前可获取红包id这个项，检查是否初始化完成
			
			logger.Logger.Error("redis不存在可获取红包id这一项")
			return -1
		}
	}
	atoi, err := strconv.ParseInt(result,10,64)

	//错误处理
	if err!=nil{
		
		logger.Logger.Error("字符串到int转换错误")
		return -1
	}

	return atoi

}

//初始化记录红包是否开启的bitmap，此方法仅在InitRedPacket中调用
func (rc *redisClient)initRedPacketOpenedBitMap()bool{
	key := splicingString(rc.keyPre, redPacketOpenedBitMapKeyName)
	_, err := rc.rdb.SetBit(key, 0, 0).Result()
	if err!=nil{
		
		logger.Logger.Error("初始化记录红包是否开启的BitMap失败")
		return false
	}
	return true
}

//初始化记录用户是否存在的bitmap，此方法仅在InitRedPacket中调用
func (rc *redisClient)initUserExistedBitMap()bool{
	key := splicingString(rc.keyPre, userExistedBitMapKeyName)
	_, err := rc.rdb.SetBit(key, 0, 0).Result()
	if err!=nil{
		
		logger.Logger.Error("初始化记录用户是否存在的BitMap失败")
		return false
	}
	return true
}




//获取缓存中存储的用户抢到过的红包的次数
//此方法不会被用到，已废弃
func(rc *redisClient)GetNumOfRedPacketByUserId (userId int64)int64  {
	buffer := bytes.Buffer{}
	buffer.WriteString(rc.keyPre)
	//这个或许可以抽取成配置文件
	buffer.WriteString("countNumOfPacket:")
	buffer.WriteString(strconv.FormatInt(userId,10))
	key:=buffer.String()

	result, err := rc.rdb.Get(key).Result()
	if err!=nil{
		if err==redis.Nil{
			//当前redis里面没有当前可获取红包id这个项，检查是否初始化完成
			
			logger.Logger.Error("redis中没有用户获取过的红包次数这一项")
			//如果从这里返回就得考虑是真的第一次进还是key过期了，后面从数据库里读
			return -1
		}
	}

	parseInt, err := strconv.ParseInt(result, 10, 64)
	if err!=nil{
		
		logger.Logger.Error("获取用户开过的红包次数时转化为数值失败")
	}
	return parseInt

}