package internalgrpc

type Config struct {
	Type string `mapstructure:"type,omitempty"`
	Host string `mapstructure:"host,omitempty"`
	Port string `mapstructure:"port,omitempty"`
}
