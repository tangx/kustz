# 4.1. 使用 cobrautils 为命令添加更实用的命令参数

> 大家好， 我是老麦。


之前的章节， 我们陆陆续续给 kustz 库添加了很多丰富服务的配置

但 kustz 命令， 还是处于一个很原始的命令状态。 
接下来我们给 kustz 添加一些更丰富的参数 ， 使 kustz 用起来更顺手。

在 CICD 的中， 一般情况下 **变量，健康检查， 镜像策略** 等很难发生变动。 而镜像名称 **经常性** 的在每次打包后发生变化。 

每次CI触发都去修改 `kustz.yml` 配置显然是不可能的。 因此， 我们需要绑定更丰富的参数来支持我们 CI 的运行。

## cobra flag

之前在 `/cmd/kustz/cmd/render.go` 中， 我们为命令添加了一个指定配置文件的参数。

```go
func init() {
	cmdRender.Flags().StringVarP(&config, "config", "c", "kustz.yml", "kustz config")
}

var config string
```

这种方法是 cobra 官方提供的基本模式。 在绑定的时候， 需要一行写一个， 并且不支持 **指针参数** 。

## cobrautils 库

接下来我们使用自己封装的 `cobrautils` 库。

```bash
$ go get -u github.com/go-jarvis/cobrautils
```

详细描述参考 https://github.com/go-jarvis/cobrautils 

```go
func init() {
	cobrautils.BindFlags(cmdRender, flags)
}

// KustzFlag 定义 flag
type KustzFlag struct {
	Config   string `flag:"config" usage:"kustz config" shorthand:"c"`
	Image    string `flag:"image" usage:"image name"`
	Replicas *int   `flag:"replicas" usage:"pod replicas number"`
}

// 初始化默认值
var flags = &KustzFlag{
	Config: "kustz.yml",
}
```

可以看到， 使用 cobrautils 之后。 

1. 使用结构体组合了所有参数， 每个字段通过注释描述， 作用更清晰， 耦合度更高。
2. 支持 **指针参数**， 解决了 **零值** 带来的负面影响。
3. 一行命令解决了所有参数的绑定。

> 其实如果喜欢的话， 可以将 `/pkg/kustz/kustz.go` 中的完整 `Config` 做成参数。

在 `/cmd/kustz/cmd/render.go` 中， 将 `Image, Replicas` 两个参数注入到配置文件中即可。

```go

func render() {
	kz := kustz.NewKustzFromConfig(flags.Config)

	if flags.Image != "" {
		kz.Service.Image = flags.Image
	}
	if flags.Replicas != nil {
		kz.Service.Replicas = int32(*flags.Replicas)
	}

	kz.RenderAll()
}
```
