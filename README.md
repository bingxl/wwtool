# 简介
wwtool 是一个用于 `鸣潮` 的小工具

实现的功能有
- b服/官服 切换
- 启动游戏
- 获取抽卡链接

详细操作方法见软件的帮助页

# 添加manifest
安装rsrc `go get github.com/akavel/rsrc`

使用 rsrc 操作 `rsrc -manifest app.manifest -ico icon.icon -o resource.syso`

使用fyne 打包时会自动打包根目录的Icon.png, 执行rsrc时不需要指定ico参数

# 打包
1. 安装了fyne 时使用 `fyne package`
2. 使用go直接打包 `go build -ldflags="-s -w -H windowsgui" .`