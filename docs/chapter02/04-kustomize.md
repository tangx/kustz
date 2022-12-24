# 2.4. kustomize 流水线

![logo](/docs/img/kustz-logo.jpg)

> 大家好， 我是老麦， 一个运维小学生。

前面已经简单的封装了 `Deployment, Service, Ingress`， 完成了零部件的创建。

今天就通过 `Kustomization` 进行组装， 实现流水线。

## Kustomize 

开始之前， 先来安装 kustomize 库。


```bash
$ go get sigs.k8s.io/kustomize/v3
```

这里补充一下， 访问 Github https://github.com/kubernetes-sigs/kustomize/。

kustomize () 首页 README.md 并没有提到 `go get` 的包名。 通常 k8s 的代码在 github 上都是镜像。 这时候只需要进到 `go.mod` ， 包名就一目了然。

```go
// go.mod
module sigs.k8s.io/kustomize/v3

go 1.12
```

## 编码

先来看看 `kustomization.yml` 的定义， 非常的简单。

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: demo-demo
resources:
  - deployment.yml
  - service.yml
  - ingress.yml
```


今天的代码及其简单， 只需要 20 行搞定。
在 import 的时候， 可能自动补全不会自己带上 `v3`。 需要手工调整一下。

```go
package kustz

import "sigs.k8s.io/kustomize/v3/pkg/types"

func (kz *Config) Kustomization() types.Kustomization {
	k := types.Kustomization{
		TypeMeta: types.TypeMeta{
			Kind:       types.KustomizationKind,
			APIVersion: types.KustomizationVersion,
		},
		Namespace: kz.Namespace,
		Resources: []string{
			"deployment.yml",
			"service.yml",
			"ingress.yml",
		},
	}

	return k
}
```

这里已经定了 kustomization 三个外部资源名字。


## 其它

kustomize 还是很贴心的， 在 types 把 version 和 kind 已经通过常量定义好了。

在 https://github.com/kubernetes-sigs/kustomize/blob/v3.3.1/pkg/types/kustomization.go 

```go
const (
	KustomizationVersion = "kustomize.config.k8s.io/v1beta1"
	KustomizationKind    = "Kustomization"
)
```

另外我们可以看到， 虽然 TypeMeta 定义相同， 但是直接从 ` apimachinery/pkg/apis/meta/v1.TypeMeta` 复制过来的， 而不是通过引用。

```go
// TypeMeta partially copies apimachinery/pkg/apis/meta/v1.TypeMeta
// No need for a direct dependence; the fields are stable.
type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}
```

> 之前看到一句话， **`简单的拷贝`比引用可能更节约资源， 因为引用是初始化一整个包**

## 测试

执行命令， 检查结果是不是和自己期待的一样。

```bash
$ go test -timeout 30s -run ^Test_KustzKustomize$ ./pkg/kustz/ -v
```

如果不是， 就回去检查代码吧。

