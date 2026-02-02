## Skill: 远程服务器部署与日志查看（Remote Server Ops）

### 目的

提供一个最小可用的远程运维 skill 框架，覆盖两类常用动作：

- **部署（deploy）**：将当前版本发布到远程服务器并完成重启/健康检查
- **日志查看（logs）**：快速查看服务日志用于排障

本 skill **使用 `sshpass` + `ssh` 在远程服务器执行命令**来完成部署与日志查看。

---

### 能力范围（本期）

- **deploy**：在远程服务器上拉取/更新代码或镜像，执行启动脚本，做基本验收
- **logs**：按服务名/时间范围查看日志（支持 follow）
- **deploy_test_env**：测试环境部署（本地 push 当前分支 → 远端 merge 到 `test` 分支 → 执行启动脚本）
- **logs_test_env**：测试环境日志查看（文件 tail 或 journalctl）

### 明确不做（避免越权/误操作）

- 不会把任何真实账号密码/私钥/token **以明文写入并提交到仓库**
- 不会执行破坏性操作（如：删除数据库/清空目录）除非你明确要求并提供确认步骤
- 不会默认对生产环境执行变更；需要明确指定环境（dev/staging/prod）

---

### 测试服务器（已固化默认值）

- `ENV`: `test`
- `HOST`: `47.236.240.43`
- `USER`: `root`
- `PORT`: `22`
- `DEPLOY_DIR`: `/root/new-api`（约定：服务器上 new-api 仓库目录；如实际不同可覆盖）
- `REPO_URL`: `https://github.com/QuantumNous/new-api.git`

> 密码将从本地文件注入：`.claude/skills/remote_server/password`（该文件已被 `.gitignore` 忽略，不会提交到仓库）。
>
> - **方式 A（推荐，交互式）**：直接执行 `source ".claude/skills/remote_server/load_password.sh"`，若检测到 `password` 文件不存在，会提示你**静默输入**一次并写入本地 `password` 文件（权限 `600`）。
> - **方式 B（手动）**：复制示例文件：`.claude/skills/remote_server/password.example` → `password`，把密码写到第一行；然后用 `source ".claude/skills/remote_server/load_password.sh"` 加载为 `SSHPASS`

---

### 约定的输入参数

以下参数可通过“对话输入”或环境变量提供（具体采集方式后续统一）：

- `ENV`：`test` | `dev` | `staging` | `prod`
- `HOST`：目标主机
- `PORT`：SSH 端口（默认 22）
- `USER`：SSH 用户
- `SSHPASS`：SSH 密码（给 `sshpass -e` 用；默认从 `password` 文件加载）
- `APP_NAME`：服务/应用名（例如：`new-api`）
- `DEPLOY_DIR`：远程仓库/部署目录（默认：`/root/new-api`）
- `REPO_URL`：仓库地址（默认：`https://github.com/QuantumNous/new-api.git`）
- `START_METHOD`：启动方式（`systemd` | `build+systemd` | `script`，默认：`systemd`）
- `SERVICE_NAME`：systemd 服务名（仅在 `START_METHOD=systemd` 或 `build+systemd` 时使用，默认：`new-api`）
- `BINARY_PATH`：二进制文件路径（仅在 `START_METHOD=build+systemd` 时使用，默认：`./new-api`）
- `START_SCRIPT`：启动脚本路径（仅在 `START_METHOD=script` 时使用，默认：`./scripts/start.sh`）
- `VERSION`：版本号/分支/commit（例如：`main` 或 `v0.1.0`）

---

### 连接与鉴权（sshpass）

#### 前置条件

- 本机/执行环境已安装：`sshpass`、`ssh`
- 远程服务器允许密码登录（测试服务器默认允许）
- 已准备本地密码文件：`.claude/skills/remote_server/password`（或使用 `load_password.sh` 交互式生成）
- 已加载：`source ".claude/skills/remote_server/load_password.sh"`（会将密码导出为 `SSHPASS`）

#### 建议的 SSH 选项

- `-o StrictHostKeyChecking=no`（便于自动化；测试机可用）
- `-o UserKnownHostsFile=/dev/null`（避免污染 known_hosts；测试机可用）

#### 远程执行命令模板（推荐）

