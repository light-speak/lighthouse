type AppConfig struct {
	Name       string
	Port       string
	RPCPort    string
	Env        Env
	QueueRedis *QueueRedisConfig
}

type Env string

const (
	EnvDevelopment Env = "development"
	EnvStaging     Env = "staging"
	EnvProduction  Env = "production"
)

type QueueRedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	DB       int
}


var Config *AppConfig

func init() {
	Config = &AppConfig{
		Name: 	  "DefaultApp",
		Port:     "8080",
		RPCPort:  "8081",
		Env:  EnvDevelopment,
		QueueRedis: &QueueRedisConfig{
			Enabled:  false,
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
	}

	if cp, err := os.Getwd(); err == nil {
		_ = godotenv.Load(filepath.Join(cp, ".env"))
	}

	Config.Name = utils.GetEnv("APP_NAME", Config.Name)
	Config.Port = utils.GetEnv("APP_PORT", Config.Port)
	Config.RPCPort = utils.GetEnv("APP_RPC_PORT", Config.RPCPort)
	Config.Env = Env(utils.GetEnv("APP_ENV", string(Config.Env)))
	Config.QueueRedis.Enabled = utils.GetEnvBool("QUEUE_REDIS_ENABLED", Config.QueueRedis.Enabled)
	if Config.QueueRedis.Enabled {
		Config.QueueRedis.Host = utils.GetEnv("QUEUE_REDIS_HOST", Config.QueueRedis.Host)
		Config.QueueRedis.Port = utils.GetEnv("QUEUE_REDIS_PORT", Config.QueueRedis.Port)
		Config.QueueRedis.Password = utils.GetEnv("QUEUE_REDIS_PASSWORD", Config.QueueRedis.Password)
		Config.QueueRedis.DB = utils.GetEnvInt("QUEUE_REDIS_DB", Config.QueueRedis.DB)
	}
}
