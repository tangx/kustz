# 介绍

![logo](/docs/static/logo/kustz.jpg)

如果要在 Kubernets 发布一个应用， 并对外提供服务， 需要配置诸如 `Dep, Ing, Svc` 等 Config API。 
他们之间又是通过 `Label` 组合选择而实现的 **松耦合**。

1. 如果想要这些 Config API 之间的关系更加紧密， 我们可以自己再向上抽象， 通过自己的配置将他们整合在一起。
2. 更重要的是， 我们可以通过这层抽象， 屏蔽不同版本 API 之间差异。 实现同一个 `kustz.yml` 配置兼容多版本集群， 实现部署或迁移的丝滑。

## Kustomize

> kustomize: https://kubectl.docs.kubernetes.io/guides/introduction/kustomize/

现在这个官网的引导页面看起来已经有点乱了。

简单的说， 以下是一个最基本的 `kustomization.yaml` 文件。

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: demo-demo
resources:
  - deployment.yml
  - service.yml
  - ingress.yml
```

1. `ApiVersion` 和 `Kind` : 对文件的作用定义
2. `Namespace` : 服务部署的运行环境。
3. `Resources` : 从外部引入的资源， 最终由 `kustomize` 统一渲染管理。 比如 patch 操作等。


## Deployment, Pod 和 Container

先来说说 Deployment， 这个应该是最常见的 **工作负载 workload**， **定义** Pod 状态 。

我们都知道， Pod 是 Kubernetes 中最小的 **调度** 单元， 定义网络、存储、权限等信息。 

换而言之， 最小的 **执行** 单元其实还是 Container， 定义了执行

![pod](/docs/img/pod.png)


通过 kubectl 命令，生成的最简单的 Deployment 模版。

```bash
$ kubectl create deployment my-nginx --image nginx:alpine --dry-run=client -o yaml
```

![dep-pod-c](/docs/img/dep-pod-container.jpg)

1. 最外层**红色**是 deployment.
2. 中间层**蓝色**是 pod.
3. 最内存**绿色**是 container.

没错， 他们之间的关系就是套娃， 一层套一层。


## kustz

`kustz` 的作用就如同 `Deployment` 一样， 将 **应用** 视作一个整体， 通过 **某种** 组织方式， 在一个文件中定义一个 **应用/服务**。

1. 将所有的配置都集中到 **同一个文件** 中， 多个资源更方便管理。
2. 将原本复杂的 API 结构 **语义化**， 配置起来更简单。

```yaml
# kustz.yml

namespace: demo-demo # 运行的命名空间

service:  # 定义一个应用
  name: srv-webapp-demo
  image: docker.io/library/nginx:alpine
  replicas: 1
  envs:   # 容器变量配置
    pairs:
      key1: value1-123
  resources:
    cpu: 10m/20m
    memory: 10Mi/20Mi

  probes: # 容器探针
    liveness:
      action: http://:80/liveness

  ports: # Service 端口
    - "80:80"       # cluster ip

## 对外暴露
ingress:
  - http://www.example.com/*
```

> 注意: 以上配置文件结构，可能会随着代码开发进行调整。 

如此这边，我们就可以通过 `kustz.yml` 定义完成一个服务的 **完整配置定义** ， 之后再通过 `kustz` 工具将起转化为 `kustomization.yml, deployment.yml ...` 等文件， 最后通过 `kubectl` 工具进行发布管理。

