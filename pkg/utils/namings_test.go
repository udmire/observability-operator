package utils

import "testing"

func TestUpdateImageRegistry(t *testing.T) {
	type args struct {
		registry string
		image    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "a",
			args: args{
				registry: "registry.udmire.cn",
				image:    "quay.io/udmire/alpine:3.12",
			},
			want: "registry.udmire.cn/udmire/alpine:3.12",
		},
		{
			name: "b",
			args: args{
				registry: "registry.udmire.cn",
				image:    "udmire/alpine:3.12",
			},
			want: "registry.udmire.cn/udmire/alpine:3.12",
		},
		{
			name: "c",
			args: args{
				registry: "registry.udmire.cn",
				image:    "192.168.0.1/udmire/alpine:3.12",
			},
			want: "registry.udmire.cn/udmire/alpine:3.12",
		},
		{
			name: "d",
			args: args{
				registry: "registry.udmire.cn",
				image:    "192.168.0.1:1234/udmire/alpine:3.12",
			},
			want: "registry.udmire.cn/udmire/alpine:3.12",
		},
		{
			name: "e",
			args: args{
				registry: "registry.udmire.cn",
				image:    "[::]:1234/udmire/alpine:3.12",
			},
			want: "registry.udmire.cn/udmire/alpine:3.12",
		},
		{
			name: "f",
			args: args{
				registry: "registry.udmire.cn",
				image:    "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:1234/udmire/alpine:3.12",
			},
			want: "registry.udmire.cn/udmire/alpine:3.12",
		},
		{
			name: "g",
			args: args{
				registry: "registry.udmire.cn",
				image:    "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]/udmire/alpine:3.12",
			},
			want: "registry.udmire.cn/udmire/alpine:3.12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpdateImageRegistry(tt.args.registry, tt.args.image); got != tt.want {
				t.Errorf("UpdateImageRegistry() = %v, want %v", got, tt.want)
			}
		})
	}
}
