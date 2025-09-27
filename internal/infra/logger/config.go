package logger

type Config struct {
	Level            string   `mapstructure:"level"`
	OutputPaths      []string `mapstructure:"out_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
}
