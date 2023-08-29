package manifest

import (
	"regexp"

	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/udmire/observability-operator/pkg/templates/template"
	"github.com/udmire/observability-operator/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Builder interface {
	Build() *CapsuleManifests
}

func New(template *template.AppTemplate) Builder {
	if template == nil {
		return nil
	}
	return &templateBuilder{
		template: template,
	}
}

type templateBuilder struct {
	template *template.AppTemplate
}

func (b *templateBuilder) Build() *CapsuleManifests {
	cms := &CapsuleManifests{}
	var capsule []*Capsule
	files := make(map[string][]byte)
	for _, tempFile := range b.template.TemplateFiles {
		if tempFile.FileName == CapsuleFile {
			capsule = b.buildCapsules(tempFile.Content)
			continue
		}
		files[tempFile.FileName] = tempFile.Content
	}
	appLabels := appLabels(b.template.Name)
	cms.Manifest = *b.buildManifests(appLabels, capsule, files)

	for _, comp := range b.template.Workloads {
		cms.CompsManifests = append(cms.CompsManifests, b.buildComp(b.template.Name, comp))
	}

	return cms
}

func (b *templateBuilder) buildComp(app string, template *template.WorkloadTemplate) *CompManifests {
	compLabels := componentLabels(app, template.Name)

	var capsules []*Capsule
	files := make(map[string][]byte)
	for _, tempFile := range b.template.TemplateFiles {
		if tempFile.FileName == CapsuleFile {
			capsules = b.buildCapsules(tempFile.Content)
			continue
		}
		files[tempFile.FileName] = tempFile.Content
	}

	manifests := b.buildManifests(compLabels, capsules, files)
	return &CompManifests{
		Manifest: *manifests,
	}
}

func (b *templateBuilder) buildManifests(labels map[string]string, capsules []*Capsule, refs map[string][]byte) *Manifest {
	if len(capsules) < 1 {
		return nil
	}

	manifests := &Manifest{}

	for _, cap := range capsules {
		switch cap.Type {
		case ConfigmapType:
			manifests.ConfigMaps = append(manifests.ConfigMaps, &core_v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cap.Name,
					Namespace: "",
					Labels:    labels,
				},
				BinaryData: b.buildBinaryDatas(cap, refs),
			})
		case SecretType:
			manifests.Secrets = append(manifests.Secrets, &core_v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cap.Name,
					Namespace: "",
					Labels:    labels,
				},
				Data: b.buildBinaryDatas(cap, refs),
			})
		}
	}

	return manifests
}

func (b *templateBuilder) buildBinaryDatas(cap *Capsule, refs map[string][]byte) map[string][]byte {
	result := make(map[string][]byte)

	if cap.DynamicItems != nil {
		regex := regexp.MustCompile(*cap.DynamicItems)
		for key, val := range refs {
			if regex.Match([]byte(key)) {
				result[key] = val
			}
		}
	}

	for key, value := range cap.Items {
		result[key] = refs[value]
	}

	return result
}

func (b *templateBuilder) buildCapsules(content []byte) []*Capsule {
	cap := []*Capsule{}
	err := yaml.Unmarshal(content, &cap)
	if err != nil {
		panic(err)
	}
	return cap
}

// func (b *templateBuilder) BuildComp(template *template.WorkloadTemplate) *CompManifests {
// 	manifests := &CompManifests{Name: template.Name}
// 	for _, tempFile := range template.TemplateFiles {
// 		resType, _ := recognize(tempFile)
// 		var err error

// 		switch resType {
// 		case ConfigMap:
// 			cm := &core_v1.ConfigMap{}
// 			err = yaml.Unmarshal(tempFile.Content, cm)
// 			if err != nil {
// 				panic(err)
// 			}
// 			manifests.ConfigMaps = append(manifests.ConfigMaps, cm)
// 		case Secret:
// 			sec := &core_v1.Secret{}
// 			err = yaml.Unmarshal(tempFile.Content, sec)
// 			if err != nil {
// 				panic(err)
// 			}
// 			manifests.Secrets = append(manifests.Secrets, sec)

// 		default:
// 		}
// 	}

// 	return manifests
// }

//	func recognize(file *template.TemplateFile) (ManifestType, string) {
//		for i := 0; i < len(filePatterns); i++ {
//			match := regexp.MustCompile(filePatterns[i]).FindStringSubmatch(file.FileName)
//			if len(match) > 1 {
//				firstGroupContent := match[1]
//				return ManifestTypes[i], firstGroupContent
//			}
//		}
//		return -1, ""
//	}

func appLabels(app string) map[string]string {
	return map[string]string{
		utils.AppLabel: app,
	}
}

func componentLabels(app, component string) map[string]string {
	return map[string]string{
		utils.AppLabel:       app,
		utils.ComponentLabel: component,
	}
}
