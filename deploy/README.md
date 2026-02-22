# New-API 部署说明 (Main 分支)

## 环境变量（从 test 分支迁移）

生产环境变量已记录在：
- `env.production.example` - 可复制为 `.env` 使用
- `new-api.service` - systemd 服务已包含环境变量

## 服务端依赖说明

**生产运行不依赖 Node.js/Bun**。Go 编译后的二进制已将前端静态资源嵌入，直接运行即可。

构建阶段才需要 Node/Bun（用于 `web/` 的 vite build），可选方式：
1. **Docker 构建**：Dockerfile 内自动完成前端构建，服务器只需 Docker
2. **预构建二进制**：在 CI 或本地构建后上传到服务器

## 部署步骤（二进制方式）

```bash
# 1. 拉取 main 分支
cd /root/new-api
git fetch origin
git checkout main
git pull origin main

# 2. 构建（需安装 Go + Bun，或使用 Docker 构建）
# 方式 A: Docker 构建
docker build -t new-api:main .
docker run -d --rm -p 3000:3000 -v $(pwd)/data:/data --env-file .env new-api:main

# 方式 B: 本地构建
cd web && bun install && bun run build && cd ..
go build -ldflags "-s -w" -o new-api

# 3. 更新 systemd 配置
sudo cp deploy/new-api.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl restart new-api
```

## 部署步骤（Docker Compose）

若使用 Docker，需将 `deploy/env.production.example` 复制为 `.env`，并修改 `docker-compose.yml` 中的 `environment` 使用外部 MySQL/Redis 地址。
