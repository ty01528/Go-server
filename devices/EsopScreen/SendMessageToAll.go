package EsopScreen

import (
	"first/ReadConfig"
	"first/SQL"
	Client "first/tcp_Client"
	"log"
	"sync"
)

//定义从数据库取数据的类型
type screen struct {
	Ip    string
	Image string
}

func SendMessageToAll() bool {
	//声明esop的IP和data
	var rowsData []screen
	conf := ReadConfig.ReadConfig()
	//声明数据库连接字符串
	conn := SQL.ConnSQL()
	defer conn.Close()
	//编写查询语句
	stmt, err := conn.Prepare(`select 设备网络IP,显示图片 from dbo.esop表单`)
	if err != nil {
		log.Println("Sql Prepare failed:", err.Error())
		return false
	}
	defer stmt.Close()

	//执行查询语句
	rows, err := stmt.Query()
	if err != nil {
		log.Println("Query failed:", err.Error())
		return false
	}
	//将数据读取到实体中
	for rows.Next() {
		var row screen
		rows.Scan(&row.Ip, &row.Image)
		rowsData = append(rowsData, row)
	}
	//读取不到信息则返回空
	if rowsData == nil {
		log.Println("Can not get Data,Please Check DataSources!!!")
		return false
	}
	//读取到信息则通过tcp传递信息
	esopPort := (*conf)["esop_port"]
	var wg sync.WaitGroup
	for _, row := range rowsData {
		ip := row.Ip
		image := row.Image
		wg.Add(1)
		go sendToEsop(ip+esopPort, "Pic:"+"ftp://ftp@192.168.2.46/home/ftp/mes/"+image, &wg)
	}
	wg.Wait()
	return true
}
func sendToEsop(Ip string, Message string, wg *sync.WaitGroup) {
	log.Println(Ip + " " + Message)
	Client.SendMessage(Ip, Message)
	defer wg.Done()
}
