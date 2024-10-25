# Application settings
APP_NAME=MyApp
APP_ENVIRONMENT=development
APP_MODE=single

# Server settings
SERVER_PORT=8080
SERVER_THROTTLE=100

# API settings
API_RESTFUL=true
API_PREFIX=/api

# Database settings
DB_HOST=localhost
DB_PORT=5432
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=mydb
DB_DRIVER=postgres
DB_ORM=gorm

# Authentication settings
AUTH_DRIVER=jwt

# Logger settings
LOGGER_LEVEL=info
LOGGER_PATH=./logs/app.log
LOGGER_STACK=false
LOGGER_DRIVER=file

# Redis settings
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Elasticsearch settings
ELASTICSEARCH_ENABLE=false
ELASTICSEARCH_HOST=localhost
ELASTICSEARCH_PORT=9200
ELASTICSEARCH_USER=elastic
ELASTICSEARCH_PASSWORD=changeme

