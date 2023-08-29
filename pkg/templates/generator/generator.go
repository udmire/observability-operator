package generator

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/udmire/observability-operator/pkg/apps/manifest"
	apps_v1 "k8s.io/api/apps/v1"
	autoscaling_v1 "k8s.io/api/autoscaling/v1"
	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

type Generator struct {
	encoder runtime.Encoder
	logger  log.Logger
}

func NewGenerator(logger log.Logger) *Generator {
	scheme := runtime.NewScheme()
	gv := core_v1.SchemeGroupVersion
	scheme.AddKnownTypes(gv, &apps_v1.Deployment{}, &apps_v1.DaemonSet{}, &apps_v1.StatefulSet{}, &apps_v1.ReplicaSet{})
	scheme.AddKnownTypes(gv, &batch_v1.Job{}, &batch_v1.CronJob{})
	scheme.AddKnownTypes(gv, &networking_v1.Ingress{})
	scheme.AddKnownTypes(gv, &autoscaling_v1.HorizontalPodAutoscaler{})
	scheme.AddKnownTypes(gv, &rbac_v1.ClusterRole{}, &rbac_v1.ClusterRoleBinding{}, &rbac_v1.Role{}, &rbac_v1.RoleBinding{})
	scheme.AddKnownTypes(gv, &core_v1.ConfigMap{}, &core_v1.Secret{}, &core_v1.ServiceAccount{}, &core_v1.Service{})
	codecs := serializer.NewCodecFactory(scheme)
	encoder := codecs.EncoderForVersion(json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme, scheme, json.SerializerOptions{
		Yaml: true, Pretty: true, Strict: false,
	}), gv)
	return &Generator{
		encoder: encoder,
		logger:  logger,
	}
}

func (g *Generator) Generate(name, version string, manifest *manifest.AppManifests, path string) error {
	if manifest == nil {
		return fmt.Errorf("invalid manifests")
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("invalid path")
	}

	root, err := os.MkdirTemp(path, name)
	if err != nil {
		return fmt.Errorf("cannot write in given folder")
	}
	defer os.RemoveAll(root)

	err = g.writeAppManifestsToGivenFolder(name, manifest, root)
	if err != nil {
		return err
	}

	return g.writeZip(root, filepath.Join(path, fmt.Sprintf("%s_%s.zip", name, version)))
}

func (g *Generator) writeAppManifestsToGivenFolder(app string, manifests *manifest.AppManifests, root string) error {
	g.writeManifestsToGivenFolder(app, &manifests.Manifests, root)
	for _, comp := range manifests.CompsMenifests {
		compDir := filepath.Join(root, comp.Name)
		_ = os.Mkdir(filepath.Join(root, comp.Name), os.ModeDir)
		g.writeManifestsToGivenFolder(comp.Name, &comp.Manifests, compDir)

		if comp.Deployment != nil {
			data, err := runtime.Encode(g.encoder, comp.Deployment)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to encoding deployment", "name", app)
				return err
			}
			err = os.WriteFile(filepath.Join(compDir, fmt.Sprintf("%s_deployment.yml", app)), data, 0644)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to write deployment to file", "name", app)
				return err
			}
		}
		if comp.DaemonSet != nil {
			data, err := runtime.Encode(g.encoder, comp.DaemonSet)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to encoding daemonSet", "name", comp.Name)
				return err
			}
			err = os.WriteFile(filepath.Join(compDir, fmt.Sprintf("%s_daemonset.yml", comp.Name)), data, 0644)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to write daemonSet to file", "name", app)
				return err
			}
		}
		if comp.StatefulSet != nil {
			data, err := runtime.Encode(g.encoder, comp.StatefulSet)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to encoding statefulSet", "name", comp.Name)
				return err
			}
			err = os.WriteFile(filepath.Join(compDir, fmt.Sprintf("%s_statefulset.yml", comp.Name)), data, 0644)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to write statefulSet to file", "name", app)
				return err
			}
		}
		if comp.ReplicaSet != nil {
			data, err := runtime.Encode(g.encoder, comp.ReplicaSet)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to encoding replicaSet", "name", comp.Name)
				return err
			}
			err = os.WriteFile(filepath.Join(compDir, fmt.Sprintf("%s_replicaset.yml", comp.Name)), data, 0644)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to write replicaSet to file", "name", app)
				return err
			}
		}
		if comp.Job != nil {
			data, err := runtime.Encode(g.encoder, comp.Job)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to encoding job", "name", app)
				return err
			}
			err = os.WriteFile(filepath.Join(compDir, fmt.Sprintf("%s_job.yml", app)), data, 0644)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to write job to file", "name", app)
				return err
			}
		}
		if comp.CronJob != nil {
			data, err := runtime.Encode(g.encoder, comp.CronJob)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to encoding cronjob", "name", app)
				return err
			}
			err = os.WriteFile(filepath.Join(compDir, fmt.Sprintf("%s_cronjob.yml", app)), data, 0644)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to write cronjob to file", "name", app)
				return err
			}
		}
		if comp.HPA != nil {
			data, err := runtime.Encode(g.encoder, comp.HPA)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to encoding HPA", "name", app)
				return err
			}
			err = os.WriteFile(filepath.Join(compDir, fmt.Sprintf("%s_horizontalpodautoscaler.yml", app)), data, 0644)
			if err != nil {
				level.Warn(g.logger).Log("msg", "failed to write HPA to file", "name", app)
				return err
			}
		}
	}
	return nil
}

