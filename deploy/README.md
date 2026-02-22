# New-API 部署说明 (Main 分支)

## 环境变量（从 test 分支迁移）

生产环境变量已记录在：
- `env.production.example` - 可复制为 `.env` 使用
- `new-api.service` - systemd 服务已包含环境变量

## 服务端依赖说明

**生产运行不依赖 Node.js/Bun/Docker**。Go 编译后的二进制已将前端静态资源嵌入，直接运行即可。

## 部署步骤

使用 Claude skill 部署：项目中已配置 `.claude/skills/deploy/`，对话时说「按 deploy skill 部署」即可按流程执行。

### 手动部署

1. 本地构建：`cd web && bun install && bun run build && cd ..` 后 `go build -o new-api`
2. 部署：`scp new-api root@47.237.158.148:/root/new-api/`，再 SSH 执行 `chmod +x new-api && systemctl restart new-api`
3. Systemd 配置见 `deploy/new-api.service`，已含生产环境变量
