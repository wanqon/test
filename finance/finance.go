package finance

import (
	"time"
)

const (
	//SourceBase = "/data0/paydata/src_data/db/"
	//TargetBase = "/data0/paydata/finance/"
	SourceBase = "/Users/wangqiong1/app/data1/paydata/src_data/db/"
	TargetBase = "/Users/wangqiong1/app/data1/paydata/finance/"
	TIME_LAYIN  = "2006-01-02"
	TIME_LAYOUT = "2006-01-02 15:04:05"
)

//var (
//	Uri		string
//	Date	string
//)
//
//func init()  {
//	flag.StringVar(&Uri, "uri", "", "请求地址")
//	flag.StringVar(&Date, "date", "", "日期")
//}

//func TaskInit()  {
//	flag.StringVar(&Uri, "uri", "", "请求地址")
//	flag.StringVar(&Date, "date", "", "日期")
//}


func NewCharge(date string) *ChargeBill {
	return &ChargeBill{
		date:      date,
		StartTime: time.Time{},
		EndTime:   time.Time{},
		Fw:        nil,
	}
}