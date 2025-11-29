# 方案 1：VS Code Remote-SSH + 原生环境开发指南

本文面向已经安装 VS Code 的开发者，帮助你在云服务器上直接构建 Lingcdn 的远程开发环境。流程包含云主机选型、SSH 密钥管理、服务器初始化脚本、VS Code 配置以及日常开发步骤。

---

## 1. 准备云服务器

1. **规格**：建议 4 vCPU / 8 GB RAM / 80 GB SSD 起步，方便同时运行多个 Lingcdn 子项目（Lingadmin、Lingapi、Lingnode）。  
2. **操作系统**：推荐 Ubuntu 22.04 LTS 或 AlmaLinux 9。保持内核较新，可直接使用 systemd、最新 glibc。  
3. **网络与安全**：启用公网 IP，放通 TCP 22（SSH）、必要的调试端口（例如 3000/8080/9000），其余端口默认拒绝。  
4. **账户与备份**：创建独立的非 root 账号（如 `lingcdn`），启用云平台快照/备份策略，保证环境可快速恢复。

---

## 2. 生成并上传 SSH 密钥

在本地电脑（Windows PowerShell 或 WSL 均可）执行：

```powershell
ssh-keygen -t ed25519 -C "lingcdn-dev" -f $env:USERPROFILE\.ssh\lingcdn_ed25519
```

- `.pub` 文件即公钥，登录云控制台 > SSH Key/密钥对，粘贴公钥内容。  
- 私钥文件务必妥善保存，权限设置为 600；不要上传到代码仓库。  
- 如果云服务器已创建，可通过云平台控制台或运营商提供的“重装系统 + 注入密钥”功能替换默认密码登录。

---

## 3. 服务器初始化

1. 首次登陆：使用云厂商提供的 root/默认账号通过 `ssh root@云服务器IP` 连接。  
2. 下载并执行仓库提供的初始化脚本：

```bash
curl -fsSL https://raw.githubusercontent.com/<your-repo>/tools/remote-dev/init-server.sh -o /tmp/init-server.sh
chmod +x /tmp/init-server.sh
sudo REMOTE_DEV_USER=lingcdn bash /tmp/init-server.sh
```

> 如果暂时无法从 GitHub 下载，可在本地 `tools/remote-dev/init-server.sh` 修改后使用 `scp` 上传至服务器执行。  

脚本功能（可根据需要调整环境变量 `GO_VERSION`、`NODE_MAJOR`、`TIMEZONE`）：

- 安装 `git`、`curl`、`build-essential` 等常用依赖；
- 安装 Go 1.22.x、Node.js 18 LTS（含 pnpm）；
- 创建开发用户 `lingcdn`（可自定义），并配置 `.bashrc`、工作目录；
- 设置时区为 `Asia/Shanghai`。

脚本执行完成后重新登录 `lingcdn` 用户：

```bash
ssh lingcdn@云服务器IP
```

---

## 4. VS Code Remote-SSH 配置

1. 安装扩展：在 VS Code 的扩展市场安装 `Remote - SSH`。  
2. 编辑 `~/.ssh/config`（Windows 在 `C:\Users\<user>\.ssh\config`）：

```
Host lingcdn-dev
    HostName <云服务器IP>
    User lingcdn
    Port 22
    IdentityFile C:\Users\<user>\.ssh\lingcdn_ed25519
    ForwardAgent yes
```

3. 在 VS Code 命令面板执行 `Remote-SSH: Connect to Host...`，选择 `lingcdn-dev`。首次连接会自动在服务器上安装 VS Code Server。  
4. 成功连接后，使用 `File > Open Folder` 打开 `~/workspace`。建议执行 `git clone git@github.com:<your-org>/Lingcdn系统开发.git` 将仓库拉至 `workspace`。

---

## 5. 日常开发流程

1. **同步代码**  
   ```bash
   cd ~/workspace/Lingcdn系统开发
   git pull --rebase
   ```

2. **安装依赖**  
   - Go 项目（Lingnode/Lingapi）：使用 `go mod tidy`、`go test ./...`。  
   - 前端或管理端（Lingadmin web）：若需要 Node 工具链，直接使用脚本安装的 Node 18 + pnpm。  

3. **运行/调试**  
   - 在 VS Code 中创建 Go/Node 启动配置，Remote-SSH 会在服务器上调试。  
   - 若需要数据库/缓存，可在服务器上通过 Docker 或系统包管理器安装 MySQL、Redis 等。  

4. **Git 工作流**  
   - `gs` / `gp` 别名已配置在 `.bashrc`，可快速查看状态。  
   - 建议通过 SSH Agent 转发（`ForwardAgent yes`）直接使用本地密钥签署 commit。  

---

## 6. 常见问题排查

| 问题 | 排查步骤 |
| ---- | -------- |
| VS Code 无法连接 | 确认本地能 `ssh lingcdn@IP`；检查安全组 22 端口；若提示 fingerprint 变化，删除 `known_hosts` 中旧记录。 |
| 终端语言乱码 | 在 `.bashrc` 中增加 `export LANG=zh_CN.UTF-8` 并安装 `language-pack-zh-hans`。 |
| Go/Node 版本不符 | 编辑 `tools/remote-dev/init-server.sh` 的默认版本或手动安装所需版本，重新加载 shell。 |
| 端口占用 | 通过 `sudo ss -tunlp` 查看，必要时在 `~/.ssh/config` 中配置 `LocalForward 8080 localhost:8080` 做端口映射。 |

---

## 7. 下一步建议

1. 将 `tools/remote-dev/init-server.sh` 上传到公司私有制品库，避免外网下载失败。  
2. 在云服务器上配置自动备份（快照 + 数据盘）并开启 `fail2ban`、`ufw` 等防暴力破解措施。  
3. 如果有多人协作需求，可为每位开发者生成独立 SSH Key，合理分配 sudo 权限并开启审计日志。  
4. 当开发流程稳定后，可考虑配合 `Dev Containers` 或 `docker-compose` 在服务器上构建半隔离环境，进一步降低“环境雪花”风险。
