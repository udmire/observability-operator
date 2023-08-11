package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHasOnlySubDirectory(t *testing.T) {
	type args struct {
		name string
		dirs []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "has only directory",
			args: args{
				name: "dir1",
				dirs: []string{"dir1"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "has only directory not match",
			args: args{
				name: "dir1",
				dirs: []string{"dir"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "has more than one directories",
			args: args{
				name: "dir1",
				dirs: []string{"dir1", "dir2"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "has three directories",
			args: args{
				name: "dir1",
				dirs: []string{"dir", "dir1", "dir2"},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, _ := os.MkdirTemp("", "xxx-")
			defer os.RemoveAll(tempDir)
			for _, name := range tt.args.dirs {
				os.Mkdir(filepath.Join(tempDir, name), os.ModeDir)
			}
			got, err := HasOnlySubDirectory(tempDir, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("HasOnlySubDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HasOnlySubDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}
