# ===========================================
# Application Settings
# ===========================================
APP_NAME=DefaultApp
APP_PORT=8080
APP_ENV=development                    # development | staging | production

# ===========================================
# Log Settings
# ===========================================
LOG_LEVEL=info                         # debug | info | warn | error
LOG_TIME_FORMAT=2006-01-02 15:04:05
LOG_CALLER=false                       # 是否显示调用位置
LOG_CONSOLE=true                       # 是否输出到控制台
LOG_FILE=false                         # 是否输出到文件
LOG_FILE_PATH=logs/logs.log
LOG_PRETTY=false                       # 是否美化输出

# ===========================================
# Database Settings
# ===========================================
DB_TIMEZONE=Asia/Shanghai              # 数据库时区
DB_LOG_LEVEL=info                      # debug | info | warn | error
DB_MAX_IDLE_CONNS=50                   # 最大空闲连接数
DB_MAX_OPEN_CONNS=200                  # 最大打开连接数
DB_CONN_MAX_LIFETIME=30                # 连接最大生命周期(分钟)
DB_CONN_MAX_IDLE_TIME=5                # 空闲连接最大生命周期(分钟)

# 单数据库模式 (DB_ENABLE_SLAVE=false 时使用)
DB_ENABLE_SLAVE=false
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=example

# 主从分离模式 (DB_ENABLE_SLAVE=true 时使用)
# DB_ENABLE_SLAVE=true
# DB_NAME=example
# DB_MAIN_HOST=localhost
# DB_MAIN_PORT=3306
# DB_MAIN_USER=root
# DB_MAIN_PASSWORD=
# DB_SLAVE_HOST=localhost,localhost2    # 支持多个从库，逗号分隔
# DB_SLAVE_PORT=3306
# DB_SLAVE_USER=root
# DB_SLAVE_PASSWORD=

# ===========================================
# Queue Settings (Redis-based)
# ===========================================
QUEUE_ENABLE=false
QUEUE_REDIS_HOST=localhost
QUEUE_REDIS_PORT=6379
QUEUE_REDIS_PASSWORD=
QUEUE_REDIS_DB=0

# ===========================================
# Messaging Settings (NATS/Redis)
# ===========================================
MESSAGING_DRIVER=nats                  # nats | redis
MESSAGING_URL=localhost:4222
HOSTNAME=default-instance              # 实例标识，用于消息订阅

# ===========================================
# JWT & Middleware Settings
# ===========================================
JWT_SECRET=IWY@*3JUI#d309HhefzX2WpLtPKtD!hn
MID_HEARTBEAT_PATH=/health             # 存活检查路径 (liveness)
MID_READINESS_PATH=/ready              # 就绪检查路径 (readiness)
MID_COMPRESS_LEVEL=5                   # gzip 压缩级别 (0-9)
MID_TIMEOUT=30                         # 请求超时时间(秒)
MID_THROTTLE=100                       # 请求限流数 (每分钟每IP)

# ===========================================
# CORS Settings
# ===========================================
# CORS_ALLOW_ORIGINS=*                 # 允许的域名，逗号分隔，默认 *
# CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
# CORS_ALLOW_HEADERS=Origin,Content-Type,Authorization

# ===========================================
# Storage Settings
# ===========================================
STORAGE_DRIVER=s3                      # s3 | cos
USE_CDN=false                          # 是否使用 CDN

# S3 Storage (MinIO/AWS S3)
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_USE_SSL=false
S3_DEFAULT_BUCKET=default
S3_CDN=                                # CDN 地址

# COS Storage (腾讯云对象存储)
# COS_SECRET_ID=
# COS_SECRET_KEY=
# COS_REGION=ap-beijing
# COS_DEFAULT_BUCKET=default
# COS_CDN=                             # CDN 地址