把要在远程执行的命令放进同一个字符串里：

- `SSHPASS="${SSHPASS}" sshpass -e ssh -p "${PORT}" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null "${USER}@${HOST}" "<REMOTE_CMD>"`

> 说明：`sshpass -e` 会从环境变量 `SSHPASS` 读取密码。请在运行前 `export SSHPASS='***'`，不要把密码写进仓库文件。

#### 本地便捷函数（可复制）

```bash
remote_exec() {
  SSHPASS="${SSHPASS}" sshpass -e ssh -p "${PORT:-22}" \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    "${USER}@${HOST}" "$@"
}
```

---

### 操作：deploy（基于远程命令执行）

#### 目标

把指定 `VERSION` 部署到 `ENV` 环境的目标服务器，并完成基本验收。

#### 步骤（命令模板）

0. 初始化（仅首次需要）
   - 若 `${DEPLOY_DIR}` 不存在或不是 git 仓库，则 clone：
     - `mkdir -p ${DEPLOY_DIR}`
     - `git clone ${REPO_URL} ${DEPLOY_DIR}`
1. 进入部署目录
   - `cd ${DEPLOY_DIR}`
2. 获取/更新版本（按你最终的发布形态二选一）
   - **代码发布（git）**：
     - `git fetch --all --prune`
     - `git checkout ${VERSION}`
     - `git pull --ff-only`
   - **二进制发布（build+systemd）**：
     - 在远程服务器编译：`go build -o new-api`
3. 安装/构建（按项目需要）
   - Python 示例（占位）：`pip install -r requirements.txt`
4. 构建/重启服务（按 `START_METHOD`）
   - systemd 示例：`sudo systemctl restart ${APP_NAME}`
   - build+systemd 示例：先编译二进制，然后 `sudo systemctl restart ${APP_NAME}`
5. 验收（最小集）
   - `sudo systemctl status ${APP_NAME} --no-pager`（若 systemd）
   - `curl -fsS http://127.0.0.1:<PORT>/health`（若提供健康检查）

#### 一键执行（示例模板）

```bash
# 0) 加载测试机密码（本地文件，不入库）
source ".claude/skills/remote_server/load_password.sh"

# 1) 提供连接信息（测试机默认值可直接用）
export ENV="test"
export HOST="47.236.240.43"
export USER="root"
export PORT="22"
export DEPLOY_DIR="/root/new-api"
export REPO_URL="https://github.com/QuantumNous/new-api.git"

# 2) 业务参数
export VERSION="main"
export APP_NAME="new-api"
export START_METHOD="systemd"  # 或 build+systemd、script

# 3) 执行：远程更新代码 + 重启（重启方式后续你确定后再收敛）
SSHPASS="${SSHPASS}" sshpass -e ssh -p "${PORT}" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
  "${USER}@${HOST}" \
  "set -euo pipefail; \
   if [[ ! -d '${DEPLOY_DIR}/.git' ]]; then mkdir -p '${DEPLOY_DIR}'; git clone '${REPO_URL}' '${DEPLOY_DIR}'; fi; \
   cd '${DEPLOY_DIR}'; \
   git fetch --all --prune; git checkout '${VERSION}'; git pull --ff-only; \
   echo 'UPDATED';"
```

#### 失败处理（占位）

- 若启动失败：先 `logs` 查看最近 200 行，再定位配置/依赖/端口
- 若需要回滚：回到上一个可用 `VERSION`，重复部署步骤（后续细化）

---

### 操作：logs（基于远程命令执行）

#### 目标

按应用名快速查看日志，支持追踪最新日志（follow）。

#### 常用参数（占位符）

- `LINES`：默认 `200`
- `FOLLOW`：`true` | `false`
- `SINCE`：例如 `10m` / `1h` / `2026-01-20 10:00:00`

#### 命令模板（按 `SERVICE_MANAGER`）

- systemd（journalctl）：
  - 最近日志：`journalctl -u ${APP_NAME} -n ${LINES} --no-pager`
  - 按时间：`journalctl -u ${APP_NAME} --since "${SINCE}" --no-pager`
  - 跟随：`journalctl -u ${APP_NAME} -f`

