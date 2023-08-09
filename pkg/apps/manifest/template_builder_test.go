package manifest

import (
	"reflect"
	"testing"

	"github.com/udmire/observability-operator/pkg/apps/templates/template"
)

func Test_recognize(t *testing.T) {
	type args struct {
		file *template.TemplateFile
	}
	tests := []struct {
		name  string
		args  args
		want  ManifestType
		want1 string
	}{
		{
			name: "recognize_configmap",
			args: args{
				file: &template.TemplateFile{
					FileName: "app_configmap.yaml",
				},
			},
			want:  ConfigMap,
			want1: "app",
		},
		{
			name: "recognize_configmap_1",
			args: args{
				file: &template.TemplateFile{
					FileName: "app-configmap.yaml",
				},
			},
			want:  ConfigMap,
			want1: "app",
		},
		{
			name: "recognize_configmap_2",
			args: args{
				file: &template.TemplateFile{
					FileName: "app_configmap.yml",
				},
			},
			want:  ConfigMap,
			want1: "app",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := recognize(tt.args.file)
			if got != tt.want {
				t.Errorf("recognize() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("recognize() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_templateBuilder_Build(t *testing.T) {
	type fields struct {
		template *template.AppTemplate
	}

	template := &template.AppTemplate{TemplateBase: template.TemplateBase{
		Name:    "app",
		Version: "version",
	}}

	tests := []struct {
		name   string
		fields fields
		want   *AppManifests
	}{
		{
			name: "",
			fields: fields{
				template: template,
			},
			want: &AppManifests{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &templateBuilder{
				template: tt.fields.template,
			}
			if got := b.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("templateBuilder.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
