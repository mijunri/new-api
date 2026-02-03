#!/usr/bin/env python3
"""
上传文件到阿里云 OSS

使用方法：
    python3 upload_to_oss.py \
        --access-key-id YOUR_AK \
        --access-key-secret YOUR_SK \
        --bucket BUCKET_NAME \
        --region REGION \
        --source-dir ./dist
"""

import argparse
import os
import sys

try:
    import oss2
except ImportError:
    print("错误: 请先安装 oss2: pip3 install oss2")
    sys.exit(1)


# Content-Type 映射
CONTENT_TYPES = {
    '.html': 'text/html; charset=utf-8',
    '.css': 'text/css; charset=utf-8',
    '.js': 'application/javascript; charset=utf-8',
    '.json': 'application/json; charset=utf-8',
    '.png': 'image/png',
    '.jpg': 'image/jpeg',
    '.jpeg': 'image/jpeg',
    '.gif': 'image/gif',
    '.svg': 'image/svg+xml',
    '.ico': 'image/x-icon',
    '.webp': 'image/webp',
    '.woff': 'font/woff',
    '.woff2': 'font/woff2',
    '.ttf': 'font/ttf',
    '.eot': 'application/vnd.ms-fontobject',
    '.map': 'application/json',
    '.txt': 'text/plain; charset=utf-8',
    '.xml': 'application/xml',
}


def get_content_type(filename):
    """根据文件扩展名获取 Content-Type"""
    ext = os.path.splitext(filename)[1].lower()
    return CONTENT_TYPES.get(ext, 'application/octet-stream')


def upload_directory(bucket, source_dir, prefix=''):
    """
    上传目录到 OSS
    
    Args:
        bucket: OSS bucket 对象
        source_dir: 本地源目录
        prefix: OSS 路径前缀
    
    Returns:
        上传的文件数量
    """
    uploaded = 0
    failed = 0
    
    for root, dirs, files in os.walk(source_dir):
        for filename in files:
            local_path = os.path.join(root, filename)
            
            # 计算相对路径作为 OSS key
            relative_path = os.path.relpath(local_path, source_dir)
            oss_key = os.path.join(prefix, relative_path) if prefix else relative_path
            
            # 统一使用正斜杠
            oss_key = oss_key.replace('\\', '/')
            
            # 获取 Content-Type
            content_type = get_content_type(filename)
            
            try:
                # 设置缓存控制
                headers = {
                    'Content-Type': content_type,
                }
                
                # 对于 HTML 文件，不缓存
                if filename.endswith('.html'):
                    headers['Cache-Control'] = 'no-cache, no-store, must-revalidate'
                # 对于带 hash 的静态资源，长期缓存
                elif '/assets/' in oss_key:
                    headers['Cache-Control'] = 'public, max-age=31536000, immutable'
                
                # 上传文件
                bucket.put_object_from_file(oss_key, local_path, headers=headers)
                uploaded += 1
                print(f"✓ {oss_key}")
                
            except Exception as e:
                failed += 1
                print(f"✗ {oss_key}: {e}")
    
    return uploaded, failed


def main():
    parser = argparse.ArgumentParser(description='上传文件到阿里云 OSS')
    parser.add_argument('--access-key-id', required=True, help='阿里云 Access Key ID')
    parser.add_argument('--access-key-secret', required=True, help='阿里云 Access Key Secret')
    parser.add_argument('--bucket', required=True, help='OSS Bucket 名称')
    parser.add_argument('--region', default='ap-southeast-1', help='OSS 区域（默认: ap-southeast-1）')
    parser.add_argument('--source-dir', required=True, help='本地源目录')
    parser.add_argument('--prefix', default='', help='OSS 路径前缀')
    
    args = parser.parse_args()
    
    # 验证源目录
    if not os.path.isdir(args.source_dir):
        print(f"错误: 源目录不存在: {args.source_dir}")
        sys.exit(1)
    
    # 构建 endpoint
    endpoint = f'https://oss-{args.region}.aliyuncs.com'
    
    print(f"OSS Endpoint: {endpoint}")
    print(f"Bucket: {args.bucket}")
    print(f"Source: {args.source_dir}")
    print(f"Prefix: {args.prefix or '(none)'}")
    print("")
    
    # 创建 OSS 客户端
    auth = oss2.Auth(args.access_key_id, args.access_key_secret)
    bucket = oss2.Bucket(auth, endpoint, args.bucket)
    
    # 上传文件
    print("正在上传文件...")
    print("-" * 40)
    
    uploaded, failed = upload_directory(bucket, args.source_dir, args.prefix)
    
    print("-" * 40)
    print(f"上传完成: {uploaded} 成功, {failed} 失败")
    
    if failed > 0:
        sys.exit(1)


if __name__ == '__main__':
    main()
