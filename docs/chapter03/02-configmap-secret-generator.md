# 3.2. [kustz] ConfigMap 和 Secret 的生成器

> 大家好， 我是老麦， 一个运维小学生。
> 今天我们通过 `kustomize` 管理 **ConfigMap** 和 **Secret**。

![logo](/docs/img/kustz-logo.jpg)

上一节我们通过 `k-v` 和 `YAML文件` 为容器添加环境变量。 同时也提到了可以通过 `envFrom` 这个关键字， 直接读取 ConfigMap 或 Secret 中的 `k-v` 作为容器的环境变量。

除了环境变量之外， ConfigMap 和 Secret 还能管理的东西还很多。 所以我个人觉得单应用管理部署的话， 对于配置的管理，还是比较重要的。


## Kustomize 中的 ConfigMap Env File

在 kustzomize 中， ConfigMap 和 Secret 都是通过 **生成器 Generator** 管理的， 有很多配置。

> https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/generatoroptions/


先切到 ConfigMapGenerator， 可以看到有三种模式提供数据, `files` , `literals`, `envs`。

如果按照我们之前说的， 为容器提供环境变量， 使用 `envs` 是最方便的。 从名字就可以看到， 就是为了环境变量而提供的。

> https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/configmapgenerator/#configmap-from-env-file

但是这种模式提供数据有也有限制的

1. 必须使用 `key=value` 这种结构
    1. **但这并不是 SHELL 变量赋值的形式**
2. 每一对 `k-v` 只能是单行。 key 作为变量名还好说， value 作为值就 **不能支持多行** 数据了。
    1. 另外 value 中的所有字符都是字面值。 

举个例子

```bash
HTTPS_CERT="---RSA---\nasdfjal\n---END"
```

通常在 Shell 中 `" 双引号` 是可以保留 `\n 换行符` 转义的含义的。 但是在这里 `" 和 \n` 都是字面意思， 没有任何特殊。


## ConfigMap / Secret 生成器

看看定义， ConfigMapArgs 和 SecretArgs 
1. 都是通过 `GeneratorArgs` 管理数据的。
2. Secret 比 ConfigMap 多了一个 Type 字段。

```go
type ConfigMapArgs struct {
	GeneratorArgs `json:",inline,omitempty" yaml:",inline,omitempty"`
}

// SecretArgs contains the metadata of how to generate a secret.
type SecretArgs struct {
	GeneratorArgs `json:",inline,omitempty" yaml:",inline,omitempty"`
	// This is the same field as the secret type field in v1/Secret:
	// It can be "Opaque" (default), or "kubernetes.io/tls".
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}
```

对于数据源我计划都从文件中读取。 这样三个模式就有了相同过的抽象结构体。 抽象结果为 Generator 结构体， 在 `/pkg/kustz/kustz.go` 中可以看到。


## `kustz.yml` 配置

1. 在 `kustz.yml` 中新增加了两个字段 `configmaps, secrets`。 
2. 每个字段都有三个模式 `envs, files, literals`。
3. 每个模式都有三个字段
    1. name: 最终生成的 ConfigMap 或 Secret 名字。
    2. files: 数据源。 `[target_name=]source_name`。 target_name 就是 ConfigMap 中的文件 key。 如省略， 默认与 source_name 相同。
    3. type: 类型。 Secret 专有。 取值范围参考 https://kubernetes.io/docs/concepts/configuration/secret/#secret-types

![types-of-secret](/docs/img/types-of-secret.jpg)

```yml
# kustz.yml

# https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/configmapgenerator/
configmaps:
  envs:
    - name: srv-webapp-demo-envs
      files:
        - src_name.txt

# https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/secretgenerator/
secrets:
  literals:
    - name: srv-webapp-demo-literals
      files:
        - foo.yml
      # type: Opaque # default
  files:
    - name: srv-webapp-demo-files
      files:
      - tls.crt=catsecret/tls.crt
      - tls.key=secret/tls.key
      type: "kubernetes.io/tls"
```


## 编码

```go
type Config struct {
	ConfigMaps Generator `json:"configmaps"`
	Secrets    Generator `json:"secrets"`
}

// Generator 定义数据源种类
type Generator struct {
	Literals []GeneratorArgs `json:"literals,omitempty"`
	Envs     []GeneratorArgs `json:"envs,omitempty"`
	Files    []GeneratorArgs `json:"files,omitempty"`
}

//  GeneratorArgs 定义数据源类型参数
type GeneratorArgs struct {
	Name  string   `json:"name,omitempty"`
	Files []string `json:"files,omitempty"`
	Type  string   `json:"type,omitempty"`
}
```

`Generator` 也就承担关于 ConfigMap 和 Secret 所有工作。

在 `/pkg/kustz/k_kustomize.go` 中， 为 Generator 创建了两个方法创建对应参数。

```go
// toConfigMapArgs 返回 ConfigMap 参数
func (genor *Generator) toConfigMapArgs() []types.ConfigMapArgs {
	args := []types.ConfigMapArgs{}
	for _, data := range genor.datas() {
		for _, garg := range data.gargs {
			arg := tokust.ConfigMapArgs(garg.Name, garg.Files, data.mode)
			args = append(args, arg)
		}
	}
	return args
}
// toSecretArgs 返回 Secret 参数
func (genor *Generator) toSecretArgs() []types.SecretArgs {}
```

由于 ConfigMap 和 Secret 确实太过相似， 所以对于处理 `GeneratorArgs` 使用循环， 从而添加了一个 `Mode` 类型的概念。 这个 Mode 的取值范围就是 `envs, literals, files`。

```go
type GeneratorArgsData struct {
	mode  tokust.GeneratorMode
	gargs []GeneratorArgs
}

// datas 整合生成器数据
func (genor *Generator) datas() []GeneratorArgsData {
	return []GeneratorArgsData{
		{mode: tokust.GeneratorMode_Envs, gargs: genor.Envs},
		{mode: tokust.GeneratorMode_Files, gargs: genor.Files},
		{mode: tokust.GeneratorMode_Literals, gargs: genor.Literals},
	}
}
```

在 `/pkg/tokust/generator.go` 文件中， 定义了几个函数创建 kustomize 对象的方法。 

```go
func ConfigMapArgs(name string, files []string, mode GeneratorMode) types.ConfigMapArgs {
}
func SecretArgs(name string, files []string, typ string, mode GeneratorMode) types.SecretArgs {
	// 处理默认类型
	if typ == "" {
		typ = "Opaque"
	}
}
```

相应的， 也创建三种模式的对应的方法。

```go
func generatorArgs_literals(name string, files []string) types.GeneratorArgs {
	data := make(map[string]string, 0)
	for _, file := range files {
		err := marshalYaml(file, data)
		if err != nil {
			panic(err)
		}
	}
	sources := mapToSlice(data)
	// ... 省略
}

func generatorArgs_files(name string, files []string) types.GeneratorArgs {
}

func generatorArgs_envs(name string, files []string) types.GeneratorArgs {
}
```

在 literals 中， 由于我们传入的是 **文件**， 但是在 `kustomization.yml` 是键值对。 
所以多了一个读取数据的步骤， 并且定义了一个规则， **如果出现同名变量， 后面出现的覆盖先出现的**。


## 测试

执行命令， 查看结果。

```bash
$ make test.kustomize
```

> 这里不会直接生成 ConfigMap 和 Secret， 而是生成 `Kustomization.yml` 规则。

