package manifest

import (
	"reflect"
	"testing"

	"github.com/udmire/observability-operator/pkg/templates/template"
)

func Test_templateBuilder_buildCapsules(t *testing.T) {
	type fields struct {
		template *template.AppTemplate
	}
	type args struct {
		content []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name:   "",
			fields: fields{},
			args: args{
				content: []byte(`- name: abc
  type: secrets
  items:
    abc: def
    xxx: yyy
- name: abcd
  type: configmaps
  dynamics: .*`),
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &templateBuilder{
				template: tt.fields.template,
			}
			if got := b.buildCapsules(tt.args.content); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("templateBuilder.buildCapsules() = %v, want %v", got, tt.want)
			}
		})
	}
}
