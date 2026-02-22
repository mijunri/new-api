---
name: new-api-deploy
description: 当需要将 new-api 部署到生产服务器时使用。支持 GitHub Actions 全自动部署，或本地构建后 SCP 部署。
allowed-tools: Read, Grep, Glob, Run
---

# New-API 部署

## 方式一：GitHub Actions（推荐，可查看进度）

推送 main 或手动触发 workflow `Build & Deploy`，分两步执行：

- **1️⃣ 构建**：前端 + 后端，进度在 Actions 页可见
- **2️⃣ 部署**：构建成功后自动 SCP 上传并重启服务

### 首次配置（GitHub prod 环境 Secrets）

在 **Settings → Environments → prod** 下配置：

| Secret | 值 |
|--------|-----|
| DEPLOY_HOST | 47.237.158.148 |
| DEPLOY_USER | root |
| DEPLOY_PASSWORD | foxrouter@1234 |

或使用 SSH 私钥：`DEPLOY_SSH_KEY`（与密码二选一）

### 查询进度

```bash
# 需要 gh CLI 且已登录
gh run list --workflow=build-and-deploy.yml --limit 1
gh run watch  # 实时跟踪最新一次运行
```

或直接打开：**GitHub → Actions → Build & Deploy**

---

## 方式二：本地构建 + SCP

当 Actions 不可用或需手动控制时使用。

| 项 | 值 |
|----|-----|
| 服务器 | 47.237.158.148 |
| 用户 | root |
| 密码 | foxrouter@1234 |
| 部署路径 | /root/new-api |

### 1. 本地构建

```bash
cd /workspace
git checkout main && git pull origin main
cd web && bun install && DISABLE_ESLINT_PLUGIN='true' bun run build && cd ..
VERSION="main-$(date +'%Y%m%d')-$(git rev-parse --short HEAD)"
go build -ldflags "-s -w -X 'github.com/QuantumNous/new-api/common.Version=$VERSION'" -o new-api
```

### 2. 部署

```bash
sshpass -p 'foxrouter@1234' scp -o StrictHostKeyChecking=no new-api root@47.237.158.148:/root/new-api/
sshpass -p 'foxrouter@1234' ssh -o StrictHostKeyChecking=no root@47.237.158.148 "cd /root/new-api && chmod +x new-api && systemctl restart new-api && systemctl status new-api"
```

### 3. 验证

```bash
curl -s http://47.237.158.148:3000/api/status
```

---

systemd 配置见 `deploy/new-api.service`
