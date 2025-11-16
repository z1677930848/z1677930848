# LingCDN 双端口架构设计

## 架构说明

在同一个 `ling-admin` 程序中运行两个端口：

- **80 端口** → 用户端界面（普通 CDN 用户）
- **7788 端口** → 管理端界面（系统管理员）

## 实现方案

### 1. 配置文件修改

修改 `configs/server.yaml`：

```yaml
env: prod

http:
  "on": true
  listen:
    - "0.0.0.0:80"      # 用户端
    - "0.0.0.0:7788"    # 管理端
```

### 2. 目录结构

```
web/
├── views/
│   ├── @default/          # 管理端界面（7788端口）
│   │   ├── index/         # 管理员登录
│   │   ├── dashboard/     # 管理后台
│   │   └── ...
│   └── @user/             # 用户端界面（80端口）
│       ├── index/         # 用户登录
│       ├── dashboard/     # 用户仪表盘
│       ├── domains/       # 域名管理
│       ├── ssl/           # 证书管理
│       ├── stats/         # 统计数据
│       └── profile/       # 个人设置
```

### 3. 路由分发逻辑

在中间件中根据端口判断：

```go
func PortMiddleware(ctx *actions.ActionObject) {
    port := ctx.Request.URL.Port()

    if port == "80" || port == "8080" {
        // 用户端访问
        ctx.Data["portal"] = "user"
        ctx.ViewDir = "views/@user"
    } else {
        // 管理端访问（7788）
        ctx.Data["portal"] = "admin"
        ctx.ViewDir = "views/@default"
    }

    ctx.Next()
}
```

### 4. 权限隔离

- 用户端（80端口）：只能访问自己的域名、统计、证书
- 管理端（7788端口）：可以管理所有用户和系统设置

### 5. 数据库表设计

```sql
-- 用户表（users 表已存在）
-- 用户可以登录 80 端口

-- 管理员表（admins 表已存在）
-- 管理员可以登录 7788 端口
```

## 优势

1. **共享代码**：用户端和管理端共享同一套后端逻辑
2. **统一维护**：只需要维护一个程序
3. **数据一致**：共享同一个数据库
4. **部署简单**：只需要部署一个服务

## 访问示例

- 用户登录：http://your-domain:80
- 管理登录：http://your-domain:7788

---

**创建时间**: 2025-10-31
