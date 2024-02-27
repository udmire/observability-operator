package specs

import (
	"testing"

	"github.com/go-kit/log"
	"github.com/udmire/observability-operator/pkg/apps/manifest"
	"github.com/udmire/observability-operator/pkg/templates/provider"
	v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
)

func Test_appHandler_updateImagesWithRegistry(t *testing.T) {
	type fields struct {
		logger   log.Logger
		provider provider.TemplateProvider
	}
	type args struct {
		registry string
		manifest *manifest.AppManifests
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{},
			args: args{
				registry: "registry.udmire.cn",
				manifest: &manifest.AppManifests{
					CompsMenifests: []*manifest.CompManifests{
						{
							Name: "comp",
							Deployment: &v1.Deployment{
								Spec: v1.DeploymentSpec{
									Template: core_v1.PodTemplateSpec{
										Spec: core_v1.PodSpec{
											InitContainers: []core_v1.Container{
												{
													Name:  "init",
													Image: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]/udmire/alpine:3.12",
												},
											},
											Containers: []core_v1.Container{
												{
													Name:  "agent",
													Image: "query.io/udmire/alpine:3.12",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: "registry.udmire.cn/udmire/alpine:3.12",
		},
		{
			fields: fields{},
			args: args{
				registry: "",
				manifest: &manifest.AppManifests{
					CompsMenifests: []*manifest.CompManifests{
						{
							Name: "comp",
							Deployment: &v1.Deployment{
								Spec: v1.DeploymentSpec{
									Template: core_v1.PodTemplateSpec{
										Spec: core_v1.PodSpec{
											InitContainers: []core_v1.Container{
												{
													Name:  "init",
													Image: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]/udmire/alpine:3.12",
												},
											},
											Containers: []core_v1.Container{
												{
													Name:  "agent",
													Image: "query.io/udmire/alpine:3.12",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: "query.io/udmire/alpine:3.12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &appHandler{
				logger:   tt.fields.logger,
				provider: tt.fields.provider,
			}
			h.updateImagesWithRegistry(tt.args.registry, tt.args.manifest)
			if got := tt.args.manifest.CompsMenifests[0].Deployment.Spec.Template.Spec.Containers[0].Image; got != tt.want {
				t.Errorf("UpdateImageRegistry() = %v, want %v", got, tt.want)
			}
		})
	}
}
