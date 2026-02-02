#!/usr/bin/env bash
set -euo pipefail

# 从本地忽略文件加载密码并导出 SSHPASS，供 sshpass -e 使用。
# - 密码文件：.claude/skills/remote_server/password（已被 .gitignore 忽略）

_skill_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
_pw_file="${_skill_dir}/password"

_is_interactive="false"
if [[ -t 0 && -t 1 ]]; then
  _is_interactive="true"
fi

# 允许在非交互环境中直接通过环境变量注入（避免落盘保存密码）
if [[ -n "${SSHPASS:-}" ]]; then
  export SSHPASS
  return 0 2>/dev/null || exit 0
fi

_ensure_password_file() {
  if [[ -f "${_pw_file}" ]]; then
    return 0
  fi

  if [[ "${_is_interactive}" != "true" ]]; then
    echo "ERROR: password 文件不存在：${_pw_file}" >&2
    echo "请在本地创建该文件（内容为密码第一行），或在交互式终端执行后按提示录入密码。" >&2
    return 1
  fi

  echo "未检测到本地 password 文件，将提示你输入一次并仅保存到本地：" >&2
  echo "  ${_pw_file}" >&2
  echo "注意：该文件已被 .gitignore 忽略，不会被提交。" >&2

  local _pw=""
  # 静默输入，不回显
  read -r -s -p "请输入 SSH 密码: " _pw </dev/tty
  echo >&2
  if [[ -z "${_pw}" ]]; then
    echo "ERROR: 密码为空，已取消。" >&2
    return 1
  fi

  # 写入并收紧权限（避免被其他用户读取）
  umask 077
  printf "%s\n" "${_pw}" > "${_pw_file}"
  chmod 600 "${_pw_file}" || true
}

_ensure_password_file

# 读取第一行作为密码（允许末尾换行）
SSHPASS="$(head -n 1 "${_pw_file}")"
export SSHPASS

