# New-API 部署说明 (Main 分支)

## 环境变量（从 test 分支迁移）

生产环境变量已记录在：
- `env.production.example` - 可复制为 `.env` 使用
- `new-api.service` - systemd 服务已包含环境变量

## 服务端依赖说明

**生产运行不依赖 Node.js/Bun/Docker**。Go 编译后的二进制已将前端静态资源嵌入，直接运行即可。

## 部署步骤

### 方式一：GitHub Actions（推荐）

推送 main 或手动触发 workflow **Build & Deploy**：

- **1️⃣ 构建** → **2️⃣ 部署**，分步执行可查看进度
- 构建成功后才自动部署

**首次需配置 Secrets**：`DEPLOY_HOST`、`DEPLOY_USER`、`DEPLOY_PASSWORD`（见 `.github/workflows/build-and-deploy.yml`）

查询进度：`gh run watch` 或打开 Actions 页面

### 方式二：Claude skill / 本地

项目内有 `.claude/skills/deploy/`，对话时说「按 deploy skill 部署」即可。或按 skill 内命令手动执行。
