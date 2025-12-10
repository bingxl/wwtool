# WWTool 项目概述

WWTool 是一个多游戏辅助工具集合，使用 Go 语言和 Fyne UI 框架开发的桌面应用程序。主要支持《鸣潮》和《原神》两款游戏，提供服务器切换、游戏启动、抽卡链接获取和库街区签到等功能。

## 项目架构

- **语言**: Go 1.25.3
- **UI 框架**: Fyne v2.7.0
- **主要依赖**: 
  - golang.org/x/text (国际化支持)
  - gopkg.in/yaml.v3 (配置文件解析)
  - gopkg.in/ini.v1 (INI 配置文件解析)

### 目录结构

```
├── cmd/              # 命令行工具
│   ├── binary_search/ # 二分搜索工具
│   └── diff/          # 文件差异工具
├── genshin/          # 原神相关功能
│   └── lib.go         # 原神服务器切换实现
├── i18n/             # 国际化文件
│   └── locales/       # 语言文件目录
├── lib/              # 核心功能库
├── model/            # 数据模型
├── view/             # UI 界面
├── viewmodel/        # 视图模型
└── wuwa/             # 鸣潮相关功能（独立子模块）
    ├── card/         # 抽卡记录解析
    ├── cmd/          # 鸣潮命令行工具
    ├── kujiequ/      # 库街区功能
    ├── req/          # 请求处理
    └── server/       # 服务器切换
```

## 核心功能

### 1. 鸣潮游戏功能
- **服务器切换**: 支持在官服、B服和WeGame之间切换
  - 通过创建软链接方式实现服务器切换
  - 配置文件: `model/default_config.yaml`
- **游戏启动**: 
  - 支持多游戏路径管理
  - 可配置游戏可执行文件名
  - 路径配置存储在用户配置中
- **抽卡链接获取**: 
  - 从游戏日志文件中提取抽卡记录链接
  - 自动复制到剪贴板
- **库街区功能**: 
  - 库街区签到
  - Token 管理
  - 用户信息管理
  - 支持鸣潮和战双帕弥什

### 2. 原神游戏功能
- **服务器切换**: 
  - 支持官服和B服之间切换
  - 通过修改配置文件实现
  - 自动管理B服专用SDK (PCGameSDK.dll)
- **配置文件管理**: 
  - 修改游戏目录下的 config.ini
  - 自动处理 channel、cps、sub_channel 等参数

### 3. 通用工具功能
- **第三方工具集成**: 支持管理其他游戏辅助工具
- **文件差异工具**: `cmd/diff/` 提供文件对比功能
- **二分搜索工具**: `cmd/binary_search/` 提供搜索功能

## 构建和运行

### 开发环境要求
- Go 1.25.3+
- Fyne CLI (可选，用于打包)

### 构建命令

#### 主应用构建
使用 Task (推荐):
```bash
# Windows
task build

# Linux
task build:linux

# Android
task build:android
```

使用 Go 直接构建:
```bash
# Windows
go build -ldflags="-s -w -H windowsgui" -o ./wwtool.exe .

# Linux
go build -ldflags="-s -w" -o ./wwtool .
```

使用 Fyne 打包:
```bash
fyne package -os android -app-id cn.bingxl.wwtool -icon Icon.ico
```

#### 鸣潮子模块构建
鸣潮子模块有独立的构建系统:
```bash
# 进入鸣潮目录
cd wuwa

# 构建签到工具
task build:windows  # Windows版本
task build:linux    # Linux版本

# 部署到服务器
task deploy
```

### 运行应用
```bash
# 运行主应用
go run .

# 运行鸣潮签到工具
cd wuwa && go run ./cmd/signin/main.go

# 获取抽卡链接
cd wuwa && go run ./cmd/getgachalink/getgachalink.go
```

## 开发约定

### 配置管理
- 默认配置嵌入在二进制文件中 (`model/default_config.yaml`)
- 用户配置存储在 Fyne 应用存储目录
- 配置格式为 YAML
- 原神使用 INI 格式配置文件

### 国际化
- 使用 `golang.org/x/text` 实现多语言支持
- 语言文件位于 `i18n/locales/` 目录
- 翻译函数: `i18n.T(key)`
- 鸣潮子模块有独立的国际化系统

### 代码结构
- **Model**: 数据结构和配置管理
- **View**: UI 界面组件
- **ViewModel**: 业务逻辑和数据绑定
- **Lib**: 核心功能实现
- **Genshin**: 原神特定功能
- **Wuwa**: 鸣潮特定功能（独立模块）

### 测试
- 单元测试文件以 `_test.go` 结尾
- 测试命令: `go test ./...`
- 鸣潮子模块有独立的测试: `cd wuwa && go test ./...`

## 项目配置

### 应用清单
- Windows 应用清单: `app.manifest`
- 应用图标: `Icon.ico`

### 资源嵌入
- 使用 Go embed 嵌入默认配置文件
- 资源文件: `rsrc_windows_amd64.syso`

### 环境变量
- 鸣潮子模块支持 `.env` 文件配置
- 支持通过环境变量配置部署服务器

## 子模块

### wuwa 模块
独立的鸣潮相关功能模块，包含:
- 库街区 API 交互
- 抽卡记录解析
- 用户数据管理
- 独立的构建和部署系统
- 支持战双帕弥什游戏

### 命令行工具
- `cmd/binary_search/`: 二分搜索工具
- `cmd/diff/`: 文件差异工具
- `wuwa/cmd/signin/`: 库街区签到工具
- `wuwa/cmd/getgachalink/`: 抽卡链接获取工具

### 原神模块
- `genshin/lib.go`: 原神服务器切换实现
- 支持 B服 SDK 自动管理
- 配置文件自动修改

## 注意事项

1. **平台限制**: 由于使用了 CGO，交叉编译受限
2. **Windows 特定**: 部分功能仅在 Windows 平台可用
3. **路径配置**: 游戏路径需要指向游戏根目录，而非启动器目录
4. **软链接**: 服务器切换功能需要管理员权限创建软链接
5. **子模块独立性**: 鸣潮子模块可以独立开发和部署
6. **多游戏支持**: 现在同时支持鸣潮和原神两款游戏

## 扩展开发

### 添加新游戏支持
1. 在项目根目录创建游戏专用目录 (如 `genshin/`)
2. 在 `model/` 中定义游戏配置结构
3. 在 `lib/` 中实现通用功能
4. 在 `view/` 中添加游戏专用 UI
5. 在 `viewmodel/` 中添加游戏业务逻辑

### 添加新功能
1. 在 `model/` 中定义数据结构
2. 在 `lib/` 中实现核心逻辑
3. 在 `viewmodel/` 中添加业务逻辑
4. 在 `view/` 中创建 UI 组件

### 国际化支持
1. 在 `i18n/locales/` 中添加语言文件
2. 使用 `i18n.T()` 函数标记需要翻译的文本
3. 运行 `go generate` 更新翻译文件
4. 鸣潮子模块需要单独处理国际化