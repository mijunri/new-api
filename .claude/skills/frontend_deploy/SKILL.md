## Skill: 前端部署到阿里云 OSS（Frontend Deploy to Aliyun OSS）

### 目的

将前端构建产物部署到阿里云 OSS，配合 CDN 加速分发。

---

### 部署目标

- **OSS Bucket**: `claude-code-agent-web`
- **Region**: 新加坡 (`oss-ap-southeast-1`)
- **前端目录**: `/workspace/web`
- **构建工具**: bun

---

### 凭证配置

凭证通过本地文件注入，**不会提交到仓库**：

- **凭证文件**: `.claude/skills/frontend_deploy/credentials`
- **格式**:
  ```
  ACCESS_KEY_ID=your_access_key_id
  ACCESS_KEY_SECRET=your_access_key_secret
  ```

> 该文件已被 `.gitignore` 忽略。

#### 首次配置

```bash
# 复制示例文件
cp .claude/skills/frontend_deploy/credentials.example .claude/skills/frontend_deploy/credentials

# 编辑填入真实的 AK/SK
vim .claude/skills/frontend_deploy/credentials
```

---

### 依赖

- **bun**: 前端构建工具（脚本会自动安装）
- **Python 3**: 运行 OSS 上传脚本
- **oss2**: 阿里云 OSS Python SDK（脚本会自动安装）

---

### 一键部署

```bash
# 在仓库根目录执行
bash .claude/skills/frontend_deploy/deploy.sh
```

#### 可选参数

- `SKIP_BUILD=true`: 跳过构建，只上传现有的 dist 目录
- `OSS_BUCKET`: 指定 OSS bucket 名称（默认: `claude-code-agent-web`）
- `OSS_REGION`: 指定 OSS 区域（默认: `ap-southeast-1`）

```bash
# 示例：跳过构建，直接上传
SKIP_BUILD=true bash .claude/skills/frontend_deploy/deploy.sh

# 示例：部署到其他 bucket
OSS_BUCKET=my-bucket bash .claude/skills/frontend_deploy/deploy.sh
```

---

### 部署流程

1. **加载凭证**: 从 `credentials` 文件读取 AK/SK
2. **安装依赖**: 检查并安装 bun（如果不存在）
3. **构建前端**: 执行 `bun install && bun run build`
4. **上传到 OSS**: 使用 Python oss2 SDK 上传 `dist` 目录

---

### 文件说明

| 文件 | 说明 |
|------|------|
| `SKILL.md` | 本说明文档 |
| `deploy.sh` | 一键部署脚本 |
| `upload_to_oss.py` | OSS 上传 Python 脚本 |
| `credentials.example` | 凭证文件示例 |
| `credentials` | 实际凭证文件（不入库） |

---

### 常见问题

#### 1. 构建失败

```bash
# 清理 node_modules 后重试
cd web && rm -rf node_modules && bun install && bun run build
```

#### 2. 上传失败

- 检查 AK/SK 是否正确
- 检查 bucket 名称和区域是否正确
- 检查网络连接

#### 3. CDN 缓存

部署后如果没有生效，可能需要：
- 在阿里云 CDN 控制台刷新缓存
- 或等待 CDN 缓存自动过期

---

### 注意事项

- **不要**将 `credentials` 文件提交到仓库
- 建议在部署前先在本地测试构建是否成功
- 大文件上传可能需要较长时间，请耐心等待
