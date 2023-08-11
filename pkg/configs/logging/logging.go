package logging

import (
	"flag"

	"github.com/go-kit/log/level"
)

type Config struct {
	Level    string      `yaml:"level"`
	Format   string      `yaml:"format"`
	LogLevel level.Value `yaml:"-"`
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&c.Level, "log.level", "info", "Level of the application log.")
	c.LogLevel = level.ParseDefault(c.Level, level.InfoValue())
	f.StringVar(&c.Format, "log.format", "", "Format of the application log, default to logfmt, supported values are ['json', 'logfmt'].")
}
