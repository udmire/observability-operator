# 创建golang 项目
```
go mod init github.com/udmire/observability-operator
```

# 下载kubebuilder
```
cd ..
curl -OL https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
cd -
```

# 问Bito的问题及执行的操作

## 如何使用kubebuilder创建一个operator?
```
# 初始化Operator项目
../kubebuilder init --domain udmire.cn

# 创建Agent API
../kubebuilder create api --version v1alpha1 --kind Agents

# 创建Exporter API
../kubebuilder create api --version v1alpha1 --kind Exporters

# 生成应用模板结构
mkdir pkg/templates/app_version -p
cd pkg/templates/app_version/
touch app_configmap.yaml  app_hpa.yaml  app_ingress.yaml  app_rbac.yaml  app_secret.yaml  app_serviceaccount.yaml
mkdir comp
cd comp/
touch comp_configmap.yaml comp_deployment.yaml comp_ingress.yaml comp_secret.yaml comp_serviceaccount.yaml comp_daemonset.yaml comp_hpa.yaml comp_rbac.yaml comp_service.yaml comp_statefulset.yaml

```
* 
