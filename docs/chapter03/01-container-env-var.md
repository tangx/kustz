# 3.1. 为 Container 添加环境变量

> 大家好， 我是老麦。 一个运维小学生。
> 今天为容器添加环境变量。

![logo](/docs/static/logo/kustz.jpg)

再前面一章中， 我们已经完成了 `Deployment, Service, Ingress 和 Kustomization` API 的封装。 
并通过 `cobra` 库创建了属于我们自己的 `kustz` 命令。

然而 kustz 的功能还简陋。 今天我们就先来为容器添加环境变量。


## 为容器设置环境变量

在官方文档中， 提高了两种为容器设置环境变量的方法

> https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/

1. `env`: 提供 `k-v 模式` 键值对。
    1. 值可以直接 `value` 提供。
    2. 也可以通过 `valueFrom` 从 secret 或 configmap 引用。
2. `envFrom`: 从 secret 或 configmap 中读取键值对， 注入到容器中。。
    1. https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/


### `kubez.yml` 配置

首先来看看 `kubez.yml` 的配置

```yaml
# kubez.yml
service:
  envs:
    pairs:
      key1: value1
    files:
      - foo.yml
      - bar.yml
```

我设计了两种方式为容器提供环境变量。 都是提供 k-v 模式。

```yaml
# deployment.yml
    env:
    - name: DEMO_GREETING
      value: "Hello from the environment"
```

1. `pairs`: k-v 模式。 **优先级更高**， 可以覆盖 files 中出现的同名 k-v。
2. `files`: 从文件中读取 k-v。 
    1. 多个 `kustz.yml` 可以复用。
    2. 可以按类型分类， **更直观**。 例如工程变量和数据库变量。
    3. 选择 `YAML` 格式是为了更好的管理 **值为多行的变量**。 比如证书。
    4. **同名变量，后者覆盖前者**

> 挖个坑， 以后实现 `2.2` 中提到的数据库变量文件的加解密。 让 GitOPS 更安全一点。


**最后强调一下变量优先级顺序**， 用链条表示: **后者覆盖前者**。

```
foo.yml <- bar.yml <- pairs
```

### 编码实现

在 `/pkg/kustz/kustz.go` 中， 增加配置字段。  这个很简单， 就不赘述了。

```go
type ServiceEnvs struct {
	Pairs map[string]string `json:"pairs,omitempty"`
	Files []string          `json:"files,omitempty"`
}
```

在引用 `ServiceEnvs` 的时候没有始终指针。 这样 Service 在初始化的时候， ServiceEnvs 即使在 `kustz.yml` 没有定义也会被初始化为 **空** 的零值。

```go
type Service struct {
	Envs     ServiceEnvs `json:"envs,omitempty"`
}
```

代码还是很简单的， 在 `/pkg/kustz/k_container.go` 定义了方法 `kubeContainerEnv`。
该方法中， 指定了变量的优先级。

```go
func (kz *Config) kubeContainerEnv() []corev1.EnvVar {
	pairs := make(map[string]string, 0)

	for _, file := range kz.Service.Envs.Files {
		b, _ := os.ReadFile(file)
		_ = yaml.Unmarshal(b, &pairs)
	}

	for k, v := range kz.Service.Envs.Pairs {
		pairs[k] = v
	}

	return tokube.ContainerEnv(pairs)
}
```

新建了一个包 `tokube`， 这里面的函数通过 **接受** 参数返回 API 配置。 

在 `/pkg/tokube/container_env.go` 定义了函数 `ContainerEnvs` 创建容器变量键值对。

```go
func ContainerEnv(pairs map[string]string) []corev1.EnvVar {
	envs := []corev1.EnvVar{}
	for k, v := range pairs {
		envs = append(envs, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	return envs
}
```

### 说一个容易错的点

在 `Container API` 中， 变量保存在 `[]corev1.EnvVar`， 这是一个切片。 切片的另一个 **隐藏含义** 就是可能出现 **同名 KEY**。

第一次的代码如下， 
```go
func (kz *Config) kubeContainerEnv_Error() []corev1.EnvVar {
	envs := []corev1.EnvVar{}

// 注意这里所有变量全部假如了 envs 切片。
	envs = append(envs, tokube.ContainerEnv(kz.Service.Envs.Pairs)...)
	for _, file := range kz.Service.Envs.Files {
		b, _ := os.ReadFile(file)
		mm := make(map[string]string, 0)
		_ = yaml.Unmarshal(b, &mm)
		envs = append(envs, tokube.ContainerEnv(mm)...)
	}

	return envs
}
```

![error](/docs/img/env-duplicate-key.jpg)

并没有测试这样的变量结构会出现什么情况， 因为这种情况就不应该出现。


## 测试

觉得之前的测试命令不方便， 更新了 Makefile， 添加了测试命令。
执行命令测试吧。

```bash 
$ make test.deployment
```