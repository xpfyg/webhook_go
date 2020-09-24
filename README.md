# webhook_go

### 描述
    gitlab 的webhook 服务，可自定义shell脚本实现自动发布
    sbin目录命名规则为 [project_id].sh,project_id指gitlab项目ID

### 使用
```
# go build
# p 为端口 t为X-Gitlab-Token
# ./webhook_go -p=8080 -t=1234
```