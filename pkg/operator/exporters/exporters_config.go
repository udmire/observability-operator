package exporters

import "flag"

type Config struct {
	Concurrency int
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.IntVar(&c.Concurrency, "exporters.concurrency", 3, "Max concurrent deploying exporters.")
}
