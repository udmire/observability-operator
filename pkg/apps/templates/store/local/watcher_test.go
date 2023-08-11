package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-kit/log"
)

func TestLocalStore_watching(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curDir, _ := os.Getwd()
			config := Config{Directory: filepath.Join(curDir, "..", "..", "template")}
			l := New(config, log.NewNopLogger())
			_ = l.StartAsync(context.Background())
			_ = l.AwaitRunning(context.Background())

			copyFile(filepath.Join(config.Directory, "app_v1.0.1.tar.gz"), filepath.Join(config.Directory, "app_v1.0.3.tar.gz"))

			time.Sleep(time.Second * 3)

			os.Rename(filepath.Join(config.Directory, "app_v1.0.3.tar.gz"), filepath.Join(config.Directory, "app_v1.0.4.tar.gz"))

			time.Sleep(time.Second * 3)

			os.Remove(filepath.Join(config.Directory, "app_v1.0.4.tar.gz"))

			time.Sleep(time.Second * 3)
			l.StopAsync()
		})
	}
}

func copyFile(from, to string) {
	source, _ := os.Open(from)
	defer source.Close()

	destination, _ := os.Create(to)
	defer destination.Close()
	_, _ = io.Copy(destination, source)
}
