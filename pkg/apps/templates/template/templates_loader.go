package template

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/udmire/observability-operator/pkg/utils"
)

const (
	EXT_ZIP    = ".zip"
	EXT_TGZ    = ".tgz"
	EXT_TAR_GZ = ".tar.gz"
)

type TemplateLoader interface {
	LoadTemplate(path string) (*AppTemplate, error)
	TemplateName(path string) (string, string)
}

type templatesLoader struct {
	logger log.Logger
}

func NewTemplateLoader(logger log.Logger) TemplateLoader {
	tl := &templatesLoader{logger: logger}
	return tl
}

func (l *templatesLoader) LoadTemplate(path string) (app *AppTemplate, err error) {
	appVer, ext := l.TemplateName(path)
	switch ext {
	case EXT_ZIP:
		app, err = l.handleZipFile(path, appVer)
	case EXT_TAR_GZ, EXT_TGZ:
		app, err = l.handleTarGzFile(path, appVer)
	default:
		return nil, nil
	}

	if err != nil {
		level.Warn(l.logger).Log("msg", "load to template failed", "path", path, "err", err)
	}

	return app, err
}

func (l *templatesLoader) handleZipFile(path, appVer string) (*AppTemplate, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		level.Warn(l.logger).Log("msg", "invalid zip template file", "err", err)
		return nil, err
	}
	defer zipFile.Close()

	tempDir, err := os.MkdirTemp("", "template-")
	if err != nil {
		level.Warn(l.logger).Log("msg", "cannot create temp space for unzip template", "path", path, "err", err)
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	for _, file := range zipFile.File {
		destPath := filepath.Join(tempDir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(destPath, file.Mode())
			continue
		}
		fileReader, err := file.Open()
		if err != nil {
			level.Warn(l.logger).Log("msg", "cannot open zipfile reader", "file", file.Name, "err", err)
			return nil, err
		}
		defer fileReader.Close()

		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			level.Warn(l.logger).Log("msg", "cannot create temp file for content", "err", err)
			return nil, err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, fileReader)
		if err != nil {
			level.Warn(l.logger).Log("msg", "cannot copy content to temp file", "err", err)
			return nil, err
		}
	}

	return l.loadTemplateWithFolder(appVer, tempDir)
}

func (l *templatesLoader) handleTarGzFile(path, appVer string) (*AppTemplate, error) {
	file, err := os.Open(path)
	if err != nil {
		level.Warn(l.logger).Log("msg", "cannot open template file", "path", path, "err", err)
		return nil, err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		level.Warn(l.logger).Log("msg", "invalid tar.gz template file", "path", path, "err", err)
		return nil, err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	tempDir, err := os.MkdirTemp("", "template-")
	if err != nil {
		level.Warn(l.logger).Log("msg", "cannot create temp space for tar.gz template", "path", path, "err", err)
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			level.Warn(l.logger).Log("msg", "create tag.gz template content failed", "path", path, "err", err)
			return nil, err
		}

		target := filepath.Join(tempDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(target, os.FileMode(header.Mode))
			if err != nil {
				level.Warn(l.logger).Log("msg", "cannot create directory in temp space", "target", target, "err", err)
				return nil, err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				level.Warn(l.logger).Log("msg", "cannot create file in temp space ", "target", target, "err", err)
				return nil, err
			}
			defer file.Close()

			_, err = io.Copy(file, tarReader)
			if err != nil {
				level.Warn(l.logger).Log("msg", "cannot copy content in temp space", "target", target, "err", err)
				return nil, err
			}
		}
	}

	return l.loadTemplateWithFolder(appVer, tempDir)
}

