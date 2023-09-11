package sync

import (
	"flag"
	"regexp"
	"strings"
	"time"
)

const (
	// templateFilePattern should be "<type>/<name>_<version>.<ext>
	TemplateFilePattern = "(\\./)?([a-zA-Z]+)/([a-zA-Z0-9-]+)_v([0-9\\.]+)(-(alpha|beta))?\\.(zip|tgz|tar\\.gz)"
)

// Config holds configuration for accessing long-term storage.
type Config struct {
	Enabled   bool          `yaml:"enabled"`
	Address   string        `yaml:"address"`
	IndexFile string        `yaml:"index"`
	Interval  time.Duration `yaml:"interval"`
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.BoolVar(&c.Enabled, "templates.store.sync.enabled", false, "Weather syncing templates from remote or not.")
	f.DurationVar(&c.Interval, "templates.store.sync.interval", 10*time.Minute, "Interval of syncing templates from remote")
	f.StringVar(&c.Address, "templates.store.sync.address", "", "Remote address of the templates.")
	f.StringVar(&c.IndexFile, "templates.store.sync.index", "index.list", "List of templates at the remote address. Should be each oneline")
}

func normalizeFilePattern(content string) string {
	reg := regexp.MustCompile(TemplateFilePattern)
	if reg.MatchString(content) {
		return strings.TrimLeft(content, "./")
	}
	return ""
}
