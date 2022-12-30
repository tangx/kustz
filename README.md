# Kustz 让应用在 Kubernetes 中管理更简单

![logo](/docs/static/logo/kustz.jpg)


## `kustz` 的设计思想和定义

`kustz` 的一个核心理念就是 **语义话**， 换句话说就是具有 **可读性** 高， **见名知义**。 

力求 `kustz.yml` 之于 **应用**， 就像 **域名** 至于 **IP**。

对于一个服务应用来说， 所有的定义都在同一个配置文件里面， 不再割裂。

从 [kustz 的完整配置][3] 中可以看到， 主要的参数都进行了 **语义化** 的处理和简化， 更贴近生活语言。

```yaml
## 1. k8s Deployment API 定义
  name: nginx
  image: docker.io/library/nginx:alpine
  replicas: 2
  envs:
    pairs:
      key1: value1
    configmaps:
      - srv-webapp-demo-envs:true
  resources:
    cpu: 10m/20m
    memory: 10Mi/20Mi
    nvidia.com/gpu: 1/1
  probes:
    liveness:
      action: http://:8080/healthy
```

```yaml
## 2. k8s Service API 定义
  ports:
    - "80:8080" # cluster ip
    - "udp://!9998:8889" # 随机 nodeport
    # - "!20080:80:8080" # 指定 nodeport
```

```yaml
## 3. k8s Ingress API 定义
ingress:
  rules:
    - http://api.example.com/ping?tls=star-example-com&svc=srv-webapp-demo:8080
```

### 说明: ConfigMap 配置说明

由于默认的 `kustomize` 的生成器支持 `k=v` 格式， 不支持多行变量。  因此使用 liternals 实现。

> 注意， 变量文件值支持 `YAML` 格式文件。

```yaml
# kustz.yml
# https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/secretgenerator/
secrets: # 或 configmaps
  literals:
    - name: srv-webapp-demo-literals
      files:
        - foo.yml
      # type: Opaque # default
```

变量文件

```yaml
## foo.yaml
JAVA_HOME: /opt/java/jdk
JAVA_TOOL_OPTIONS: -agentlib:hprof
HTTPS_CERT: |
  ---- RSA ----
  asdflalsdjflasdjfl
  ---- RSA END ----
```


既然现在的工具满足不了我们， 我们就自己抽象一层， 自己实现一个工具。

## 为什么会有 `kustz`

你有没有想过， 如果要在 kubernetes 集群中 **发布** 一个最基本的 **无状态服务**， 并 **提供** 给用户访问， 最少需要配置几个 `K8S Config API` ?

1. `Deployment`: 管理应用本身。 
2. `Service`: 管理应用在集群内的访问地址， 也是应用在在集群累的负载均衡器。
3. `Ingress`: 管理应用对外暴露的入口， 通俗点说， 就是 URL。

前三个是最基本的的 API。 

如果还有配置文件或或者其他密钥管理， 可能你还需要。

4. `Secret` / `ConfigMap`: 管理应用配置。

这些配置文件的存在， 本身都独立存在， 并没什么关系。

为了让他们在一起， 你还需要定义 `Label` 信息， 并且通过 `LabelSelector` 将他们组合起来。 

只是将这些 `Config API` 文件组合在一起， 都是一件麻烦事情了。 这还不包括各个 `Config API` 本身的复杂结构， 以及不同版本之间的差别。

社区也注意到这件事情了， 有很多工具帮我们组合管理， 例如我们今天要说的 [`Kustomize`][2]。 

除此之外， 还有微软和阿里云一起搞的 [`Open Application Model`(简称 `OAM`)][1]。

### Kustomize

下面是 `kustomize` 最基本的配置文件 `kustomization.yaml`

```yaml
# kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: demo-demo
resources:
  - deployment.yml
  - service.yml
  - ingress.yml
configMapGenerator:
- name: my-application-properties
  files:
  - application.properties
```

更多参数， 可以到 [`kustomize` 官网][2] 查看。 

可以看到 `kustomize` 也只是帮我们完成了文件的组合， 并没有解决 `Config API` 复杂结构的问题。


## 引用

[1]: https://oam.dev/
[2]: https://kubectl.docs.kubernetes.io/guides/introduction/kustomize/
[3]: https://github.com/tangx/kustz/blob/main/pkg/kustz/kustz.yml
[4]: https://tangx.in/books/kustz/chapter01/01-introduce/

