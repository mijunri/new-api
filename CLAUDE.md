## Claude Skills（本仓库约定）

本仓库使用 **Claude Skills** 来沉淀可复用的“操作手册 + 执行约束 + 示例”。每个 skill 都是一个独立目录，核心入口文件为 `SKILL.md`。

### 目录结构

- skills 根目录：`.claude/skills/`
- 单个 skill 的入口文档：`.claude/skills/<skill_name>/SKILL.md`

示例路径（与 Claude Skills 的常见组织方式一致）：

- `.claude/skills/elevenlabs/SKILL.md`

### SKILL.md 推荐内容（轻量规范）

每个 `SKILL.md` 建议包含以下信息，方便 Claude 在需要时“按手册执行”，并减少误操作风险：

- **目的 / 能力范围**：skill 能做什么，不能做什么
- **前置条件**：依赖（如：SSH 可用、目标机器已装 docker/systemd 等）
- **输入参数**：需要哪些变量/参数（例如：环境、服务名、版本号）
- **安全与权限**：不允许的操作（例如：禁止直接打印密钥、禁止删除生产数据）
- **执行步骤（可复制）**：命令模板 + 注意事项（用占位符表达，避免写死）
- **验收方式**：怎么确认部署成功/日志正常
- **失败处理**：常见报错定位思路与回滚策略

### 如何新增一个 skill

1. 在 `.claude/skills/<skill_name>/` 下创建 `SKILL.md`
2. 用清晰的命令模板描述流程（先用占位符，连接/鉴权细节可后补）
3. 尽量做到：**可复用、可审计、可回滚**

---

## 部署约定

### 后端部署

- **部署方式**：SSH 到测试服务器，拉取代码、编译、重启 systemd 服务
- **部署脚本**：`.claude/skills/remote_server/deploy_test_env.sh`
- **凭证文件**：`.claude/skills/remote_server/password`（不入库）
- **如果缺少凭证**：向用户索取 SSH 密码

### 前端部署

- **部署方式**：构建后上传到阿里云 OSS，通过 CDN 分发
- **部署脚本**：`.claude/skills/frontend_deploy/deploy.sh`
- **凭证文件**：`.claude/skills/frontend_deploy/credentials`（不入库）
- **如果缺少凭证**：向用户索取阿里云 ACCESS_KEY_ID 和 ACCESS_KEY_SECRET

> ⚠️ **重要**：前端不是部署到服务器，而是上传到 OSS！

### 输出规范（重要）

- **默认中文输出**：对用户的说明、结论、操作步骤、风险提示一律使用中文。
- **推理呈现方式**：可以用中文给出“推理摘要/决策依据/关键假设”（可公开、可复现、可审计），但不要输出逐字的内部思考过程或详细推理链路。

