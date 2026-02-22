# New-API 部署说明 (Main 分支)

## 环境变量（从 test 分支迁移）

生产环境变量已记录在：
- `env.production.example` - 可复制为 `.env` 使用
- `new-api.service` - systemd 服务已包含环境变量

## 服务端依赖说明

**生产运行不依赖 Node.js/Bun/Docker**。Go 编译后的二进制已将前端静态资源嵌入，直接运行即可。

构建由 **GitHub Actions** 完成，推送 main 分支或手动触发工作流即可生成 Linux amd64 二进制。

## 部署步骤（推荐：GitHub Actions 构建 + 二进制运行）

### 1. 获取构建产物

- 推送代码到 main 分支后，前往 **GitHub → Actions → "Build Main (Linux amd64)"**
- 打开最新的成功运行，在 **Artifacts** 中下载 `new-api-linux-amd64`
- 解压得到 `new-api` 可执行文件

或使用 GitHub API 下载最新构建（需 token）：

```bash
# 获取最新 workflow run 的 artifact 下载 URL
# 可在 Actions 页面手动下载
```

### 2. 部署到服务器

```bash
# 上传二进制到服务器
scp new-api root@47.237.158.148:/root/new-api/

# SSH 到服务器
ssh root@47.237.158.148

# 停止旧服务，替换二进制，启动
cd /root/new-api
systemctl stop new-api
chmod +x new-api
systemctl start new-api
```

### 3. Systemd 配置

确保 `/etc/systemd/system/new-api.service` 使用 `deploy/new-api.service` 中的配置（已包含生产环境变量）。
