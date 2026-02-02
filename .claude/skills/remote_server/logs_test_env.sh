#!/usr/bin/env bash
set -euo pipefail

# 测试环境日志查看（通过 sshpass 在远端执行 tail/journalctl）
#
# 支持的日志模式：
# - journalctl（默认）：systemd 服务日志
# - file：文件日志（适用于直接运行二进制文件）

_skill_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "${_skill_dir}/load_password.sh"

HOST="${HOST:-47.236.240.43}"
# 不要使用通用环境变量 USER（本机通常为 ubuntu），避免误连到错误账号
REMOTE_USER="${REMOTE_USER:-root}"
PORT="${PORT:-22}"
DEPLOY_DIR="${DEPLOY_DIR:-/root/new-api}"

LOG_MODE="${LOG_MODE:-journalctl}" # journalctl | file
LINES="${LINES:-200}"
FOLLOW="${FOLLOW:-false}"

# journalctl mode：systemd 服务名
APP_NAME="${APP_NAME:-new-api}"

# file mode：文件日志路径
LOG_FILE="${LOG_FILE:-${DEPLOY_DIR}/logs/new-api.log}"

SINCE="${SINCE:-}"

remote_cmd=""
if [[ "${LOG_MODE}" == "journalctl" ]]; then
  # systemd journalctl 模式
  if [[ -n "${SINCE}" ]]; then
    remote_cmd="journalctl -u \"${APP_NAME}\" --since \"${SINCE}\" --no-pager"
  else
    remote_cmd="journalctl -u \"${APP_NAME}\" -n \"${LINES}\" --no-pager"
  fi
  if [[ "${FOLLOW}" == "true" ]]; then
    remote_cmd="journalctl -u \"${APP_NAME}\" -f"
  fi
else
  # file 模式：直接读取日志文件
  if [[ "${FOLLOW}" == "true" ]]; then
    remote_cmd="if [[ ! -f \"${LOG_FILE}\" ]]; then echo \"ERROR: LOG_FILE not found: ${LOG_FILE}\" >&2; ls -la \"$(dirname "${LOG_FILE}")\" || true; exit 2; fi; tail -f \"${LOG_FILE}\""
  else
    remote_cmd="if [[ ! -f \"${LOG_FILE}\" ]]; then echo \"ERROR: LOG_FILE not found: ${LOG_FILE}\" >&2; ls -la \"$(dirname "${LOG_FILE}")\" || true; exit 2; fi; tail -n \"${LINES}\" \"${LOG_FILE}\""
  fi
fi

SSHPASS="${SSHPASS}" sshpass -e ssh -p "${PORT}" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
  -o PreferredAuthentications=password -o PubkeyAuthentication=no \
  "${REMOTE_USER}@${HOST}" \
  "${remote_cmd}"

