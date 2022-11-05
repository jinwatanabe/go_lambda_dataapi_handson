package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/krotscheck/go-rds-driver"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ginLambda *ginadapter.GinLambda

type User struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func init() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {

		conf := &rds.Config{
			ResourceArn: "RDSのリソースARN",
			SecretArn:   "シークレットARN",
			Database:    "handson",
			AWSRegion:   "ap-northeast-1",
			SplitMulti:  false,
			ParseTime:   true,
		}
		dsn := conf.ToDSN()

		DB, err := gorm.Open(mysql.New(mysql.Config{
			DriverName: rds.DRIVERNAME,
			DSN:        dsn,
		}), &gorm.Config{})

		if err != nil {
			fmt.Println(err)
		}

		err = DB.AutoMigrate(&User{})

		if err != nil {
			fmt.Println(err)
		}

		c.JSON(200, gin.H{
			"message": "Success!",
		})

	})

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