- 文件日志（file mode）：
  - 最近日志：`tail -n ${LINES} ${LOG_FILE}`
  - 跟随：`tail -f ${LOG_FILE}`

#### 一键执行（systemd / journalctl 示例）

```bash
source ".claude/skills/remote_server/load_password.sh"

export HOST="47.236.240.43"
export USER="root"
export PORT="22"
export APP_NAME="new-api"
export LINES="200"

SSHPASS="${SSHPASS}" sshpass -e ssh -p "${PORT}" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
  "${USER}@${HOST}" \
  "journalctl -u '${APP_NAME}' -n '${LINES}' --no-pager"
```

---

### 操作：deploy_test_env（测试环境：push → merge 到 test → 启动脚本）

#### 目标

把**当前本地分支** push 到 GitHub，然后在测试服务器上把该分支合并到 `test` 分支，最后执行启动脚本完成测试环境部署。

#### 核心约定

- **本地**：你正在开发的分支（默认自动识别当前分支）
- **远端服务器**：在 `${DEPLOY_DIR}` 目录已 clone 该仓库（我们已完成）
- **测试分支**：`test`（可通过 `REMOTE_TEST_BRANCH` 覆盖）
- **启动方式**：默认 `systemd`（可通过 `START_METHOD` 覆盖）
  - `systemd`（默认）：使用 `systemctl restart new-api` 重启服务（需要预先编译好二进制文件）
  - `build+systemd`：先编译前端和后端二进制文件，然后重启 systemd 服务（需要远端服务器安装 go 和 bun）
  - `script`：执行自定义启动脚本（通过 `START_SCRIPT` 指定，默认：`./scripts/start.sh`）

#### 一键执行（推荐：用脚本）

仓库内已提供可直接执行的脚本：

- `.claude/skills/remote_server/deploy_test_env.sh`

```bash
# 在仓库根目录执行
bash .claude/skills/remote_server/deploy_test_env.sh
```

可选环境变量：

- `LOCAL_BRANCH`：要发布的本地分支（默认当前分支）
- `REMOTE_TEST_BRANCH`：测试环境目标分支（默认 `test`）
- `START_METHOD`：启动方式（`systemd` | `build+systemd` | `script`，默认 `systemd`）
- `START_SCRIPT`：远端启动脚本（仅在 `START_METHOD=script` 时使用，默认 `./scripts/start.sh`）
- `SERVICE_NAME`：systemd 服务名（仅在 `START_METHOD=systemd` 或 `build+systemd` 时使用，默认 `new-api`）
- `BINARY_PATH`：二进制文件路径（仅在 `START_METHOD=build+systemd` 时使用，默认 `./new-api`）
- `DEPLOY_DIR`：远端仓库目录（默认 `/root/new-api`）
- `HOST/USER/PORT`：默认已指向测试服务器

---

### 操作：logs_test_env（测试环境日志）

#### 一键执行（推荐：用脚本）

- `.claude/skills/remote_server/logs_test_env.sh`

```bash
# 文件日志（默认）
bash .claude/skills/remote_server/logs_test_env.sh

# 跟随最新日志
FOLLOW=true bash .claude/skills/remote_server/logs_test_env.sh

# 指定日志文件（file 模式）
LOG_MODE=file LOG_FILE="/root/new-api/logs/new-api.log" bash .claude/skills/remote_server/logs_test_env.sh

# systemd / journalctl 模式（默认）
LOG_MODE=journalctl APP_NAME="new-api" bash .claude/skills/remote_server/logs_test_env.sh
```

> 注意：
> - 日志查看方式取决于实际部署方式：
>   - systemd（默认）：使用 `journalctl -u new-api` 查看
>   - 直接运行：日志文件通常在 `${DEPLOY_DIR}/logs/` 目录下，使用 file 模式
> - 项目为 Go 语言项目，使用 systemd 部署时，日志通过 journalctl 查看
> - 使用自己的 Redis 和 MySQL 时，确保环境变量 `REDIS_CONN_STRING` 和 `SQL_DSN` 已正确配置

### 产出与记录（建议）

- 每次 deploy 记录：`ENV / HOST / VERSION / 时间 / 结果`
- 若发生故障：记录日志片段与最终修复动作（避免重复踩坑）

