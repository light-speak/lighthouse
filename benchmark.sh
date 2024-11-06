#!/bin/bash

# 设置压测参数
CONCURRENT=200  # 并发数
REQUESTS=10000  # 总请求数
URL="http://localhost:8000/query"  # GraphQL endpoint

# 检查是否安装了 ab
if ! command -v ab &> /dev/null; then
    echo "Apache Benchmark (ab) not found. Please install it first."
    echo "On Ubuntu/Debian: sudo apt-get install apache2-utils"
    echo "On MacOS: brew install apache2-utils"
    exit 1
fi

# 检查 query.json 是否存在
if [ ! -f "query.json" ]; then
    echo "query.json not found!"
    exit 1
fi

echo "Starting benchmark..."
echo "URL: $URL"
echo "Concurrent users: $CONCURRENT"
echo "Total requests: $REQUESTS"
echo "----------------------------------------"

# 执行压测
ab -n $REQUESTS \
   -c $CONCURRENT \
   -T 'application/json' \
   -p query.json \
   -H "Accept: application/json" \
   -v 4 \
   $URL

echo "----------------------------------------"
echo "Benchmark completed!" 