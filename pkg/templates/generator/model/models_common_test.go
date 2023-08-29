package model

import (
	"testing"
)

func TestConfigMap_String(t *testing.T) {
	type fields struct {
		GenericModel GenericModel
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "cm",
			fields: fields{
				GenericModel: GenericModel{
					Base: Base{
						Name:   "cm1",
						Labels: make(map[string]string),
					},
					DefaultNamespace: "default",
				},
			},
			want: "configmap: default/cm1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigMap{
				GenericModel: tt.fields.GenericModel,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("ConfigMap.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}
