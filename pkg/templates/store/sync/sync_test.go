package sync

import "testing"

func Test_normalizeFilePattern(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				content: "./abc/app_v1.0.0.zip",
			},
			want: "abc/app_v1.0.0.zip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeFilePattern(tt.args.content); got != tt.want {
				t.Errorf("normalizeFilePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
