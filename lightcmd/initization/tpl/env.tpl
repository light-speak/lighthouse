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
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=example

# Elasticsearch settings
ELASTICSEARCH_ENABLE=false
ELASTICSEARCH_HOST=localhost
ELASTICSEARCH_PORT=9200
ELASTICSEARCH_USER=
ELASTICSEARCH_PASSWORD=

# Queue settings
QUEUE_ENABLE=false
QUEUE_REDIS_HOST=localhost
QUEUE_REDIS_PORT=6379
QUEUE_REDIS_PASSWORD=
QUEUE_REDIS_DB=0

# Kafka settings
LIGHT_KAFKA_ENABLED=false
LIGHT_KAFKA_BROKERS=localhost:9092
LIGHT_KAFKA_GROUP_ID=default