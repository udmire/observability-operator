package model

import (
	"reflect"
	"testing"

	"github.com/udmire/observability-operator/pkg/apps/manifest"
)

func TestApp_Build(t *testing.T) {
	type fields struct {
		GenericModel GenericModel
		Common       Common
		Version      string
		Components   []*Component
	}
	tests := []struct {
		name   string
		fields fields
		want   *manifest.AppManifests
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				GenericModel: tt.fields.GenericModel,
				Common:       tt.fields.Common,
				Version:      tt.fields.Version,
				Components:   tt.fields.Components,
			}
			if got := a.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("App.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
