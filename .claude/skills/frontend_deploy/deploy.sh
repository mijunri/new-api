#!/usr/bin/env bash
set -euo pipefail

# 前端部署脚本：构建并上传到阿里云 OSS
#
# 使用方法：
#   bash .claude/skills/frontend_deploy/deploy.sh
#
# 环境变量：
#   SKIP_BUILD=true   - 跳过构建，只上传
#   OSS_BUCKET        - OSS bucket 名称（默认: claude-code-agent-web）
#   OSS_REGION        - OSS 区域（默认: ap-southeast-1）

_skill_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
_repo_root="$(cd -- "${_skill_dir}/../../.." && pwd)"

# 默认配置
OSS_BUCKET="${OSS_BUCKET:-claude-code-agent-web}"
OSS_REGION="${OSS_REGION:-ap-southeast-1}"
WEB_DIR="${_repo_root}/web"
DIST_DIR="${WEB_DIR}/dist"

echo "=========================================="
echo "  前端部署到阿里云 OSS"
echo "=========================================="
echo "OSS Bucket: ${OSS_BUCKET}"
echo "OSS Region: ${OSS_REGION}"
echo "Web 目录: ${WEB_DIR}"
echo ""

# 1. 加载凭证
echo "==> 步骤 1: 加载凭证"
CREDENTIALS_FILE="${_skill_dir}/credentials"

if [[ ! -f "${CREDENTIALS_FILE}" ]]; then
    echo "错误: 凭证文件不存在: ${CREDENTIALS_FILE}"
    echo "请复制 credentials.example 并填入 AK/SK:"
    echo "  cp ${_skill_dir}/credentials.example ${CREDENTIALS_FILE}"
    exit 1
fi

# shellcheck source=/dev/null
source "${CREDENTIALS_FILE}"

if [[ -z "${ACCESS_KEY_ID:-}" ]] || [[ -z "${ACCESS_KEY_SECRET:-}" ]]; then
    echo "错误: ACCESS_KEY_ID 或 ACCESS_KEY_SECRET 未设置"
    exit 1
fi

echo "✓ 凭证已加载"
echo ""

# 2. 检查并安装 bun
if [[ "${SKIP_BUILD:-false}" != "true" ]]; then
    echo "==> 步骤 2: 检查构建工具"
    
    # 添加 bun 到 PATH
    export BUN_INSTALL="${HOME}/.bun"
    export PATH="${BUN_INSTALL}/bin:${PATH}"
    
    if ! command -v bun &> /dev/null; then
        echo "bun 未安装，正在安装..."
        curl -fsSL https://bun.sh/install | bash
        # 重新加载 PATH
        export PATH="${BUN_INSTALL}/bin:${PATH}"
    fi
    
    echo "✓ bun 版本: $(bun --version)"
    echo ""
    
    # 3. 构建前端
    echo "==> 步骤 3: 构建前端"
    cd "${WEB_DIR}"
    
    echo "安装依赖..."
    bun install
    
    echo "构建..."
    DISABLE_ESLINT_PLUGIN=true bun run build
    
    echo "✓ 构建完成"
    echo ""
else
    echo "==> 跳过构建（SKIP_BUILD=true）"
    echo ""
fi

# 4. 检查 dist 目录
if [[ ! -d "${DIST_DIR}" ]]; then
    echo "错误: dist 目录不存在: ${DIST_DIR}"
    echo "请先构建前端"
    exit 1
fi

# 5. 安装 Python 依赖
echo "==> 步骤 4: 检查 Python 依赖"
if ! python3 -c "import oss2" &> /dev/null; then
    echo "安装 oss2..."
    pip3 install oss2 -q
fi
echo "✓ oss2 已安装"
echo ""

# 6. 上传到 OSS
echo "==> 步骤 5: 上传到 OSS"
python3 "${_skill_dir}/upload_to_oss.py" \
    --access-key-id "${ACCESS_KEY_ID}" \
    --access-key-secret "${ACCESS_KEY_SECRET}" \
    --bucket "${OSS_BUCKET}" \
    --region "${OSS_REGION}" \
    --source-dir "${DIST_DIR}"

echo ""
echo "=========================================="
echo "  ✅ 部署完成！"
echo "=========================================="
