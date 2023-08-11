package apps

import "flag"

type Config struct {
	Concurrency int
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.IntVar(&c.Concurrency, "apps.concurrency", 3, "Max concurrent deploying apps.")
}
