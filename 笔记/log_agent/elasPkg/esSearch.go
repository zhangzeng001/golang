package elasPkg

import (
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/vigneshuvi/GoDateFormat"
	"log"
	"regexp"
	"strings"
	"time"
)

type EsObj struct {
	Server string
	IndexName string
	client *elastic.Client
}

// 初始化客户端,相当于一个构建函数
func NewesObj(server,indexname string)*EsObj{
	// 初始化es客户端
	client, err := elastic.NewClient(elastic.SetURL(server),elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
	//ping(server,client)
	log.Println("connect to es success")
	indexname = DateFormat(indexname)
	return &EsObj{
		Server: server,
		IndexName: indexname,
		client: client,
	}
}

// 检查es可否连通
func ping(server string,client *elastic.Client){
	_, _, err := client.Ping(server).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	return
}

// DateFormat 处理索引时间后缀
func DateFormat(dateformat string)(string){
	// 通过特定格式切割配置文件传参
	splitStr := strings.Split(dateformat,"%")
	// 用户传入dateLayout
	lastStr := splitStr[len(splitStr)-1]
	// 具体key信息
	indexStr := splitStr[0]
	// 匹配dateLayout
	regFormat := regexp.MustCompile("[^{\\+]\\w.*[^}]")
	dateFormat := regFormat.FindAllString(lastStr,-1)
	//fmt.Println(dateFormat[0])  // YYYY.MM.dd

	// 转换为go语言time Layout
	dateLayout := GoDateFormat.ConvertFormat(dateFormat[0])
	//fmt.Println(dateLayout)  // 2006.04.02

	// 根据 time Layout格式化当前时间
	now := time.Now()
	formatTime := now.Format(dateLayout)
	//fmt.Println("dateformat--------",lastStr)
	//fmt.Println("lastStr-----------",lastStr)
	//fmt.Println("indexStr----------",indexStr)

	return indexStr+formatTime
}

// PutData 插入数据
func (e *EsObj)PutData(data interface{})error{
	return e.Exputdata(data)
}

// PutMapping 创建mapping
func (e *EsObj)PutMapping(data string)error{

	return nil
}

// PostData 修改数据
func (e *EsObj)PostData(data string)error{

	return nil
}

// DelData 删除数据
func (e *EsObj)DelData(data string)error{

	return nil
}