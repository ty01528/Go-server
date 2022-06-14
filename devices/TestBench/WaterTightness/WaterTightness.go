package TestBench

import (
	"errors"
	"first/SQL"
	"log"
	"time"
)

//#--SN:0123456789011
//--ItemName:Water_tightness_test
//--Result:OK
//--Water_tightness_ST:2022-06-13 15:49:06
//--Water_tightness_ET:2022-06-13 15:49:37
//--ALM:0,0.015,10#
func WaterTightnessTestFunc(message map[string]string) error {
	log.Println("=========================正在写入水路气密性检测报告===============================")
	conn := SQL.ConnSQL()
	defer conn.Close()
	SN, Result, St, Et, ItemNm :=
		message["SN"], message["Result"],
		message["Water_tightness_ST"], message["Water_tightness_ET"],
		message["ALM"]

	//将时间字符串序列化为标准时间格式
	StartTime, err := time.ParseInLocation("2006-01-02 15:04:05", St, time.Local)
	if err != nil {
		log.Printf("综合性能检测台转换时间格式失败！请检查TCP连接与传递信息是否正确")
	}
	EndTime, err := time.ParseInLocation("2006-01-02 15:04:05", Et, time.Local)
	if err != nil {
		log.Printf("综合性能检测台转换时间格式失败！请检查TCP连接与传递信息是否正确")
	}
	//从数据库查询订单所属的编号等信息
	var (
		orderNumber    string
		productNo      string
		productionLine string
		productCode    string
	)
	rows, err := conn.Query(`select 订单编号,产品型号,生产线别,产品代码 from dbo.设备ID生成情况 where 设备序列号 = ?`, SN)
	if err != nil {
		log.Println("Query failed:", err.Error())
		return errors.New("从设备ID生成情况拉取相关信息失败")
	}
	for rows.Next() {
		err := rows.Scan(&orderNumber, &productNo, &productionLine, &productCode)
		if err != nil {
			//log.Fatal(err)
			return errors.New("遍历行错误！请查看错误信息")
		}
	}

	//将接收到的结果插入数据库
	stmt, err := conn.Prepare("INSERT INTO dbo.水路气密性检测 (订单编号,产品代码,工序代码,设备序列ID,开始时间,结束时间,测试流程文件名,测试结果) VALUES (?,?,?,?,?,?,?,?)")
	if len(orderNumber) == 0 || len(productCode) == 0 || stmt == nil {
		return errors.New("从和数据库查询的订单信息为空，请检查数据库是否已录入该订单！")
	}
	_, err = stmt.Exec(orderNumber, productCode, "P0110", SN, StartTime, EndTime, ItemNm, Result)
	if err != nil {
		log.Print(err)
		return errors.New("水路气密性检测执行数据库写入错误！")
	}

	log.Println("=========================水路气密性检测数据写入操作完成=============================")
	return nil
}
