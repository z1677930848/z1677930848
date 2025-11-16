# LingCDN 用户端与管理端集成说明

## 概述

LingCDN 双端口架构已完成集成，用户端登录与管理端用户管理功能完全对接。

## 架构说明

### 端口分配
- **8080端口**：用户端界面（自动重定向到 /user）
- **7788端口**：管理端界面

### 访问地址
- 用户端登录：http://10.0.54.2:8080/ 或 http://10.0.54.2:8080/user
- 管理端登录：http://10.0.54.2:7788/

## 用户管理流程

### 1. 管理员创建用户

管理员通过管理端界面创建用户：

1. 登录管理端：http://10.0.54.2:7788/
2. 进入"用户管理"模块
3. 点击"创建用户"
4. 填写用户信息：
   - 用户名（只能包含英文、数字、下划线）
   - 密码
   - 全名
   - 手机号
   - 邮箱
   - 关联集群
   - 备注

### 2. 用户登录

用户使用管理员创建的账号登录：

1. 访问用户端：http://10.0.54.2:8080/
2. 自动跳转到用户登录页面（紫色渐变背景）
3. 输入用户名和密码
4. 登录后进入用户控制面板

## 技术实现

### 端口重定向

`/root/Lingadmin-master/internal/web/actions/default/index/index.go`

```go
// 第32-36行：端口检测
requestHost := this.Request.Host
if strings.Contains(requestHost, ":8080") {
    this.RedirectURL("/user")
    return
}
```

### 用户认证对接

`/root/Lingadmin-master/internal/web/actions/user/index/index.go`

```go
// 通过RPC验证用户账号（第65-85行）
rpcClient, err := rpc.SharedRPC()
if err != nil {
    this.Fail("系统错误，请稍后重试")
    return
}

// 调用UserRPC.LoginUser进行验证
resp, err := rpcClient.UserRPC().LoginUser(rpcClient.Context(0), &pb.LoginUserRequest{
    Username: params.Username,
    Password: params.Password,
})
if err != nil {
    this.FailField("username", "用户名或密码错误")
    return
}

var userId = resp.UserId
if userId <= 0 {
    this.FailField("username", "用户名或密码错误")
    return
}

// 使用真实的用户ID创建会话
params.Auth.StoreAdmin(userId, params.Remember)
```

### 界面优化

用户端登录页面已移除管理员登录链接：

- 位置：`/opt/lingcdn/web/views/@user/index/index.html`
- 纯净的用户登录界面
- 紫色渐变背景设计
- 移动端自适应

## 测试验证

### 测试步骤

1. **端口监听测试**
```bash
netstat -tlnp | grep ling-admin
# 应该显示8080和7788两个端口
```

2. **端口重定向测试**
```bash
curl http://localhost:8080/
# 应该返回重定向到 /user
```

3. **用户登录测试**
   - 访问 http://10.0.54.2:8080/
   - 应该显示用户登录页面
   - 页面应该是紫色渐变背景
   - 不应该有"管理员登录"链接

4. **管理端测试**
   - 访问 http://10.0.54.2:7788/
   - 应该显示管理员登录页面
   - 登录后可以创建用户

## 文件清单

### 修改的文件
1. `/root/Lingadmin-master/internal/web/actions/default/index/index.go`
   - 添加端口检测和重定向逻辑

2. `/root/Lingadmin-master/internal/web/actions/user/index/index.go`
   - 实现RPC用户认证对接

3. `/root/Lingadmin-master/web/views/@user/index/index.html`
   - 移除管理员登录链接
   - 优化界面样式

4. `/root/Lingadmin-master/configs/server.yaml`
   - 配置双端口：8080和7788

### 部署的文件
1. `/opt/lingcdn/bin/ling-admin` - 编译后的二进制文件
2. `/opt/lingcdn/web/views/@user/` - 用户端视图文件

## 注意事项

1. **用户必须由管理员创建**
   - 用户端不提供注册功能
   - 所有用户账号由管理员在管理端创建

2. **用户名规范**
   - 只能包含英文字母、数字和下划线
   - 管理端创建时会自动验证

3. **密码安全**
   - 密码在RPC传输时加密
   - 数据库存储使用哈希

4. **集群关联**
   - 每个用户必须关联到一个集群
   - 用户只能管理自己集群内的资源

## 当前状态

✅ 双端口架构运行正常
✅ 端口重定向功能正常
✅ 用户认证已对接管理端
✅ 用户登录页面已优化
✅ 管理端用户创建功能可用

## 下一步开发建议

1. 完善用户端Dashboard功能（域名管理、统计数据等）
2. 添加用户个人资料修改功能
3. 实现SSL证书管理
4. 添加流量统计图表
5. 实现域名配置管理界面

---

更新时间：2025-10-31
版本：v1.0.11
