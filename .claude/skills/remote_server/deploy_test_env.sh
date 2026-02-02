#!/usr/bin/env bash
set -euo pipefail

# 测试环境部署（本地 push -> 远端 merge 到 test -> 启动服务）
#
# 默认目标：
# - HOST=47.236.240.43  REMOTE_USER=root  PORT=22
# - DEPLOY_DIR=/root/new-api
#
# 依赖：
# - 本地：git、sshpass、ssh
# - 远端：git、bash、go（如果使用 build 模式）、systemd（如果使用 systemd 模式）
#
# 启动方式（通过 START_METHOD 环境变量控制）：
# - systemd（默认）：使用 systemctl restart new-api 重启服务（需要预先编译好二进制）
# - build+systemd：先编译二进制文件，然后重启 systemd 服务
# - script：使用自定义启动脚本（通过 START_SCRIPT 指定）

_skill_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "${_skill_dir}/load_password.sh"

HOST="${HOST:-47.236.240.43}"
# 不要使用通用环境变量 USER（本机通常为 ubuntu），避免误连到错误账号
REMOTE_USER="${REMOTE_USER:-root}"
PORT="${PORT:-22}"
DEPLOY_DIR="${DEPLOY_DIR:-/root/new-api}"

REMOTE_TEST_BRANCH="${REMOTE_TEST_BRANCH:-test}"
LOCAL_BRANCH="${LOCAL_BRANCH:-$(git rev-parse --abbrev-ref HEAD)}"
REMOTE_NAME="${REMOTE_NAME:-origin}"

# 启动方式：systemd | build+systemd | script
START_METHOD="${START_METHOD:-systemd}"

# 启动脚本：仅在 START_METHOD=script 时使用，在远端仓库目录下的相对路径
START_SCRIPT="${START_SCRIPT:-./scripts/start.sh}"

# systemd 服务名：仅在 START_METHOD=systemd 或 build+systemd 时使用
SERVICE_NAME="${SERVICE_NAME:-new-api}"

# 二进制文件路径（build+systemd 模式编译后的输出路径）
BINARY_PATH="${BINARY_PATH:-./new-api}"

echo "==> 本地 push 分支：${LOCAL_BRANCH}"
git push -u "${REMOTE_NAME}" "${LOCAL_BRANCH}"

echo "==> 远端 merge 到 ${REMOTE_TEST_BRANCH} 并启动（启动方式：${START_METHOD}）"
SSHPASS="${SSHPASS}" sshpass -e ssh -p "${PORT}" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
  -o PreferredAuthentications=password -o PubkeyAuthentication=no \
  "${REMOTE_USER}@${HOST}" \
  "set -euo pipefail
   cd '${DEPLOY_DIR}'
   # 避免首次连接 GitHub 时触发 Host key verification failed
   export GIT_SSH_COMMAND='ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null'
   git remote -v
   git fetch '${REMOTE_NAME}' --prune

   # 确保 test 分支存在：优先跟踪 origin/test；否则从 origin/main 创建
   if git show-ref --verify --quiet refs/remotes/${REMOTE_NAME}/${REMOTE_TEST_BRANCH}; then
     git checkout -B '${REMOTE_TEST_BRANCH}' '${REMOTE_NAME}/${REMOTE_TEST_BRANCH}'
   else
     git checkout -B '${REMOTE_TEST_BRANCH}' '${REMOTE_NAME}/main'
   fi

   # 合并待发布分支（使用远端分支，避免服务器本地分支漂移）
   git merge --no-edit '${REMOTE_NAME}/${LOCAL_BRANCH}'

   # 根据启动方式执行相应的启动命令
   echo \"==> 启动服务（方式：${START_METHOD}）\"
   case \"${START_METHOD}\" in
     systemd)
       if ! systemctl is-enabled \"${SERVICE_NAME}\" >/dev/null 2>&1; then
         echo \"WARNING: systemd service ${SERVICE_NAME} not found or not enabled\" >&2
         echo \"请确保已配置 systemd service: /etc/systemd/system/${SERVICE_NAME}.service\" >&2
       fi
       systemctl daemon-reload || true
       systemctl restart \"${SERVICE_NAME}\"
       echo '==> 检查服务状态'
       systemctl status \"${SERVICE_NAME}\" --no-pager -l || true
       ;;
     build+systemd|build-systemd)
       echo '==> 编译前端...'
       if command -v bun >/dev/null 2>&1; then
         cd web && bun install && DISABLE_ESLINT_PLUGIN='true' VITE_REACT_APP_VERSION=\$(cat ../VERSION) bun run build && cd ..
       else
         echo 'WARNING: bun not found, skipping frontend build' >&2
       fi
       echo '==> 编译后端...'
       if ! command -v go >/dev/null 2>&1; then
         echo \"ERROR: go command not found, cannot build binary\" >&2
         exit 2
       fi
       go mod download
       VERSION=\$(cat VERSION 2>/dev/null || echo 'dev')
       CGO_ENABLED=0 go build -ldflags \"-s -w -X 'github.com/QuantumNous/new-api/common.Version='\"\${VERSION}\"\"\" -o \"${BINARY_PATH}\"
       if [[ ! -f \"${BINARY_PATH}\" ]]; then
         echo \"ERROR: build failed, binary not found: ${BINARY_PATH}\" >&2
         exit 2
       fi
       chmod +x \"${BINARY_PATH}\"
       echo \"==> 二进制文件已编译: ${BINARY_PATH}\"
       if ! systemctl is-enabled \"${SERVICE_NAME}\" >/dev/null 2>&1; then
         echo \"WARNING: systemd service ${SERVICE_NAME} not found or not enabled\" >&2
         echo \"请确保已配置 systemd service: /etc/systemd/system/${SERVICE_NAME}.service\" >&2
       fi
       systemctl daemon-reload || true
       systemctl restart \"${SERVICE_NAME}\"
       echo '==> 检查服务状态'
       systemctl status \"${SERVICE_NAME}\" --no-pager -l || true
       ;;
     script)
       if [[ ! -f \"${START_SCRIPT}\" ]]; then
         echo \"ERROR: START_SCRIPT not found: ${START_SCRIPT}\" >&2
         exit 2
       fi
       if [[ ! -x \"${START_SCRIPT}\" ]]; then
         chmod +x \"${START_SCRIPT}\" || true
       fi
       \"${START_SCRIPT}\"
       ;;
     *)
       echo \"ERROR: 未知的启动方式: ${START_METHOD}\" >&2
       echo '支持的方式: systemd, build+systemd, script' >&2
       exit 2
       ;;
   esac

   echo '==> 部署完成'
   git status -sb
   git log -1 --oneline
  "

