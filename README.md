# 如何开发一个k8s的operator

## 初始化项目

```shell
go mod init lfm-operator
```

## 使用code-generator

```shell
K8S_VERSION=v0.22.2
go get k8s.io/code-generator@$K8S_VERSION
go mod vendor
```

```shell
go get k8s.io/client-go@$K8S_VERSION
go get k8s.io/apimachinery@$K8S_VERSION
go get sigs.k8s.io/controllers-runtime@v0.10.3
```
