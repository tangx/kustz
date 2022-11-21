# 2.5. 使用 cobra 实现 kustz 命令

> 大家好， 我是老麦。 一个运维小学生。

有了前面几章的努力， 我们的命令行工具 kustz 终于要问世了。

## kustz 命令

当前命令功能就很简单。

1. `default`: 输出 kustz **默认配置**。
2. `render`: 读取 kustz 配置并生成 kustomize 配置四件套。

```bash
$ kustz -h

Available Commands:
  default     在屏幕上打印 kustz 默认配置
  render      读取 kustz 配置， 生成 kustomize 所需文件
```


## 编码

本章的代码都很简单， 就是设计的文件比较多。

### 使用 cobra 创建命令

cobra 真的是一个非常好用的命令行工具。

```bash
go get -u github.com/spf13/cobra
```

在 `/cmd/kustz/cmd/root.go` 中创建 **根命令** `rootCmd`。
并定义 **执行函数** Execute。

```go
func Execute() error {
	return rootCmd.Execute()
}
```

在 `/cmd/kustz/cmd/default.go` 中创建 **子命令** `default`， 无任何参数。
在 `/cmd/kustz/cmd/render.go` 中创建 **子命令** `render`， 有一个参数 `config`， 实现根据配置管理应用。

而在外部 `/cmd/kustz/main.go` 中， 只有一个入口函数 `main` 调用 rootCmd 的执行。 保持文件清洁干爽。

### 

花开两朵。 在 `/pkg/kustz/cmd.go` 文件中， 提供了 **函数** 或 **方法** 供之前的命令调用


#### default

使用 `go:embed` 将配置文件 `kustz.yml` 嵌入到应用中。 配置随着代码走， 测试分发两不误。

```go
//go:embed kustz.yml
var defaultConfig string

func DefaultConfig() {
	fmt.Println(defaultConfig)
}
```

这里只是简单的配置文件打印标准输出。 如果需要保存到文件， 用户可以自行使用 **重定向符**。


#### render

通过 RenderAll 方法将之间的 `Deployment, Ingress, Service, Kustomization` 都保存成了对应文件。

```go
func (kz *Config) RenderAll() error {
}
```

在 Kustomize 章节已经硬编码 Resources 的资源文件名称。 因此这里可就定义了这几个文件名的常量。

```go
const (
	FileDeployment    = "deployment.yml"
	FileIngress       = "ingress.yml"
	FileService       = "service.yml"
	FileKustomization = "kustomization.yml"
)
```

在 `/pkg/kubeutils/yaml.go` 中， 将 `dep, svc, ing` 等编码成 YAML 的时候， 用到了 `sigs.k8s.io/yaml` (k8syaml) 库， 而非 `gopkg.in/yaml.v2` (pkgyaml)。
跟踪一下 k8syaml 的代码就很容易知道， 前者在 pkgyaml 的基础上， 针对性的为 k8s 做了很多优化。


## 编译及测试

```bash
$ go build ./cmd/kustz

$ ./kustz default > abc.yml
$ ./kustz render -c abc.yml
```

如果不行， 检查一下代码吧。

