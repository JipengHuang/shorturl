package main

import (
	"db"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
	"regexp"
)

//isbnRegexp 定义了正则，用于判断post传入的参数是否合法
//region ，获取aws环境变量中自定义的区域字段
var isbnRegexp = regexp.MustCompile(`.*`)
var isreghttp = regexp.MustCompile(`(?i:^https).*|(?i:^http).*`)
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var region string = os.Getenv("AWSREGION")
var forwarddomain string = "api.t3t3.top"

//var Awsredisurl = os.Getenv("AWS_REDIS_URL")

// Data : 这是用个用于嵌套的jsion类型
type Data struct {
	Shortenurl string `json:"shortenUrl"`
}

//rejsion : 这里定义一个结构体用于返回jsion
type rejsion struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    `json:"data"`
}

// transTo62 函数用于从10进制生成一个62禁止的字符串并返回
func transTo62(id int64) string {
	// 1 -- > 1
	// 10-- > a
	// 61-- > Z
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var shortURL []byte
	for {
		var result byte
		number := id % 62
		result = charset[number]
		var tmp []byte
		tmp = append(tmp, result)
		shortURL = append(tmp, shortURL...)
		id = id / 62
		if id == 0 {
			break
		}
	}
	//fmt.Println(string(shortURL))
	return string(shortURL)
}

// router 函数为判断 get or post请求方法
func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return show(req)
	case "POST":
		//fmt.Println("+++++++++++++++++++++++++++++++router POST")
		return create(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 获取请求中的路径参数shorturl or 从请求中获取查询 `shorturl` 的字符串参数并校验。
	//shorturl := req.QueryStringParameters["shorturl"]
	shorturl := req.PathParameters["shorturl"]
	if !isbnRegexp.MatchString(shorturl) {
		return clientError(http.StatusBadRequest)
	}

	// 根据 路径参数 shorturl 值从数据库中取出 longurl 记录
	surl, err := db.DygetItem(shorturl)
	if err != nil {
		return serverError(err)
	} else if surl == "" {
		fmt.Println("surl is null ++++++++++++++++++++++++", surl)
		return events.APIGatewayProxyResponse{
			StatusCode: 302,
			Headers:    map[string]string{"Location": fmt.Sprintf("https://s3.t3t3.top")},
		}, nil
	}

	// APIGatewayProxyResponse.Body 域是个字符串，所以 我们将 shorturl 记录解析成 JSON。
	js, err := json.Marshal(surl)
	if err != nil {
		return serverError(err)
	}

	// 返回一个响应，带有代表成功的 200 状态码和 JSON 格式的  记录响应体。
	//jslongurl := println(js.longurl)
	fmt.Println("sur++++++++++++++++++++++++", surl)
	fmt.Println("js=++++++++++++++++++++++++", js)
	//通过短域名重定向到长域名
	return events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers:    map[string]string{"Location": fmt.Sprintf("%v", surl)},
	}, nil
}

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" && req.Headers["Content-Type"] != "application/x-www-form-urlencoded" {
	//	fmt.Printf("Headers is %s", req.Headers)
	//	return clientError(http.StatusNotAcceptable)
	//}

	inShorturldb := new(db.Urljsion)
	err := json.Unmarshal([]byte(req.Body), inShorturldb)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	if !isbnRegexp.MatchString(inShorturldb.SHORTURL) {
		return clientError(http.StatusBadRequest)
	}
	if inShorturldb.LONGURL == "" {
		return clientError(http.StatusBadRequest)
	}

	if !isreghttp.MatchString(inShorturldb.LONGURL) {
		fmt.Println("长域名不是以http开头，处理++", inShorturldb.LONGURL)
		inShorturldb.LONGURL = "https://" + inShorturldb.LONGURL
		fmt.Println("长域名不是以http开头，处理后++", inShorturldb.LONGURL)
		//return clientError(http.StatusBadRequest)
	}

	//strings.HasPrefix(inShorturldb.LONGURL, "https://|http://" string) bool

	shorturlkey := db.RedisKeyId()
	shorturldomain := transTo62(shorturlkey)
	fmt.Println("打印生成生的短域名62进制 is +++++++++++\n", shorturldomain)
	//插入shorturldomain 到数据结构
	inShorturldb.SHORTURL = shorturldomain
	err = db.DyputItem(inShorturldb)
	if err != nil {
		return serverError(err)
	}
	//定义返回的headers ，解决跨域问题
	reheaders := map[string]string{
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "OPTIONS,POST,GET",
		"content-type":                 "application/json",
	}

	bingo := rejsion{"success", "success", Data{fmt.Sprintf("https://%s/%s", forwarddomain, shorturldomain)}}
	jsonBytes, err := json.Marshal(bingo)
	if err != nil {
		panic(err)
	}
	fmt.Printf("打印出需要返回的Response Body 内容++++++++++++++++%s", jsonBytes)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("%s", jsonBytes),
		//Headers:    map[string]string{"Location": fmt.Sprintf("/books?isbn=%s", bk.ISBN)},
		//Headers: map[string]string{"Location": fmt.Sprintf("https://aip.t3t3.top/%s", shorturldomain)},
		Headers: reheaders,
	}, nil
}

// 添加一个用来处理错误的帮助函数。它会打印错误日志到 os.Stderr
// 并返回一个 AWS API 网关能够理解的 500 服务器内部错误
// 的响应。
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// 加一个简单的帮助函数，用来发送和客户端错误相关的响应。
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}
