# New API 项目启动指南

## 快速启动

### 1. 启动后端服务

```bash
# 在项目根目录执行
cd new-api

# 方式一：直接运行（开发模式）
go run main.go

# 方式二：构建后运行
go build -o new-api && ./new-api
```

后端默认运行在：`http://localhost:3000`

### 2. 启动前端开发服务器（可选）

```bash
# 在 web 目录执行
cd web

# 启动 Vite 开发服务器（支持热更新）
npm run dev
```

前端开发服务器运行在：`http://localhost:5173`

---

## 完整启动流程

### 第一次启动

```bash
# 1. 检查环境配置
cat .env
# 如果不存在，会使用默认配置（SQLite 本地数据库）

# 2. 启动后端（会自动创建 root 用户）
go run main.go

# 3. 访问系统
# 浏览器打开：http://localhost:3000
# 使用 root / 123456 登录
```

### 日常开发

```bash
# 终端 1：启动后端
cd new-api
go run main.go

# 终端 2：启动前端（可选，支持热更新）
cd new-api/web
npm run dev
```

然后访问 http://localhost:5173 进行开发。

---

## 配置说明

### 环境变量

| 文件 | 用途 |
|------|------|
| `.env` | 本地开发配置（默认使用 SQLite） |
| `prod.env` | 生产环境配置（MySQL） |
| `test.env` | 测试环境配置（MySQL） |

### 创建 .env 文件（可选）

```bash
# 使用 SQLite（默认，无需配置）
SQL_DSN=local
PORT=3000

# 或使用 MySQL
SQL_DSN=bitmodel:bitmodel@1234@tcp(47.236.240.43:3306)/bitmodel?parseTime=true
PORT=3000
```

---

## 前端配置

| 文件 | API 地址 |
|------|----------|
| `web/.env.development` | 空（走 Vite 代理到 localhost:3000） |
| `web/.env.production` | https://test-api.foxrouter.com |

---

## 默认账号

| 用户名 | 密码 | 权限 |
|--------|------|------|
| `root` | `123456` | 超级管理员（系统自动创建） |

---

## 常见问题

### 端口被占用

```bash
# 查找占用 3000 端口的进程
lsof -i :3000

# 杀掉进程
kill -9 <PID>
```

### 数据库问题

```bash
# SQLite 数据库位置
ls -la one-api.db

# 重置数据库（删除后重启服务会自动创建）
rm one-api.db
```

### 前端构建失败

```bash
cd web

# 清理依赖重新安装
rm -rf node_modules package-lock.json
npm install --legacy-peer-deps

# 构建
npm run build
```

---

## 服务管理

### 查看服务状态

```bash
# 后端健康检查
curl http://localhost:3000/api/status

# 查看进程
ps aux | grep "go run main.go" | grep -v grep
ps aux | grep "vite" | grep -v grep
```

### 停止服务

```bash
# 如果使用 go run，直接 Ctrl+C 停止

# 或查找并杀掉进程
pkill -f "go run main.go"
pkill -f "npm run dev"
```

---

## 部署到测试服务器

```bash
# 使用部署脚本
bash .claude/skills/remote_server/deploy_test_env.sh

# 或手动部署
ssh root@47.236.240.43
cd /root/new-api
git checkout test
git merge --no-edit origin/main
CGO_ENABLED=0 go build -ldflags '-s -w' -o new-api
systemctl restart new-api
systemctl status new-api
```

---

## 项目结构

```
new-api/
├── main.go              # 后端入口
├── .env                 # 环境配置
├── one-api.db           # SQLite 数据库（自动生成）
├── web/
│   ├── src/            # 前端源码
│   ├── dist/           # 前端构建产物
│   ├── .env.development # 开发环境配置
│   └── .env.production # 生产环境配置
└── .claude/skills/      # 部署技能
    ├── remote_server/  # 服务器部署
    └── frontend_deploy/ # 前端部署到 OSS
```

---

## 技能说明

项目内置了 Claude Code 技能用于部署：

### 远程服务器部署

- **技能目录**：`.claude/skills/remote_server/`
- **功能**：部署到测试服务器、查看日志
- **密码文件**：`.claude/skills/remote_server/password`

```bash
# 部署到测试环境
bash .claude/skills/remote_server/deploy_test_env.sh

# 查看测试环境日志
bash .claude/skills/remote_server/logs_test_env.sh

# 实时跟踪日志
FOLLOW=true bash .claude/skills/remote_server/logs_test_env.sh
```

### 前端部署到 OSS

- **技能目录**：`.claude/skills/frontend_deploy/`
- **功能**：上传前端到阿里云 OSS
- **凭证文件**：`.claude/skills/frontend_deploy/credentials`

```bash
# 部署前端到 OSS
bash .claude/skills/frontend_deploy/deploy.sh
```

---

**最后更新**：2026-02-07