func (l *templatesLoader) loadTemplateWithFolder(appVer, tempDir string) (app *AppTemplate, err error) {
	appVerArr := strings.Split(appVer, "_")
	if len(appVerArr) != 2 {
		level.Warn(l.logger).Log("msg", "the template filename not match pattern 'app_version'", "template", appVer)
		return nil, fmt.Errorf("invalid app package name %s", appVer)
	}

	appFS := os.DirFS(tempDir)

	rootPath := "."
	var exists bool
	if exists, err = utils.HasOnlySubDirectory(tempDir, appVer); exists {
		rootPath = appVer
	}
	if err != nil {
		level.Warn(l.logger).Log("msg", "inliad application package", "err", err)
		return nil, err
	}
	if exists, err = utils.HasOnlySubDirectory(tempDir, appVerArr[0]); exists {
		rootPath = appVerArr[0]
	}
	if err != nil {
		level.Warn(l.logger).Log("msg", "inliad application package", "err", err)
		return nil, err
	}

	app = &AppTemplate{
		TemplateBase: TemplateBase{Name: appVerArr[0], Version: appVerArr[1]},
		Workloads:    map[string]*WorkloadTemplate{},
	}
	err = fs.WalkDir(appFS, rootPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			level.Warn(l.logger).Log("msg", "failed to walk content of template", "template", appVer, "err", err)
			return err
		}

		if entry.IsDir() {
			if path == appVer {
				return nil
			}

			work := WorkloadTemplate{
				TemplateBase: TemplateBase{Name: entry.Name(), Version: appVerArr[1]},
			}
			app.Workloads[entry.Name()] = &work
		} else {
			content, err := fs.ReadFile(appFS, path)
			if err != nil {
				level.Warn(l.logger).Log("msg", "cannot read template file", "path", path, "err", err)
				return err
			}
			if len(content) == 0 {
				return nil
			}

			base := filepath.Dir(path)
			var templateBase *TemplateBase
			if base == appVer {
				templateBase = &app.TemplateBase
			} else {
				templateBase = &app.Workloads[base].TemplateBase
			}

			templateBase.TemplateFiles = append(templateBase.TemplateFiles, &TemplateFile{
				FileName: entry.Name(),
				Content:  content,
			})
		}
		return nil
	})

	if err == nil {
		return app, nil
	}

	return nil, err
}

func (l *templatesLoader) TemplateName(path string) (string, string) {
	_, file := filepath.Split(path)
	lower_file := strings.ToLower(file)
	if strings.HasSuffix(lower_file, EXT_ZIP) {
		return file[:len(file)-len(EXT_ZIP)], EXT_ZIP
	} else if strings.HasSuffix(file, EXT_TGZ) {
		return file[:len(file)-len(EXT_TGZ)], EXT_TGZ
	} else if strings.HasSuffix(lower_file, EXT_TAR_GZ) {
		return file[:len(file)-len(EXT_TAR_GZ)], EXT_TAR_GZ
	}
	return "", ""
}

// func (l *templatesLoader) startWatching() {
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		l.log.Error(err, "cannot start the templates loader watcher")
// 		return
// 	}
// 	defer watcher.Close()

// 	go func() {
// 		for {
// 			select {
// 			case event, ok := <-l.Watcher.Events:
// 				if !ok {
// 					return
// 				}
// 				path, err := filepath.Abs(event.Name)
// 				if err != nil {
// 					continue
// 				}

// 				fileinfo, err := os.Stat(path)
// 				if err != nil {
// 					continue
// 				}

// 				if fileinfo.IsDir() {
// 					continue
// 				}
// 				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
// 					err = l.loadFile(path)
// 					if err != nil {
// 						continue
// 					}
// 					l.log.Info("load template succeed", "path", path)
// 				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
// 					err = l.unloadFile(path)
// 					if err != nil {
// 						continue
// 					}
// 					l.log.Info("unload template succeed", "path", path)
// 				}
// 			case err, ok := <-l.Watcher.Errors:
// 				if !ok {
// 					return
// 				}
// 				l.log.Error(err, "watch templates error found.")
// 			}
// 		}
// 	}()
// 	l.Watcher.Add(l.Dir)

// 	sigchan := make(chan os.Signal, 1)
// 	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
// 	<-sigchan
// }
