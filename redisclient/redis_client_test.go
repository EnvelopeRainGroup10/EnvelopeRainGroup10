package redisClient

import (
	"bytes"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"reflect"
	"strconv"
	"testing"
)




func Test_redisClient_InitRedPacket(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}


	defer s.Close()
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	type args struct {
		redPacketList []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				redPacketList: []interface{}{
					1,2,3,4,5,6,7,8,9,10,1,2,3,4,5,6,7,8,9,10,
				},
			},
			want: true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}
			got, err := rc.InitRedPacket(tt.args.redPacketList)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.InitRedPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			list, _ := s.List(rc.redPacketListKeyName)

			if !reflect.DeepEqual(list, []string{
				"1","2","3","4","5","6","7","8","9","10","1","2","3","4","5","6","7","8","9","10",
			}){
				t.Errorf("初始化未能成功插入")
			}


			if got != tt.want {
				t.Errorf("redisClient.InitRedPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_InitCurrentRedPacketID(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}


	defer s.Close()
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}
			s.Set(rc.currentRedPacketIdKeyName,"100")
			rc.InitCurrentRedPacketID()
			if rc.currentRedPacketId!=100{
				t.Errorf("未能成功从redis中获取红包id")
			}

		})
	}
}

func Test_redisClient_GetRedPacketMoney(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	type args struct {
		redPacketId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				redPacketId: 1,
			},
			want: 9,
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				redPacketId: 2,
			},
			want: 10,
			wantErr: false,
		},
		{
			name: "test3",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				redPacketId: 3,
			},
			want: 11,
			wantErr: false,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}
			s.RPush(rc.redPacketListKeyName,[]string{
				"9","10","11",
			}...)
			got, err := rc.GetRedPacketMoney(tt.args.redPacketId)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetRedPacketMoney() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("redisClient.GetRedPacketMoney() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_GetCountWithNextRedPacketByUserId(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	buffer := bytes.Buffer{}
	buffer.WriteString("testpre:")
	//这个或许可以抽取成配置文件
	buffer.WriteString("countNumOfPacket:")
	buffer.WriteString(strconv.FormatInt(123,10))
	key:=buffer.String()
	s.Set(key,"4")
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				userId: 123,
			},
			want: 5,
			wantErr: false,

		},
		{
			name: "test2",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				userId: 123,
			},
			want: -1,
			wantErr: false,

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}

			got, err := rc.GetCountWithNextRedPacketByUserId(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetCountWithNextRedPacketByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("redisClient.GetCountWithNextRedPacketByUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_ReduceUserGetRedPacketCount(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	buffer := bytes.Buffer{}
	buffer.WriteString("testpre:")
	//这个或许可以抽取成配置文件
	buffer.WriteString("countNumOfPacket:")
	buffer.WriteString(strconv.FormatInt(123,10))
	key:=buffer.String()
	s.Set(key,"4")
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				userId: 123,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}
			if err := rc.ReduceUserGetRedPacketCount(tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("redisClient.ReduceUserGetRedPacketCount() error = %v, wantErr %v", err, tt.wantErr)
			}
			get,_:= s.Get(key)
			v,_:=strconv.Atoi(get)
			if v!=3{
				t.Errorf("redisClient.ReduceUserGetRedPacketCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func Test_redisClient_GetRedPacket(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	s.Set("testpre:currentPakcetRedPacketIdKeyByName","99")
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 99,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			want: 100,
			wantErr: false,

		},
		{
			name: "test2",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 99,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			want: -1,
			wantErr: false,

		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}
			got, err := rc.GetRedPacket()
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetRedPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("redisClient.GetRedPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_AddToUserRedPacketList(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	type args struct {
		userId      int64
		redPacketId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				userId: 123,
				redPacketId: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}

			buffer := bytes.Buffer{}
			buffer.WriteString(rc.keyPre)
			buffer.WriteString("redPacketList:")
			buffer.WriteString(strconv.FormatInt(tt.args.userId,10))
			key:=buffer.String()

			if err := rc.AddToUserRedPacketList(tt.args.userId, tt.args.redPacketId); (err != nil) != tt.wantErr {
				t.Errorf("redisClient.AddToUserRedPacketList() error = %v, wantErr %v", err, tt.wantErr)
			}
			list, _ := s.List(key)
			if !reflect.DeepEqual(list,[]string{"1"}){
				t.Errorf("redisClient.AddToUserRedPacketList() error = %v, wantErr %v", err, tt.wantErr)
			}


		})
	}
}

