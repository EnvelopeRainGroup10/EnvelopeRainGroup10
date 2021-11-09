package sql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

type User struct {
	ID    int64
	Count int64
}

func (User) TableName() string {
	//return "users"
	return "user"
}

type Envelope struct {
	ID         int64 `json:"envelope_id"`
	UID        int64 `json:"uid"`
	Opened     bool  `json:"opened"`
	Value      int64 `json:"value"`
	SnatchTime int64 `json:"snatch_time"`
}

func (Envelope) TableName() string {
	//return "envelopes"
	return "red_envelope"
}

var DB *gorm.DB

const DbType string = "mysql"
const DbAddress string = "root:123456@(180.184.71.7:8066)/envelope_db?charset=utf8&parseTime=True&loc=Local"
//const DbAddress string = "root:3306@tcp(mysql:3306)/test?charset=utf8&parseTime=True&loc=Local"
//const DbAddress string = "root:123456@(192.168.0.41:8066)/envelope_db?charset=utf8&parseTime=True&loc=Local"

func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(DbType, DbAddress)
	if err == nil {
		DB = db
		//db.AutoMigrate(&User{}, &Envelope{})
		return db, err
	}
	log.Println(err)
	return nil, err
}

// GetUser 查询用户并返回，如果不存在则创建用户
func GetUser(uid int64) (user User) {
	DB.FirstOrCreate(&user, User{ID: uid})
	return
}

func UpdateCount(user *User) {
	user.Count++
	DB.Model(&user).Update("count", user.Count)
}

// UpdateCountByUid 根据用户id更新count
func UpdateCountByUid(uid int64) {
	DB.Model(&User{}).Where(User{ID: uid}).Update("count", gorm.Expr("count + ?", 1))
}

//Envelope
func GetAllEnvelopesByUID(uid int64) ([]*Envelope, error) {

	var envelopes []*Envelope
	conditions := map[string]interface{}{
		"uid": uid,
	}
	if err := DB.Table(Envelope{}.TableName()).Where(conditions).Find(&envelopes).Error; err != nil {
		return nil, err
	}
	return envelopes, nil
}

func GetEnvelopeByEnvelopeID(envelopeId int64) (envelope Envelope) {
	DB.Where("id = ?", envelopeId).First(&envelope)
	return
}

// GetEnvelopeByEnvelopeIDAndUid 根据用户id与红包id查询红包
func GetEnvelopeByEnvelopeIDAndUid(envelopeId int64, uid int64) (envelope Envelope) {
	DB.Where("id = ? and uid = ? ", envelopeId, uid).First(&envelope)
	return
}

func CreateEnvelope(user User) (envelope Envelope) {

	snatchTime := time.Now().UnixNano()
	var a int64 = 10
	envelope = Envelope{UID: user.ID, Opened: false, Value: a, SnatchTime: snatchTime}
	DB.Create(&envelope)
	return envelope
}

func UpdateState(envelopeId int64) (envelope Envelope) {

	//查询条件
	envelope.ID = envelopeId
	DB.Model(&envelope).Update("opened", true)
	return
}

// CreateEnvelopeDetail 根据redis中的红包id, 用户id, 红包价值value 与创建时间snatch_time创建红包记录
func CreateEnvelopeDetail(envelopeId int64, uid int64, value int64, snatchTime int64) (envelope Envelope) {

	envelope = Envelope{ID: envelopeId, UID: uid, Opened: false, Value: value, SnatchTime: snatchTime}
	if err := DB.Create(&envelope).Error; err != nil {
		return Envelope{ID: 0}
	}
	return envelope
}

// UpdateStateByEidAndUid 根据用户id与红包id更新红包状态，将红包设置为已经打开状态
func UpdateStateByEidAndUid(envelopeId int64, uid int64) {
	envelope:=Envelope{ID: envelopeId,UID: uid}
	DB.Model(&envelope).Where(envelope).Update("opened", true)
}