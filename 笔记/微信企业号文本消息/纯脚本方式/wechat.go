package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// 请求地址：https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=ID&corpsecret=SECRET
const (
	// 企业ID
	corpid = ""
	// 用户id
	userid = ""
	// 部门id
	toparty = 
	// agentid
	agentid = 
	// secret
	secret = ""
	// token请求地址
	getTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	// 缓存tonken文件名
	tokenFileName = ".wechat.txt"
	// 表示是否是保密消息，0表示可对外分享，1表示不能分享且内容显示水印，默认为0
	safe = 0
	// 表示是否开启重复消息检查，0表示否，1表示是，默认0
	enable_duplicate_check = 0
	// 表示是否重复消息检查的时间间隔，默认1800s，最大不超过4小时
	duplicate_check_interval = 3600
)

// 传参
var (
	msg string
	printHelp    bool
	printVersion bool
)

// wechat
type wechat struct {
	corpid string
	userid string
	toparty int
	agentid int
	secret string
	getTokenUrl string
	tokenFileName string
	safe int
	enable_duplicate_check int
	duplicate_check_interval int
	tokenResponse tokenResponse
	sendTextBody sendTextBody
}

// tokenResponse token返回结构
type tokenResponse struct {
	Errcode int
	Errmsg string
	Access_token string
	Expires_in int
}

// 发送消息体
type sendTextBody struct {
	text map[string]string
}

func init() {
	flag.StringVar(&msg,"s","null","发送消息传参！")
	flag.BoolVar(&printHelp,"p",false,"Print this help")
	flag.BoolVar(&printVersion,"v",false,"Print version")
}

// NewWechat 结构体构造函数
func NewWechat(corpid,userid,secret,getTokenUrl,tokenFileName string,toparty,agentid,safe,enable_duplicate_check,duplicate_check_interval int) *wechat {
	return &wechat{
		corpid: corpid,
		userid: userid,
		toparty: toparty,
		agentid: agentid,
		secret: secret,
		safe: safe,
		enable_duplicate_check: enable_duplicate_check,
		duplicate_check_interval: duplicate_check_interval,
		tokenFileName: tokenFileName,
		getTokenUrl: getTokenUrl,
	}
}

// getToken 获取token
//	{"errcode":0,"errmsg":"ok","access_token":"H2s58xxfg9-DQHmkJw","expires_in":7200}
//	errcode	        出错返回码，为0表示成功，非0表示调用失败
//	errmsg	        返回码提示语
//	access_token	获取到的凭证，最长为512字节
//	expires_in	    凭证的有效时间（秒）
func (w *wechat) getToken() {
	// log.Println("--> getToken")
	// 构建token请求url
	requestGetUrl := fmt.Sprintf("%s?corpid=%s&corpsecret=%s",w.getTokenUrl,corpid,secret)
	// log.Println(requestGetUrl)

	// 请求tokon，并转换为结构体
	res,err := http.Get(requestGetUrl)
	defer res.Body.Close()
	checkErr(err)
	body , err := ioutil.ReadAll(res.Body)
	checkErr(err)

	//var data *tokenResponse
	err = json.Unmarshal(body,&w.tokenResponse)
	checkErr(err)

	// 判断token是否返回成功
	if w.tokenResponse.Errmsg != "ok"{
		// log.Fatal(w.tokenResponse.Errmsg)
        log.Println("getToken --> token返回失败！")
		return
	}
	// 缓存到本地文件一份token
    w.cacheToken()

}

// 将token缓存到本地文件
func (w *wechat)cacheToken(){
	f,err := os.OpenFile(w.tokenFileName,os.O_WRONLY|os.O_CREATE,0600)
	// log.Println("cacheToken -->",w.tokenFileName)
	defer f.Close()
	checkErr(err)
	_,err = f.WriteString(w.tokenResponse.Access_token)
	checkErr(err)
	//fmt.Println(n)
}

