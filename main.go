package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func main() {
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	r.POST("/webhook", Webhook)
	r.Run(":8080")
}

type loginReq struct {
	Message string
}

func Webhook(c *gin.Context) {
	headers := c.Request.Header
	for k, _ := range headers {
		println(k + "==::==" + headers.Get(k))
	}
	body, _ := ioutil.ReadAll(c.Request.Body)
	println("---body/--- \r\n " + string(body))

}

func CheckParamAndHeader(h func(a *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.Request.Header.Get("token")
		if header == "" {
			c.JSON(200, gin.H{
				"code":   3,
				"result": "failed",
				"msg":    ". Missing token",
			})
			return
		}
	}
}
