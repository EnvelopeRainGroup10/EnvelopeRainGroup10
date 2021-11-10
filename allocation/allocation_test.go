package allocation

import (
	"testing"
)



func TestAllocation_AllocateMoney(t *testing.T) {
	type fields struct {
		total_money int
		num         int
		lower       int
		luck        []int
		luck_num    int
	}
	type args struct {
		need_number int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				total_money: 9700,
				num: 990,
				lower: 5,
				luck: []int{30,30,30,30,30,30,30,30,30,30},
				luck_num: 10,
			},
			args: args{
				need_number: 1000,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Allocation{
				total_money: tt.fields.total_money,
				num:         tt.fields.num,
				lower:       tt.fields.lower,
				luck:        tt.fields.luck,
				luck_num:    tt.fields.luck_num,
			}
			 got := a.AllocateMoney(tt.args.need_number)
			 if len(got)!=1000{
				 t.Errorf("红包个数错误，分配失败")
			 }
			 for _,v:=range got{
				 if v<tt.fields.lower{
					 t.Errorf("红包内部金额小于最小值，分配失败")
				 }
			 }
		})
	}
}
