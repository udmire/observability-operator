package sync

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/grafana/dskit/services"
)

type httpSync struct {
	*services.BasicService

	cfg    Config
	logger log.Logger

	storePath string
	templates map[string]string // template urn & path
	mutex     sync.Mutex
}

func NewHttpSynchronizer(cfg Config, storePath string, logger log.Logger) *httpSync {
	sync := &httpSync{
		cfg:       cfg,
		logger:    logger,
		templates: make(map[string]string),
		storePath: storePath,
		mutex:     sync.Mutex{},
	}
	sync.BasicService = services.NewTimerService(cfg.Interval, sync.Synchronize, sync.Synchronize, nil)
	return sync
}

func (s *httpSync) Synchronize(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	idx, err := s.getIndex()
	if err != nil {
		level.Warn(s.logger).Log("msg", "failed to download templates index", "err", err)
		return err
	}

	toAdd, toDel := compareMaps(s.templates, idx)
	for name, url := range toAdd {
		err := s.DownloadTemplate(name, url)
		if err != nil {
			level.Warn(s.logger).Log("msg", "failed to download template", "template", name, "url", url, "err", err)
			continue
		}
		s.templates[name] = url
	}

	for name := range toDel {
		s.RemoveTemplate(name)
		delete(s.templates, name)
	}
	return nil
}

func (s *httpSync) DownloadTemplate(name, url string) error {
	splits := strings.SplitN(name, string(filepath.Separator), 2)
	if len(splits) != 2 {
		return fmt.Errorf("invalid name for templates")
	}

	resp, err := http.Get(url)
	if err != nil {
		level.Warn(s.logger).Log("msg", "cannot download template.", "url", url, "err", err)
		return err
	}
	defer resp.Body.Close()

	tempFile, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("%s.temp", splits[1]))
	if err != nil {
		level.Warn(s.logger).Log("msg", "cannot create temp file for downloading template.", "url", url, "err", err)
		return err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		level.Warn(s.logger).Log("msg", "cannot download template content.", "url", url, "err", err)
		return err
	}
	resp.Body.Close()

	err = os.MkdirAll(filepath.Join(s.storePath, splits[0]), 0755)
	if err != nil {
		level.Warn(s.logger).Log("msg", "cannot create target dir for content.", "dir", splits[0], "err", err)
		return err
	}

	err = os.Rename(tempFile.Name(), filepath.Join(s.storePath, splits[0], splits[1]))
	if err != nil {
		level.Warn(s.logger).Log("msg", "cannot rename to target file.", "name", name, "err", err)
		return err
	}

	return nil
}

func (s *httpSync) RemoveTemplate(name string) {
	splits := strings.SplitN(name, string(filepath.Separator), 2)
	if len(splits) != 2 {
		level.Warn(s.logger).Log("msg", "cannot remove template.", "name", name)
		return
	}
	err := os.Remove(filepath.Join(s.storePath, splits[0], splits[1]))
	if err != nil {
		level.Warn(s.logger).Log("msg", "cannot remove template.", "name", name, "err", err)
	}
}

func (s *httpSync) getIndex() (result map[string]string, err error) {
	client := http.Client{}
	idxUri, err := url.JoinPath(s.cfg.Address, s.cfg.IndexFile)
	if err != nil {
		level.Warn(s.logger).Log("msg", "invalid address for templates synchorize.", "err", err)
		return nil, err
	}
	resp, err := client.Get(idxUri)
	if err != nil {
		level.Warn(s.logger).Log("msg", "cannot download templates index.", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	result = map[string]string{}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		normalized := normalizeFilePattern(line)
		if normalized == "" {
			level.Error(s.logger).Log("msg", "invalid content of templates index.", "err", err)
			return nil, fmt.Errorf("invalid content of templates index file, provide: '%s', need match: '%s'", line, TemplateFilePattern)
		}
		remote, _ := url.JoinPath(s.cfg.Address, normalized)
		result[normalized] = remote
	}
	return result, nil
}

func compareMaps(ori, oth map[string]string) (toAdd, toDel map[string]string) {
	if len(ori) == 0 {
		return oth, toDel
	}

	if len(oth) == 0 {
		return toAdd, ori
	}

	toDel = map[string]string{}
	toAdd = map[string]string{}

	for k, v := range oth {
		nv, exists := ori[k]
		if !exists {
			toAdd[k] = v
		} else if nv != v {
			toAdd[k] = v
		}

	}

	for k, v := range ori {
		if _, exists := oth[k]; !exists {
			toDel[k] = v
		}
	}

	return toAdd, toDel
}
