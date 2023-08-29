package local

import (
	"context"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/log/level"
)

func (l *LocalStore) startup(ctx context.Context) error {
	return l.Load()
}

func (l *LocalStore) watching(ctx context.Context) error {
	l.watcher, _ = fsnotify.NewWatcher()
	l.watcher.Add(l.cfg.Directory)

	for {
		select {
		case event := <-l.watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				file, err := os.Stat(event.Name)
				if err != nil {
					level.Warn(l.logger).Log("msg", "failed to read file stat", "file", event.Name)
					continue
				}

				if file.IsDir() {
					level.Debug(l.logger).Log("msg", "watched directory change", "dir", event.Name)
					continue
				}
				l.LoadTemplate(event.Name)
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove {
				l.UnloadTemplate(event.Name)
			}

			if event.Op&fsnotify.Rename == fsnotify.Rename {
				l.UnloadTemplate(event.Name)
			}
		case err := <-l.watcher.Errors:
			fmt.Println("error : ", err)
			return err
		}
	}
}

func (l *LocalStore) shutdown(_ error) error {
	return l.watcher.Close()
}