// 读取本地tokin缓存
func (w *wechat)readToken(){
	content ,err := ioutil.ReadFile(w.tokenFileName)
	checkErr(err)
	w.tokenResponse.Access_token = string(content)
	// log.Println("readToken --> ",w.tokenResponse)
}

// checkErr ...
func checkErr(err error) {
	if err != nil {
		log.Fatal("checkErr --> ",err)
		return
	}
}

// 构造消息体
/*
{
  "toparty": "5",
  "msgtype": "text",
  "agentid": 1000001,
  "text": {
    "content": "你的快递已到，请携带工卡前往邮件中心领取。\n出发前可查看<a href=\"http://work.weixin.qq.com\">邮件中心视频实况</a>，聪明避开排队。"
  },
  "safe": 0,
  "enable_duplicate_check": 0,
  "duplicate_check_interval": 1800
}
*/
func (w *wechat)buildTextBody(messages string) []byte {
	// 构建消息体json
	w.sendTextBody.text = map[string]string{
		"content": messages,
	}
	body := make(map[string]interface{})
	body = map[string]interface{}{
		"toparty": w.toparty,
		"msgtype": "text",
		"agentid": w.agentid,
		"text": w.sendTextBody.text,
		"safe": w.safe,
		"enable_duplicate_check": w.enable_duplicate_check,
		"duplicate_check_interval": w.duplicate_check_interval,
	}
	//log.Println("sendMessage --> ",body)
	data ,err := json.Marshal(body)
	checkErr(err)
	// log.Println("sendMessage --> ",string(data))
	return data
}

// sendMessage 发送文本消息
func (w *wechat)sendTextMessage(data []byte) {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s",w.tokenResponse.Access_token)
	// log.Println("sendMessage --> ",url)
	contentType := "application/json"
	resp ,err := http.Post(url,contentType,strings.NewReader(string(data)))
	defer resp.Body.Close()
	checkErr(err)
	b,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("sendTextMessage --> 发送数据失败！",err)
		return
	}
	var r struct {
		Errcode int
		Errmsg string
		Invaliduser string
		Msgid string
	}
	err = json.Unmarshal(b,&r)
	// log.Println("sendMessage error --> ", err)
	// log.Println("sendMessage response --> ", string(b))
	// log.Println("sendMessage response --> ", r)

	switch r.Errcode {
	case 40014:
		// log.Println("sendMessage --> 42001|41001|40014",r)
		w.getToken()
		w.sendTextMessage(data)
	case 41001:
		// log.Println("sendMessage --> 42001|41001|40014",r)
		w.getToken()
		w.sendTextMessage(data)
	case 42001:
		// log.Println("sendMessage --> 42001|41001|40014",r)
		w.getToken()
		w.sendTextMessage(data)
	case 0:
		log.Println("sendTextMessage --> 发送成功！")
	default:
		log.Fatal("sendTextMessage --> 请求失败！",r)
		return
	}
}

// runTextMessage 发送text消息
func (w *wechat)RunTextMessage(messages string)  {
	// 构建消息体
	data := w.buildTextBody(messages)
	// 发送消息
	w.sendTextMessage(data)
}

func main() {
	flag.Parse()
	if printHelp {
		flag.Usage()
		os.Exit(0)
	} else if printVersion {
		fmt.Printf("%s\n", "version v1")
		os.Exit(0)
	} else if msg == "null" {
		fmt.Println("未传参...")
		flag.Usage()
		os.Exit(0)
	}

	secretData := NewWechat(corpid,userid,secret,getTokenUrl,tokenFileName,toparty,agentid,safe,enable_duplicate_check,duplicate_check_interval)
	//fmt.Println(secretData)

	// 获取token
	//secretData.getToken()

	// 读取本地缓存token
	// secretData.readToken()
	// log.Println(secretData.tokenResponse.Access_token)

	//  发送消息
	messages := msg

	secretData.RunTextMessage(messages)
}