func (g *Generator) writeManifestsToGivenFolder(app string, manifests *manifest.Manifests, parentDir string) error {
	if manifests.ServiceAccount != nil {
		data, err := runtime.Encode(g.encoder, manifests.ServiceAccount)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to encoding serviceAccount", "name", app)
			return err
		}
		err = os.WriteFile(filepath.Join(parentDir, fmt.Sprintf("%s_serviceaccount.yml", app)), data, 0644)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to write serviceAccount to file", "name", app)
			return err
		}
	}

	if manifests.ClusterRole != nil {
		data, err := runtime.Encode(g.encoder, manifests.ClusterRole)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to encoding clusterRole", "name", app)
			return err
		}
		err = os.WriteFile(filepath.Join(parentDir, fmt.Sprintf("%s_clusterrole.yml", app)), data, 0644)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to write clusterRole to file", "name", app)
			return err
		}
	}

	if manifests.ClusterRoleBinding != nil {
		data, err := runtime.Encode(g.encoder, manifests.ClusterRoleBinding)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to encoding clusterRoleBinding", "name", app)
			return err
		}
		err = os.WriteFile(filepath.Join(parentDir, fmt.Sprintf("%s_clusterrolebinding.yml", app)), data, 0644)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to write clusterRoleBinding to file", "name", app)
			return err
		}
	}

	if manifests.Role != nil {
		data, err := runtime.Encode(g.encoder, manifests.Role)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to encoding role", "name", app)
			return err
		}
		err = os.WriteFile(filepath.Join(parentDir, fmt.Sprintf("%s_role.yml", app)), data, 0644)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to write role to file", "name", app)
			return err
		}
	}

	if manifests.RoleBinding != nil {
		data, err := runtime.Encode(g.encoder, manifests.RoleBinding)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to encoding rolebinding", "name", app)
			return err
		}
		err = os.WriteFile(filepath.Join(parentDir, fmt.Sprintf("%s_rolebinding.yml", app)), data, 0644)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to write roleBinding to file", "name", app)
			return err
		}
	}

	if manifests.Ingress != nil {
		data, err := runtime.Encode(g.encoder, manifests.Ingress)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to encoding ingress", "name", app)
			return err
		}
		err = os.WriteFile(filepath.Join(parentDir, fmt.Sprintf("%s_ingress.yml", app)), data, 0644)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to write ingress to file", "name", app)
			return err
		}
	}

	for _, cm := range manifests.ConfigMaps {
		data, err := runtime.Encode(g.encoder, cm)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to encoding configMap", "name", cm.Name)
			return err
		}
		err = os.WriteFile(filepath.Join(parentDir, fmt.Sprintf("%s_cm.yml", cm.Name)), data, 0644)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to write configMap to file", "name", cm.Name)
			return err
		}
	}

	for _, secret := range manifests.Secrets {
		data, err := runtime.Encode(g.encoder, secret)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to encoding secret", "name", app)
			return err
		}
		err = os.WriteFile(filepath.Join(parentDir, fmt.Sprintf("%s_secret.yml", app)), data, 0644)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to write secret to file", "name", app)
			return err
		}
	}

	for _, svc := range manifests.Services {
		data, err := runtime.Encode(g.encoder, svc)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to encoding service", "name", svc.Name)
			return err
		}
		err = os.WriteFile(filepath.Join(parentDir, fmt.Sprintf("%s_service.yml", svc.Name)), data, 0644)
		if err != nil {
			level.Warn(g.logger).Log("msg", "failed to write service to file", "name", svc.Name)
			return err
		}
	}

	return nil
}

func (g *Generator) writeZip(source, dest string) error {
	archive, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer archive.Close()

	// New zip writer.
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	err = fs.WalkDir(os.DirFS(source), ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		// If a file is a symbolic link it will be skipped.
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// Create a local file header.
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Set compression method.
		header.Method = zip.Deflate

		// Set relative path of a file as the header name.
		header.Name = path
		header.Name = strings.ReplaceAll(header.Name, string(filepath.Separator), "/")

		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		}

		// Create writer for the file header and save content of the file.
		headerWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(filepath.Join(source, path))
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(headerWriter, f)
		return err
	})
	return err
}
