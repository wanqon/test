package finance

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"pay.sc.weibo.com/accounts/tool"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	//充值状态_关闭
	CHARGE_STATUS_CLOSED = 0
	//充值状态_新创建
	CHARGE_STATUS_CREATE = 1
	//充值状态_充值成功
	CHARGE_STATUS_CHARGED = 2
)

type ChargeBill struct {
	date	string
	StartTime	time.Time
	EndTime		time.Time
	Fw			*os.File
}

type HandleInfo func([]string)

func WbpayCharge(date string) {
	chargeBill := &ChargeBill{
		date:      date,
		StartTime: time.Time{},
		EndTime:   time.Time{},
		Fw:        nil,
	}
	chargeBill.Run()
}


func (bill *ChargeBill) Run() {
	l := time.FixedZone("CST", 8*3600)
	var t time.Time
	if len(bill.date) == 0 {
		t = time.Now().AddDate(0,0,-1).In(l)
	} else {
		t, _ = time.Parse(TIME_LAYIN,bill.date)
	}
	dbDate := t.Format("20060102")
	basePath := tool.GetConfString("dir", "data")
	sourceDir := fmt.Sprintf(basePath+"/src_data/db/charge/%s/",dbDate)
	year,month,day := t.Date()
	targetDir := fmt.Sprintf(basePath+"/finance/charge/%d%02d/%02d/", year, int(month), day)
	if _, err := os.Stat(targetDir);os.IsNotExist(err) {
		err := os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	bill.StartTime = time.Date(t.Year(),t.Month(),t.Day(),0,0,0,0,l)
	bill.EndTime = time.Date(t.Year(),t.Month(),t.Day(),23,59,59,0,l)

	var err error
	wFilePath := targetDir+"charge.txt"
	bill.Fw, err = os.OpenFile(wFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND,os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer bill.Fw.Close()

	swg := sync.WaitGroup{}
	for i:=0; i<128; i++ {
		fileName:=fmt.Sprintf(sourceDir+"snap_%d.txt", i)
		go ReadFile(&swg, fileName, bill.handleInfo)
	}
	swg.Wait()
}

func (bill *ChargeBill) handleInfo(charge []string) {
	chargeTime, _ := time.Parse(TIME_LAYOUT, charge[9])
	status,_ := strconv.Atoi(charge[7])
	if TimeCompare(chargeTime,bill.StartTime,bill.EndTime) && status == CHARGE_STATUS_CHARGED {
		info := strings.Join(charge[:11],"\t")
		info += "\r\n";
		bill.Fw.Write([]byte(info))
	}
}

func TimeCompare(check time.Time, start time.Time, end time.Time) bool {
	return check.Second() >= start.Second() && check.Second() <= end.Second()
}

func ReadFile(swg *sync.WaitGroup, path string, handlefuc HandleInfo)  {
	swg.Add(1)
	defer swg.Done()
	fr, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println(err)
	}
	defer fr.Close()
	reader := bufio.NewReader(fr)
	for {
		line,_,err := reader.ReadLine()
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		chargeInfo := strings.Split(string(line),"\t")
		handlefuc(chargeInfo)
	}
}
