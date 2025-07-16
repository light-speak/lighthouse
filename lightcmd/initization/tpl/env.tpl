# App
APP_NAME=DefaultApp
APP_PORT=8080
APP_ENV=development
APP_THROTTLE=100

# Log
LOG_LEVEL=info
LOG_TIME_FORMAT=2006-01-02 15:04:05
LOG_CALLER=false
LOG_CONSOLE=true
LOG_FILE=false
LOG_FILE_PATH=logs/logs.log

# Database settings
DB_ENABLE_SLAVE=false
DB_NAME=example
DB_LOG_LEVEL=info

# Main database settings
DB_MAIN_HOST=localhost
DB_MAIN_PORT=3306
DB_MAIN_USER=root
DB_MAIN_PASSWORD=

# Slave database settings (从数据库，当DB_ENABLE_SLAVE=true时使用)
DB_SLAVE_HOST=localhost
DB_SLAVE_PORT=3306
DB_SLAVE_USER=root
DB_SLAVE_PASSWORD=

# Legacy single database settings (兼容旧版本)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=

# Queue settings
QUEUE_ENABLE=false
QUEUE_REDIS_HOST=localhost
QUEUE_REDIS_PORT=6379
QUEUE_REDIS_PASSWORD=
QUEUE_REDIS_DB=0

# Redis settings
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Messaging settings
MESSAGING_DRIVER=nats
MESSAGING_URL=localhost:4222
HOSTNAME=default-instance

# JWT and Middleware settings
JWT_SECRET=IWY@*3JUI#d309HhefzX2WpLtPKtD!hn
MID_HEARTBEAT_PATH=/health
MID_COMPRESS_LEVEL=5
MID_TIMEOUT=30
MID_THROTTLE=100

# Storage settings
STORAGE_DRIVER=s3

# S3 storage settings
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_USE_SSL=false
S3_DEFAULT_BUCKET=default

# COS storage settings (腾讯云对象存储)
COS_SECRET_ID=
COS_SECRET_KEY=
COS_REGION=ap-beijing
COS_DEFAULT_BUCKET=default