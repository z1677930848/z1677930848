# LingCDN 用户端完整实现指南

## ✅ 已完成的工作

### 1. Go 后端代码

创建了完整的用户端 Go 代码结构：

```
internal/web/actions/user/
├── init.go              # 路由总注册
├── index/index.go       # 用户登录
├── dashboard/index.go   # 仪表盘
├── domains/index.go     # 域名管理
├── stats/index.go       # 统计数据
├── profile/index.go     # 个人设置
└── logout/index.go      # 退出登录
```

### 2. 路由系统

已注册的用户端路由：

| 路由 | 功能 | 文件 |
|------|------|------|
| `/user` | 用户登录 | index/index.go |
| `/user/dashboard` | 仪表盘 | dashboard/index.go |
| `/user/domains` | 域名管理 | domains/index.go |
| `/user/domains/create` | 添加域名 | domains/index.go |
| `/user/domains/delete` | 删除域名 | domains/index.go |
| `/user/stats` | 统计数据 | stats/index.go |
| `/user/profile` | 个人设置 | profile/index.go |
| `/user/logout` | 退出登录 | logout/index.go |

### 3. 配置文件

`configs/server.yaml` - 支持双端口：

```yaml
http:
  "on": true
  listen:
    - "0.0.0.0:80"      # 用户端
    - "0.0.0.0:7788"    # 管理端
```

### 4. 文档

- `DUAL_PORT_GUIDE.md` - 架构设计文档
- `DUAL_PORT_README.md` - 使用指南
- `USER_PORTAL_GO_CODE.md` - Go 代码说明

### 5. 部署脚本

- `/root/deploy-user-portal.sh` - 一键编译部署

## 🚀 快速开始

### 方式一：一键部署（推荐）

```bash
bash /root/deploy-user-portal.sh
```

### 方式二：手动部署

```bash
# 1. 编译
cd /root/Lingadmin-master
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -ldflags "-s -w" -o ling-admin ./cmd/edge-admin

# 2. 部署
cp ling-admin /opt/lingcdn/bin/ling-admin
cp configs/server.yaml /opt/lingcdn/configs/server.yaml
cp -r web/views/@user /opt/lingcdn/web/views/@user

# 3. 重启
/opt/lingcdn/bin/ling-admin stop
/opt/lingcdn/bin/ling-admin start
```

## 📍 访问地址

### 用户端（80端口）

- 登录页面: `http://your-ip:80/user`
- 测试账号: `user` / `123456`

### 管理端（7788端口）

- 登录页面: `http://your-ip:7788`
- 管理员账号: 使用原有账号

## 🔑 测试账号

当前使用硬编码的测试账号（正式环境需改为数据库验证）：

```go
// 文件：internal/web/actions/user/index/index.go
username: "user"
password: "123456"
```

## 📂 文件位置

所有创建的文件：

```
/root/Lingadmin-master/
├── internal/web/
│   ├── actions/user/           # 用户端 Go 代码 ⭐
│   └── import.go               # 已添加用户端路由导入
├── web/views/@user/            # 用户端视图
├── configs/server.yaml         # 双端口配置
├── DUAL_PORT_GUIDE.md
├── DUAL_PORT_README.md
└── USER_PORTAL_GO_CODE.md

/root/
├── deploy-user-portal.sh       # 一键部署脚本 ⭐
└── test-dual-port.sh
```

## 🎯 核心功能说明

### 1. 用户登录 (`/user`)

**代码**: `internal/web/actions/user/index/index.go`

**功能**:
- GET: 显示登录页面
- POST: 处理登录逻辑
- 登录成功后创建会话并跳转到仪表盘

**示例**:
```go
func (this *IndexAction) RunPost(params struct {
    Username string
    Password string
}) {
    // 验证用户
    if params.Username == "user" && params.Password == "123456" {
        // 创建会话
        this.CreateUserSession(1, params.Username)
        this.Data["url"] = "/user/dashboard"
    }
}
```

### 2. 仪表盘 (`/user/dashboard`)

**代码**: `internal/web/actions/user/dashboard/index.go`

**功能**:
- 显示用户的统计概览
- 域名数量、请求量、流量等
- 需要登录认证

### 3. 域名管理 (`/user/domains`)

**代码**: `internal/web/actions/user/domains/index.go`

**功能**:
- 列出用户的所有域名
- 添加新域名
- 删除域名

**路由**:
- GET `/user/domains` - 域名列表
- GET/POST `/user/domains/create` - 添加域名
- POST `/user/domains/delete` - 删除域名

### 4. 统计数据 (`/user/stats`)

**代码**: `internal/web/actions/user/stats/index.go`

**功能**:
- 今日/本周/本月请求统计
- 流量统计

### 5. 个人设置 (`/user/profile`)

**代码**: `internal/web/actions/user/profile/index.go`

**功能**:
- 查看个人信息
- 修改邮箱、手机、密码

### 6. 退出登录 (`/user/logout`)

**代码**: `internal/web/actions/user/logout/index.go`

**功能**:
- 清除会话
- 跳转到登录页

## ⚙️ 认证机制

每个需要登录的页面都使用：

```go
Auth *actionutils.UserMustAuth `action:"user"`
```

这会自动检查：
1. 用户是否已登录
2. 是否是用户端登录（区别于管理员）
3. 未登录则跳转到 `/user`

## 🔄 与管理端的对应

| 功能 | 用户端（80） | 管理端（7788） |
|------|-------------|---------------|
| 登录 | `/user` | `/` |
| 首页 | `/user/dashboard` | `/dashboard` |
| 域名 | `/user/domains` | `/servers` |
| 统计 | `/user/stats` | `/servers/stats` |
| 设置 | `/user/profile` | `/settings/profile` |

## 📋 后续待完成

### 高优先级

- [ ] 实现 RPC 调用获取真实数据
- [ ] 完善用户认证（对接数据库）
- [ ] 创建用户端的 HTML 模板文件
- [ ] 实现域名增删改查的完整功能

### 中优先级

- [ ] 添加 SSL 证书管理
- [ ] 完善统计数据展示
- [ ] 添加访问日志查询
- [ ] 实现密码修改功能

### 低优先级

- [ ] 添加双因素认证
- [ ] 添加邮件通知
- [ ] 优化界面设计

## ⚠️ 注意事项

1. **临时实现**: 当前很多功能使用临时数据，需要对接 RPC
2. **测试账号**: 登录使用硬编码账号，生产环境需要改为数据库验证
3. **权限控制**: 需要确保用户只能访问自己的资源
4. **视图文件**: Go 代码已完成，HTML 模板还需要创建

## 🧪 测试流程

1. **编译部署**:
   ```bash
   bash /root/deploy-user-portal.sh
   ```

2. **测试用户登录**:
   - 访问 `http://your-ip:80/user`
   - 输入 `user` / `123456`
   - 应该跳转到 `/user/dashboard`

3. **测试路由**:
   ```bash
   curl http://localhost:80/user
   curl http://localhost:80/user/dashboard
   curl http://localhost:80/user/domains
   ```

4. **测试管理端**:
   - 访问 `http://your-ip:7788`
   - 确保管理端正常工作

## 📖 相关文档

- [架构设计](DUAL_PORT_GUIDE.md)
- [使用指南](DUAL_PORT_README.md)
- [Go代码说明](USER_PORTAL_GO_CODE.md)

---

**版本**: 1.0.0
**更新时间**: 2025-10-31
**作者**: Claude Code