func Test_redisClient_AddToUserRedPacketTimeList(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	type args struct {
		userId    int64
		timeStamp int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				userId: 123,
				timeStamp: 12345678,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}
			buffer := bytes.Buffer{}
			buffer.WriteString(rc.keyPre)
			buffer.WriteString("redPacketTimeList:")
			buffer.WriteString(strconv.FormatInt(tt.args.userId,10))
			key:=buffer.String()
			if err := rc.AddToUserRedPacketTimeList(tt.args.userId, tt.args.timeStamp); (err != nil) != tt.wantErr {
				t.Errorf("redisClient.AddToUserRedPacketTimeList() error = %v, wantErr %v", err, tt.wantErr)
			}
			list,_:=s.List(key)
			if !reflect.DeepEqual(list,[]string{"12345678"}){
				t.Errorf("redisClient.AddToUserRedPacketTimeList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisClient_GetUserRedPackerList(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	buffer := bytes.Buffer{}
	buffer.WriteString("testpre:")
	buffer.WriteString("redPacketList:")
	buffer.WriteString(strconv.FormatInt(123,10))
	key:=buffer.String()
	s.RPush(key,"1","2","3")
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				userId: 123,
			},
			want: []string{"1","2","3"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}
			got, err := rc.GetUserRedPackerList(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetUserRedPackerList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisClient.GetUserRedPackerList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_GetUserRedPackerTimeList(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	buffer := bytes.Buffer{}
	buffer.WriteString("testpre:")
	buffer.WriteString("redPacketTimeList:")
	buffer.WriteString(strconv.FormatInt(123,10))
	key:=buffer.String()
	s.RPush(key,"11111111","22222222","33333333")
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				userId: 123,
			},
			want: []string{"11111111","22222222","33333333"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}
			got, err := rc.GetUserRedPackerTimeList(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetUserRedPackerTimeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("redisClient.GetUserRedPackerTimeList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_MakeWalletCacheInvalid(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				userId: 123,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}
			buffer := bytes.Buffer{}
			buffer.WriteString(rc.keyPre)
			buffer.WriteString("userWalletCacahe:")
			buffer.WriteString(strconv.FormatInt(123,10))
			key:=buffer.String()
			s.Set(key,"test")
			if err := rc.MakeWalletCacheInvalid(tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("redisClient.MakeWalletCacheInvalid() error = %v, wantErr %v", err, tt.wantErr)
			}
			_, err2 := s.Get(key)
			if err2==nil{
				t.Errorf("redisClient.MakeWalletCacheInvalid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisClient_AddUserWalletToRedis(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	type args struct {
		userId         int64
		walletJson     string
		expirationTime int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				userId: 123,
				walletJson: "{message:this is a test json}",
				expirationTime: 10,
			},
			wantErr: false,

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}

			if err := rc.AddUserWalletToRedis(tt.args.userId, tt.args.walletJson, tt.args.expirationTime); (err != nil) != tt.wantErr {
				t.Errorf("redisClient.AddUserWalletToRedis() error = %v, wantErr %v", err, tt.wantErr)
			}
			buffer := bytes.Buffer{}
			buffer.WriteString(rc.keyPre)
			buffer.WriteString("userWalletCacahe:")
			buffer.WriteString(strconv.FormatInt(123,10))
			key:=buffer.String()
			v,_:=s.Get(key)
			if v!=tt.args.walletJson{
				t.Errorf("redisClient.AddUserWalletToRedis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisClient_GetUserWalletInRedis(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	buffer := bytes.Buffer{}
	buffer.WriteString("testpre:")
	buffer.WriteString("userWalletCacahe:")
	buffer.WriteString(strconv.FormatInt(123,10))
	key:=buffer.String()
	s.Set(key,"{message:this is a test json}")
	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			args: args{
				userId: 123,
			},
			want: "{message:this is a test json}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}
			got, err := rc.GetUserWalletInRedis(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("redisClient.GetUserWalletInRedis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("redisClient.GetUserWalletInRedis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSliceConvertToInterfaceSlice(t *testing.T) {
	type args struct {
		listInt []int
	}
	tests := []struct {
		name    string
		args    args
		wantRes []interface{}
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				listInt: []int{
					1,2,3,4,5,6,
				},
			},
			wantRes: []interface{}{
				1,2,3,4,5,6,
			},
		},{
			name: "test2",
			args: args{
				listInt: []int{

				},
			},
			wantRes: []interface{}{

			},
		},{
			name: "test3",
			args: args{
				listInt: []int{
					1,2,3,
				},
			},
			wantRes: []interface{}{
				1,2,3,
			},
		},{
			name: "test4",
			args: args{
				listInt: []int{
					1,2,3,4,5,6,2333333333,
				},
			},
			wantRes: []interface{}{
				1,2,3,4,5,6,2333333333,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := IntSliceConvertToInterfaceSlice(tt.args.listInt); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("IntSliceConvertToInterfaceSlice() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_splicingString(t *testing.T) {
	type args struct {
		s1 string
		s2 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				s1: "abc",
				s2: "def",
			},
			want: "abcdef",
		},{
			name: "test2",
			args: args{
				s1: "aaa",
				s2: "ddd",
			},
			want: "aaaddd",
		},{
			name: "test3",
			args: args{
				s1: "",
				s2: "def",
			},
			want: "def",
		},{
			name: "test4",
			args: args{
				s1: "abc",
				s2: "",
			},
			want: "abc",
		},{
			name: "test5",
			args: args{
				s1: "a",
				s2: "d",
			},
			want: "ad",
		},{
			name: "test6",
			args: args{
				s1: "ab",
				s2: "de",
			},
			want: "abde",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splicingString(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("splicingString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisClient_getCurrentRedPacketID(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	type fields struct {
		rdb                       *redis.Client
		keyPre                    string
		currentRedPacketId        int64
		maxPacketNum              int64
		maxGetNum                 int64
		redPacketListKeyName      string
		currentRedPacketIdKeyName string
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				rdb: redis.NewClient(&redis.Options{
					Addr:     s.Addr(),
				}),
				keyPre: "testpre:",
				currentRedPacketId: 0,
				maxPacketNum: 100,
				maxGetNum: 5,
				redPacketListKeyName: "testpre:redPacketListKeyName",
				currentRedPacketIdKeyName: "testpre:currentPakcetRedPacketIdKeyByName",
			},
			want: 66,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &redisClient{
				rdb:                       tt.fields.rdb,
				keyPre:                    tt.fields.keyPre,
				currentRedPacketId:        tt.fields.currentRedPacketId,
				maxPacketNum:              tt.fields.maxPacketNum,
				maxGetNum:                 tt.fields.maxGetNum,
				redPacketListKeyName:      tt.fields.redPacketListKeyName,
				currentRedPacketIdKeyName: tt.fields.currentRedPacketIdKeyName,
			}
			s.Set(rc.currentRedPacketIdKeyName,"66")
			if got := rc.getCurrentRedPacketID(); got != tt.want {
				t.Errorf("redisClient.getCurrentRedPacketID() = %v, want %v", got, tt.want)
			}
		})
	}
}