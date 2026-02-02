# Skill: 测试环境认证 (Test Auth)

## 目的

提供快速获取测试环境 JWT Token 的工具，方便测试需要认证的接口。

## 测试账号信息

默认测试账号（已创建，User ID: 543）：

```bash
WALLET_ADDRESS="muA17CfW4rB5bva5FRhGsi3dmUx9H95yyX"
PUBLIC_KEY="032f4b546ded4cb284dce67b30b68cd02611d0fde7bc173f07bce534008900f366"
```

## 使用方法

### 基本用法

```bash
# 使用默认配置获取 token
bash .claude/skills/test_auth/get_test_token.sh
```

### 自定义配置

```bash
# 指定服务器地址
bash .claude/skills/test_auth/get_test_token.sh --base-url http://localhost:8100

# 使用自定义钱包地址
bash .claude/skills/test_auth/get_test_token.sh --address muA17CfW4rB5bva5FRhGsi3dmUx9H95yyX
```

### 提取 Token 并使用

```bash
# 提取 token
TOKEN=$(bash .claude/skills/test_auth/get_test_token.sh 2>&1 | python3 -c "
import sys, json
content = sys.stdin.read()
json_start = content.find('{')
json_end = content.rfind('}') + 1
d = json.loads(content[json_start:json_end])
print(d['token'])
")

# 使用 token 调用接口
curl -X GET "http://47.236.240.43:8100/api/model/check_name?name=test_model" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'accept: application/json'
```

## 输出格式

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": 543,
  "address": "muA17CfW4rB5bva5FRhGsi3dmUx9H95yyX",
  "base_url": "http://47.236.240.43:8100"
}
```

## 工作原理

1. **获取 nonce**: `GET /api/pub/nonce?address={address}&public_key={public_key}`
2. **开发环境登录**: `POST /api/pub/dev/login`（不验证签名）
3. **返回 token**: JWT token 及相关信息

## 注意事项

- 开发环境登录接口不验证签名，仅用于测试环境
- Token 有效期为 24 小时
- 日志输出到 stderr，JSON 结果输出到 stdout
