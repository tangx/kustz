# 2.1. 模仿 kubectl create 创建 Deployment 样例

![](/docs/static/logo/kustz.jpg)

为了简单， 我们假定所管理的 Deployment 都是 **单容器** 的。 

首先参考 `kubectl create` 命令

```bash
$ kubectl create deployment my-dep --image=busybox --replicas 1 --dry-run=client -o yaml
```

## 安装 client-go API

访问 `client-go` https://github.com/kubernetes/client-go

```bash
$ go get k8s.io/client-go@v0.25.4
```

这里直接选用最新版本 `v0.25.4`。 对于其他版本的兼容， 留在以后再做。


## 定义 Kustz Config

参考 `kubectl create` 命令， 创建配置文件 `kustz.yml` 结构如下

```yaml
# kustz.yml
namespace: demo-demo
name: srv-webapp-demo

service:
  name: nginx
  image: docker.io/library/nginx:alpine
  replicas: 2

```

在 service 中添加了 name 字段， 这是在 `kubectl create` 命令中没有的。 后者直接使用了镜像名称作为 name 的值。

> 由于我们的设计只有一个容器， 这里也可以 **省略或使用默认值**。 这里增加字段 **完全是为了展示凑 API**。

在 `/pkg/kustz/kustz.go` 中创建 `Config`， 对应所有字段。

```go
type Config struct {
	Name      string  `json:"name"`
	Namespace string  `json:"namespace"`
	Service   Service `json:"service"`
}

type Service struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	Replicas int32  `json:"replicas"`
}
```

## 拼凑 Deployment API

从这个标题就可以看出来， 这里就就没什么难度了， 就是把 `kustz.Config` 中的所有字段全部放在 `Deployment API` 中。

为 `kustz.Config` 添加 `KubeDeployment` 方法， 在 `/pkg/kustz/k_deployment.go` 中。

```go
func (kz *Config) KubeDeployment() *appv1.Deployment {
	return &appv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kz.Name,
			Namespace: kz.Namespace,
			Labels:    CommonLabels(*kz),
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &kz.Service.Replicas,
			Template: kz.KubePod(),
			Selector: &metav1.LabelSelector{
				MatchLabels: CommonLabels(*kz),
			},
		},
	}
}
```

可以看到， 拼凑主要由 3 部分， 都是很熟悉的字段。

1. `TypeMeta`: **Deployment** 的信息申明
2. `ObjectMeta`: **应用服务** 的信息申明
3. `Spec`: 就是具体信息了。

如果你看的够仔细，可以在 `Spec` 中发现 `Template` 字段就开始 **套娃** 了。

`KubeDeployment` 中调用了 `/pkg/kustz/k_pod.go` 中的 `KubePod` 方法。

```go
func (kz *Config) KubeDeployment() *appv1.Deployment {
	return &appv1.Deployment{
		Spec: appv1.DeploymentSpec{
			// ... 省略
			Template: kz.KubePod(),
		},
	}
}
```

`KubePod` 中调用了 `/pkg/kustz/k_container.go` 中的 `KubeContainer` 方法。

```go
func (kz *Config) KubePod() corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		// ... 省略
		Spec: corev1.PodSpec{
			Containers: kz.KubeContainer(),
		},
	}
}
```

最后， 在 `KubeContainer` 方法中， 创建了最里面的 container 信息。

```go
func (kz *Config) KubeContainer() []corev1.Container {
	if kz.Service.Name == "" {
		kz.Service.Name = kz.Name
	}

	// ...省略
	return []corev1.Container{c}
}
```

前面我们说到过 `kz.Service.Name` 为了展示 **API** 的拼凑。 为了以后使用方便， 这里我们设置其默认值为应用服务名 `kz.Name` 。

### 公共标签

在 Kubernetes 中， 不同 API 之间的关系都是通过 **标签选择** 关联的。 
为了方便， 在 `/pkg/kustz/common.go` 中使用函数 `CommonLabels()` 创建公共标签， 方便关联。


## 文件解析

细心的你可以已经发现了， 明明配置文件用的是 `yaml` 格式， 但是在 `Config` 中的标签却是 `json:"name"`。

这里只是为了 **单纯** 为了保障 yaml 库一致， 在 `/pkg/kubeutils/yaml.go` 中使用了 `sigs.k8s.io/yaml` 库而已, 这个库可以兼容 `json, yaml` 标签。


## 测试

进入 `kustz`， 执行命令

```bash
$ go test -v .

$ kubectl create deployment srv-webapp-demo --image=nginx -n demo-demo --dry-run=client -o  yaml 
```

对比二者内容， 基本上完全一样了。


