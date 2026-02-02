#!/usr/bin/env bash
set -euo pipefail

# 测试环境获取 JWT Token 的 Skill
BASE_URL="${TEST_BASE_URL:-http://47.236.240.43:8100}"
WALLET_ADDRESS="${TEST_WALLET_ADDRESS:-muA17CfW4rB5bva5FRhGsi3dmUx9H95yyX}"
PUBLIC_KEY="${TEST_PUBLIC_KEY:-032f4b546ded4cb284dce67b30b68cd02611d0fde7bc173f07bce534008900f366}"

while [[ $# -gt 0 ]]; do
    case $1 in
        --base-url) BASE_URL="$2"; shift 2 ;;
        --address) WALLET_ADDRESS="$2"; shift 2 ;;
        --public-key) PUBLIC_KEY="$2"; shift 2 ;;
        *) echo "未知选项: $1" >&2; exit 1 ;;
    esac
done

echo "==> 获取测试 Token" >&2
echo "    Base URL: ${BASE_URL}" >&2
echo "    Address: ${WALLET_ADDRESS}" >&2

# 步骤 1: 获取 nonce
echo "步骤 1: 获取 nonce..." >&2
NONCE_RESP=$(curl -s -X GET "${BASE_URL}/api/pub/nonce?address=${WALLET_ADDRESS}&public_key=${PUBLIC_KEY}")
NONCE=$(echo "$NONCE_RESP" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('data', {}).get('nonce', ''))" 2>/dev/null)

if [[ -z "$NONCE" ]]; then
    echo "错误: 无法从响应中提取 nonce" >&2
    exit 1
fi

echo "    Nonce: ${NONCE}" >&2

# 步骤 2: 开发环境登录
echo "步骤 2: 开发环境登录..." >&2
LOGIN_RESP=$(curl -s -X POST "${BASE_URL}/api/pub/dev/login" \
    -H 'Content-Type: application/json' \
    -d "{\"address\": \"${WALLET_ADDRESS}\", \"public_key\": \"${PUBLIC_KEY}\", \"signature\": \"test_sign_任意值\", \"nonce\": \"${NONCE}\"}")

TOKEN=$(echo "$LOGIN_RESP" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('data', {}).get('token', ''))" 2>/dev/null)
REFRESH_TOKEN=$(echo "$LOGIN_RESP" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('data', {}).get('refresh_token', ''))" 2>/dev/null)
USER_ID=$(echo "$LOGIN_RESP" | python3 -c "import sys, json; d=json.load(sys.stdin); print(d.get('data', {}).get('user_id', ''))" 2>/dev/null)

if [[ -z "$TOKEN" ]]; then
    echo "错误: 无法从响应中提取 token" >&2
    exit 1
fi

echo "登录成功！" >&2

cat <<JSON
{
  "token": "${TOKEN}",
  "refresh_token": "${REFRESH_TOKEN}",
  "user_id": ${USER_ID},
  "address": "${WALLET_ADDRESS}",
  "base_url": "${BASE_URL}"
}
JSON
