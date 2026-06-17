package mqxbinding

type MQConfig struct {
	Driver string     `yaml:"driver"`
	Mode   string     `yaml:"mode"`
	Addrs  []string   `yaml:"addrs"`
	Auth   AuthConfig `yaml:"auth,omitempty"`
	Redis  RedisConfig `yaml:"redis,omitempty"`
	Kafka  KafkaConfig `yaml:"kafka,omitempty"`
}

type AuthConfig struct {
	Password string `yaml:"password"`
	Username string `yaml:"username,omitempty"`
}

type RedisConfig struct {
	DB       int `yaml:"db,omitempty"`
	PoolSize int `yaml:"pool_size,omitempty"`
	MaxLen   int `yaml:"max_len,omitempty"`
}

type KafkaConfig struct {
	SecurityProtocol string `yaml:"security_protocol,omitempty"`
	SASLMechanism    string `yaml:"sasl_mechanism,omitempty"`
}
