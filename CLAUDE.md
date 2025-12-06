# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Lingcdn 是一个 CDN 系统，包含四个核心组件，通过 gRPC 进行通信：

| 组件 | 目录 | 进程名 | 作用 |
|------|------|--------|------|
| EdgeCommon | `EdgeCommon/` | - | 公共库，包含 RPC 定义、配置模型、工具函数 |
| Lingnode | `Lingnode-master/` | ling-node | 边缘节点，处理用户请求、缓存、WAF |
| Lingapi | `Lingapi-master/` | ling-api | API 服务，管理节点、域名、证书 |
| Lingadmin | `Lingadmin-master/` | lingcdnadmin | 管理后台，Web UI 和配置管理 |

**依赖关系：** EdgeCommon ← Lingnode/Lingapi/Lingadmin，构建时 Lingadmin 会自动构建 Lingapi，Lingapi 会自动构建 Lingnode。

## 仓库布局要求

各组件通过 `go.mod` 的 `replace` 指令引用本地兄弟目录：
```
父目录/
├── EdgeCommon/      # 必须存在
├── Lingnode-master/ # 对应 EdgeNode
├── Lingapi-master/  # 对应 EdgeAPI
└── Lingadmin-master/
```

构建脚本默认查找 `../../EdgeCommon`、`../../EdgeNode`、`../../EdgeAPI`。可通过环境变量覆盖：
- `EDGENODE_PATH` - 指定 EdgeNode 位置
- `EDGEAPI_PATH` - 指定 EdgeAPI 位置

## 构建命令

**前置要求：** Go 1.25+、bash (WSL/Linux/macOS)、zip、unzip

```bash
# 构建边缘节点 (输出: dist/ling-node-linux-amd64-v*.zip)
cd Lingnode-master/build && ./build.sh linux amd64 community

# 构建 API 服务 (会自动构建 Lingnode，输出: dist/edge-api-linux-amd64-community-v*.zip)
cd Lingapi-master/build && ./build.sh linux amd64 community

# 构建管理后台 (会自动构建 Lingapi，输出: dist/lingcdnadmin-linux-amd64-community-v*.zip)
cd Lingadmin-master/build && ./build.sh linux amd64 community
```

**构建标签：** `community` (社区版) 或 `plus` (商业版，启用额外功能如 libpcap)

**交叉编译：** 支持 `amd64` 和 `arm64`，静态链接需要 musl 工具链 (`x86_64-linux-musl-gcc`)

## 开发命令

```bash
# 运行管理后台开发服务器
go run ./cmd/lingcdnadmin/main.go
go run ./cmd/lingcdnadmin/main.go dev  # 开发模式

# 运行代码生成 (生成前端组件)
go run ./cmd/lingcdnadmin/main.go generate

# 前端构建 (压缩 JS)
cd Lingadmin-master/web && npm install && npm run build

# 运行测试
go test ./...
go test -race ./...  # 并发测试
```

## 代码架构

### RPC 通信
- Protocol Buffer 定义位于 `EdgeCommon/pkg/rpc/protos/` (100+ 服务定义)
- 数据模型位于 `EdgeCommon/pkg/rpc/protos/models/`
- 各组件通过 gRPC 通信，边缘节点定期向 API 服务上报状态

### 各组件核心模块

**Lingnode-master:**
- `internal/caches/` - 缓存管理
- `internal/waf/` - Web 应用防火墙
- `internal/firewalls/` - 防火墙规则
- `internal/nodes/` - 节点管理和 RPC 通信

**Lingapi-master:**
- `internal/db/` - 数据库操作
- `internal/acme/` - ACME 证书自动申请
- `internal/tasks/` - 后台任务调度
- `internal/installers/` - 节点安装器

**Lingadmin-master:**
- `internal/web/` - Web 框架 (基于 TeaGo)
- `web/views/` - 视图模板 (`@default/` 管理员, `@user/` 用户门户)
- `web/public/js/components.src.js` - 前端组件源码

### 版本号管理
版本号定义在各组件的 `internal/const/const.go`，构建脚本自动读取。

### 命名约定
项目存在多套历史名称：
- `EdgeAdmin/EdgeAPI/EdgeNode` - 上游项目名
- `ling-admin/ling-api/ling-node` - 进程名 (`ProcessName`)
- `LingCDN` - 产品名 (`ProductName`)

代码中应使用 `internal/const/const.go` 中的常量。

## 关键文件位置

| 用途 | 路径 |
|------|------|
| 程序入口 | `cmd/lingcdnadmin/main.go`, `cmd/ling-node/main.go`, `cmd/edge-api/main.go` |
| 版本常量 | `internal/const/const.go` |
| 运行时配置 | `configs/server.yaml` (从 `server.template.yaml` 生成) |
| 构建脚本 | `build/build.sh` |
| 前端组件 | `web/public/js/components.src.js` → 压缩为 `components.js` |

## 语言规则

1. **只允许使用中文回答** - 所有思考、分析、解释和回答都必须使用中文
2. **中文注释** - 生成的代码注释和文档都应使用中文

## IPC 机制

管理后台使用 `gosock` 进行进程间通信，CLI 命令 (`dev`/`prod`/`stop`) 会向运行中的进程发送指令。测试时注意是否已有进程监听 socket。
