package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/exec"
	"strconv"
)

var (
	h     bool
	p     string
	t     string
	TOKEN string
)

func init() {
	flag.StringVar(&t, "t", "", "X-Gitlab-Token")
	flag.StringVar(&p, "p", "", "端口")
	flag.Usage = usage
}
func usage() {
	fmt.Fprintf(os.Stderr, `webhook_go webhook工具, 支持自定义shell脚本
Usage: webhook_go [-p] [-t]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if p == "" {
		flag.Usage()
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CheckHeader())
	r.POST("/webhook", Webhook)
	TOKEN = t
	err := r.Run(fmt.Sprintf(":%s", p))
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

type WebhookReq struct {
	ProjectId         int      `json:"project_id"`
	Ref               string   `json:"ref"`
	Commits           []Commit `json:"commits"`
	TotalCommitsCount int      `json:"total_commits_count"`
}
type Commit struct {
	Id       string   `json:"id"`
	Message  string   `json:"message"`
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Removed  []string `json:"removed"`
}

func Webhook(c *gin.Context) {
	var req WebhookReq
	c.BindJSON(&req)
	Process(req)
}

func Process(req WebhookReq) {
	shell_path := "sbin/" + strconv.Itoa(req.ProjectId) + ".sh"
	println(shell_path)
	cmd := exec.Command("sh", shell_path)
	// 命令的错误输出和标准输出都连接到同一个管道
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		println("Process 1", err.Error())
		return
	}

	if err = cmd.Start(); err != nil {
		println("Process 2", err.Error())
		return
	}
	// 从管道中实时获取输出并打印到终端
	println("/---------------**************------------\\")
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		print(string(tmp))
		if err != nil {
			break
		}
	}
	println("\\---------------**************------------/")
	if err = cmd.Wait(); err != nil {
		println("Process 2", err.Error())
		return
	}
}

func CheckHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if TOKEN == "" {
			return
		}
		token := c.Request.Header.Get("X-Gitlab-Token")
		if token == "" {
			c.JSON(200, gin.H{
				"code":   3,
				"result": "failed",
				"msg":    ". Missing token",
			})
			c.Abort()
		}
	}
}
