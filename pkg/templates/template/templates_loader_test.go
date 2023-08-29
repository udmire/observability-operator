package template

import (
	"testing"

	"github.com/go-kit/log"
)

func Test_templatesLoader_handleZipFile(t *testing.T) {
	type fields struct {
		Dir string
	}
	type args struct {
		path   string
		appVer string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "zip",
			fields: fields{
				Dir: "./",
			},
			args: args{
				path:   "./app_v1.0.1.1.zip",
				appVer: "app_v1.0.1.1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &templatesLoader{
				logger: log.NewNopLogger(),
			}
			if _, err := l.handleZipFile(tt.args.path, tt.args.appVer); (err != nil) != tt.wantErr {
				t.Errorf("templatesLoader.handleZipFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_templatesLoader_handleTarGzFile(t *testing.T) {
	type fields struct {
		Dir string
	}
	type args struct {
		path   string
		appVer string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "tar.gz",
			fields: fields{
				Dir: "./",
			},
			args: args{
				path:   "./app_v1.0.2.tar.gz",
				appVer: "app_v1.0.2",
			},
			wantErr: false,
		},
		{
			name: "tgz",
			fields: fields{
				Dir: "./",
			},
			args: args{
				path:   "./app_v1.0.0.tgz",
				appVer: "app_v1.0.0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &templatesLoader{
				logger: log.NewNopLogger(),
			}
			if _, err := l.handleTarGzFile(tt.args.path, tt.args.appVer); (err != nil) != tt.wantErr {
				t.Errorf("templatesLoader.handleTarGzFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_templatesLoader_fileExt(t *testing.T) {
	type fields struct {
		logger log.Logger
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  string
	}{
		{
			name: "normal",
			fields: fields{
				logger: log.NewNopLogger(),
			},
			args: args{
				path: "/abc/app_version.tgz",
			},
			want:  "app_version",
			want1: ".tgz",
		},
		{
			name: "failed",
			fields: fields{
				logger: log.NewNopLogger(),
			},
			args: args{
				path: "/abc/app_version.abc",
			},
			want:  "",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &templatesLoader{
				logger: tt.fields.logger,
			}
			got, got1 := l.TemplateName(tt.args.path)
			if got != tt.want {
				t.Errorf("templatesLoader.fileExt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("templatesLoader.fileExt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
