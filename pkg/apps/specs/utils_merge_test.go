package specs

import (
	"reflect"
	"testing"

	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestMergePatchContainers(t *testing.T) {
	type args struct {
		base    []core_v1.Container
		patches []core_v1.Container
	}
	tests := []struct {
		name    string
		args    args
		want    []core_v1.Container
		wantErr bool
	}{
		{
			name: "",
			args: args{
				base: []core_v1.Container{{
					Name: "abc",
					Args: []string{"command"},
					Resources: core_v1.ResourceRequirements{
						Limits: core_v1.ResourceList{
							core_v1.ResourceCPU: resource.MustParse("0"),
						},
					},
				}, {
					Name: "def",
				}},
				patches: []core_v1.Container{{
					Name: "abc",
					Args: []string{"--args"},
					Resources: core_v1.ResourceRequirements{
						Requests: core_v1.ResourceList{
							core_v1.ResourceCPU: resource.MustParse("0"),
						},
					},
				}},
			},
			want: []core_v1.Container{{
				Name: "abc",
				Args: []string{"--args"},
				Resources: core_v1.ResourceRequirements{
					Requests: core_v1.ResourceList{
						core_v1.ResourceCPU: resource.MustParse("0"),
					},
					Limits: core_v1.ResourceList{
						core_v1.ResourceCPU: resource.MustParse("0"),
					},
				},
			}, {
				Name: "def",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergePatchContainers(tt.args.base, tt.args.patches)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergePatchContainers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergePatchContainers() = %v, want %v", got, tt.want)
			}
		})
	}
}
