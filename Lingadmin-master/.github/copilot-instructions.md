<!--
自动生成给 AI 编程助手的快速指导文件，目标：让模型快速理解本仓库的架构、常用工作流和查找点。
请保持简短、具体、以项目为准，不写通用编程建议。
-->
# Copilot 使用说明（给 AI 编程助手）

下面是对 `Lingadmin` / `EdgeAdmin` 仓库的要点速览与可直接执行的示例，帮助你快速定位重要代码路径并输出有用改动。

- **项目类型**：Go 后端管理平台（前端静态资源在 `web/`），模块名 `github.com/TeaOSLab/EdgeAdmin`（见 `go.mod`）。
- **三大组件（上下游）**：
  - 管理平台（本仓库，EdgeAdmin）——入口 `cmd/lingcdnadmin/main.go`。
  - API 节点（EdgeAPI）——构建/分发时会引用，构建脚本期待它位于仓库旁的 `../EdgeAPI`。
  - 公共库（EdgeCommon）——`go.mod` 使用 `replace` 指向 `../EdgeCommon`。

- **核心目录（先看这些）**
  - `cmd/lingcdnadmin/main.go`：程序入口，定义 CLI、启动/daemon/upgrade 等命令。
  - `internal/`：业务代码主干（`apps`, `configs`, `nodes`, `web`, `gen` 等）。
  - `web/`：前端静态资源与视图模板。
  - `configs/`：运行时配置示例，主要 `configs/server.yaml` / `configs/server.template.yaml`。
  - `build/`：发布构建脚本，重要脚本 `build/build.sh`（构建流程、依赖 EdgeAPI）。
  - `docker/`：容器化相关文件（`Dockerfile`,`run.sh`）。

- **重要工程约定与模式**
  - CLI & 进程控制：主程序用 `apps.NewAppCmd()` 注册动作（`start`,`stop`,`daemon`,`dev`,`prod`,`upgrade`,`generate` 等），很多行为通过临时 socket (gosock) 与运行中进程通信（查看 `internal/nodes` 和 `github.com/iwind/gosock` 的使用）。
  - 配置优先：运行/发布依赖 `configs/server.yaml`（模板为 `configs/server.template.yaml`）。本地调试可直接编辑 `configs/server.yaml`。
  - 构建依赖兄弟仓：`build/build.sh` 会调用 `../../EdgeAPI/build/build.sh` 并把 EdgeAPI 的 zip 解压到发布包。因此本地完整构建需要把 `EdgeAPI` 与 `EdgeCommon` 放到与当前仓库同级目录，或调整 `go.mod` 的 `replace`。
  - 版本与常量集中在 `internal/const/const.go`（例如 `ProcessName`、`Version`、`ProductName`）。谨慎修改。

- **常用开发命令（在 Windows PowerShell 下示例）**
  - 快速运行开发服务器（直接用 `go`）：
    ```powershell
    go run .\cmd\lingcdnadmin\main.go
    # 或传参数切换环境：
    go run .\cmd\edge-admin\main.go dev
    ```
  - 本地编译二进制（linux/交叉编译请使用 `build/build.sh`）：
    ```powershell
    go build -o .\bin\lingcdnadmin .\cmd\lingcdnadmin\main.go
    .\bin\lingcdnadmin start
    ```
  - 使用仓库自带发布脚本（在类 Unix 环境或 WSL）：
    ```bash
    ./build/build.sh linux amd64 community
    ```
  - 运行生成步骤（构建过程一部分）：
    ```powershell
    go run .\cmd\lingcdnadmin\main.go generate
    ```
  - 运行测试（若有）：
    ```powershell
    go test ./...
    ```

- **跨仓库/集成注意事项**
  - `go.mod` 中 `replace github.com/TeaOSLab/EdgeCommon => ../EdgeCommon`：如果缺少本地依赖，请 clone 对应仓库到上一级目录或删除/调整 replace。
  - 构建脚本会调用 `EdgeAPI` 的构建输出；对完整发行包的改动往往需要同时更新 `EdgeAPI`。
  - 新增：`build/build.sh` 现在接受 `EDGEAPI_PATH` 环境变量来指定 `EdgeAPI` 的位置（可用在 CI 或非标准布局下）。

- **命名与一致性（重要）**
  - 仓库中存在多套名称（历史遗留）：
    - `EdgeAdmin` / `edge-admin`：项目名或文档中可能出现的名称。
    - `ling-admin`：用于进程/脚本/systemd 的运行时服务名（`ProcessName`）。
    - `LingCDN`：产品名称，用于 UI/展示层（`ProductName`）。

  - 建议：
    - 程序内部请使用 `internal/const/const.go` 中的常量（例如 `ProcessName`, `ProductName`）作为唯一可信来源。
    - `ProcessName` 用于套接字名称、守护脚本与服务注册；`ProductName`/`ProductNameZH` 用于用户可见文本。
    - 新增全局常量时，在 `internal/const/const.go` 添加注释说明，避免出现多个散落的字符串常量。

- **快速定位审查点（PR/变更时优先看）**
  - API/协议改动：`internal/gen`、`rpc/`、`internal/nodes`。
  - 启动/部署改动：`cmd/lingcdnadmin/main.go`、`internal/const/const.go`、`configs/*`、`build/*`。
  - 前端改动与静态资源：`web/public/js`（`components.src.js` → 压缩为 `components.js`），`web/views`。

- **常见 pitfalls（项目特有）**
  - 不要只改 `build/build.sh` 假定 EdgeAPI 在线获取；发布脚本期望本地 sibling 仓库存在。修改构建流程时同时考虑 `EdgeAPI`。
  - IPC 使用 `gosock`，一些 CLI 操作会向正在运行进程发送命令（如 `dev`/`prod`/`demo`），测试时请注意是否已有运行进程监听 socket。

如果你想我合并到现有 `.github/copilot-instructions.md`（如果仓库已有），或把内容改成更简洁 / 更详细的版本，请告诉我需要强调的区域（例如：构建、测试、特定子系统）。

--
_修改建议：若有项目内部约定（例如私有 CI、远端环境变量、额外 setup 步骤），请补充，我会把它们合并进本文件。_
