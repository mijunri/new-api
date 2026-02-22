---
name: new-api-deploy
description: 当需要将 new-api 部署到生产服务器时使用。包含本地构建、上传 OSS、SCP 部署到 47.237.158.148 的完整流程与命令。
allowed-tools: Read, Grep, Glob, Run
---

# New-API 部署

当用户要求部署时，按以下步骤在本地执行。

## 前置信息

| 项 | 值 |
|----|-----|
| 服务器 | 47.237.158.148 |
| 用户 | root |
| 密码 | foxrouter@1234 |
| 部署路径 | /root/new-api |
| 服务名 | new-api (systemd) |

## 1. 本地构建

```bash
cd /workspace
git checkout main && git pull origin main

# 前端构建（需 Bun 或 npm）
cd web && bun install && DISABLE_ESLINT_PLUGIN='true' bun run build && cd ..
# 或: npm install --legacy-peer-deps && DISABLE_ESLINT_PLUGIN='true' npm run build

# 后端构建
VERSION="main-$(date +'%Y%m%d')-$(git rev-parse --short HEAD)"
go build -ldflags "-s -w -X 'github.com/QuantumNous/new-api/common.Version=$VERSION'" -o new-api
```

## 2. 上传 OSS（可选）

需配置 OSS 环境变量：`OSS_ACCESS_KEY_ID`、`OSS_ACCESS_KEY_SECRET`、`OSS_BUCKET`、`OSS_ENDPOINT`、`OSS_PATH`

```bash
ossutil64 cp new-api oss://${OSS_BUCKET}/${OSS_PATH} --force
```

## 3. 部署到服务器（SCP）

```bash
sshpass -p 'foxrouter@1234' scp -o StrictHostKeyChecking=no new-api root@47.237.158.148:/root/new-api/
sshpass -p 'foxrouter@1234' ssh -o StrictHostKeyChecking=no root@47.237.158.148 "cd /root/new-api && chmod +x new-api && systemctl restart new-api && systemctl status new-api"
```

## 4. 备选：服务端从 OSS 拉取（需服务端配置 ossutil）

```bash
sshpass -p 'foxrouter@1234' ssh -o StrictHostKeyChecking=no root@47.237.158.148 "cd /root/new-api && ossutil64 cp oss://\$OSS_BUCKET/\$OSS_PATH ./new-api --force && chmod +x new-api && systemctl restart new-api && systemctl status new-api"
```

## 5. 验证

```bash
curl -s http://47.237.158.148:3000/api/status
```

## 注意事项

- 执行前确认在 main 分支且已拉取最新代码
- systemd 配置见 `deploy/new-api.service`，已含生产环境变量
