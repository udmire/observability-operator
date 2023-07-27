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

# 初始化operator项目
```
../kubebuilder init --domain udmire.cn
```

# 如何使用kubebuilder创建一个operator

