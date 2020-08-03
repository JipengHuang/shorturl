package db

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
)

// 声明一个新的 DynamoDB 实例。注意它在并发调用时是安全的。

var ctx = context.Background()
var region = os.Getenv("AWSREGION")
var Awsredisurl = os.Getenv("AWS_REDIS_URL")
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion(region))

//Urljsion : 定义一个Urljsion的数据类型用于存储到dynamoDB
type Urljsion struct {
	SHORTURL string `json:"shorturl"`
	LONGURL  string `json:"longurl"`
}

func DygetItem(shorturl string) (string, error) {
	// 准备查询的输入
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Shorturldb"),
		Key: map[string]*dynamodb.AttributeValue{
			"shorturl": {
				S: aws.String(shorturl),
			},
		},
	}

	// 从 DynamoDB 检索数据。如果没有符合的数据 返回 nil.
	result, err := db.GetItem(input)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	// 返回的 result.Item 对象具有隐含的
	// map[string]*AttributeValue 类型。我们可以使用 UnmarshalMap helper
	// 解析成对应的数据结构。注意：
	// 当你需要处理多条数据时，可以使用
	// UnmarshalListOfMaps。
	Shorturl := new(Urljsion)
	err = dynamodbattribute.UnmarshalMap(result.Item, Shorturl)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	//
	//if Shorturl.SHORTURL == "" {
	//	fmt.Println("Could not find shorturl", shorturl)
	//	return Shorturl.SHORTURL, nil
	//}
	fmt.Println("返回前打印======shorturl:  ", Shorturl.SHORTURL)
	fmt.Println("返回前打印======longurl: ", Shorturl.LONGURL)

	return Shorturl.LONGURL, nil
}

// DynamoDB 插入短域名62进制ID和长域名
func DyputItem(inShorturldb *Urljsion) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Shorturldb"),
		Item: map[string]*dynamodb.AttributeValue{
			"shorturl": {
				S: aws.String(inShorturldb.SHORTURL),
			},
			"longurl": {
				S: aws.String(inShorturldb.LONGURL),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}
