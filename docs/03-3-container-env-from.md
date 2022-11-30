# 3.3. [kustz] 注入 ConfigMap 和 Secrets 到容器环境变量

![logo](./img/kustz-logo.jpg)

> 大家好， 我是老麦。 一个运维小学生。

有了前面两张的铺垫， 今天这个很简单。 我们说说另外一种为容器注入环境变量的方式。

## 容器变量注入 EnvFrom

前面我们提到过， Container 有两种方式定义环境变量， 其中一种就是 `envFrom`， 从 ConfigMap 或 Secret 中读取所有键值对作为容器的变量。

ConfigMap 和 Secret 看起来是这样的。 数据都在 data 字段。

```yaml
apiVersion: v1
data:
  APP_NAME: gin-demo
  LOG_LEVEL: debug
kind: ConfigMap
metadata:
  name: config-demo
```

在定义引用的时候， 使用 `envFrom` 关键字， 参考官网案例 https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/configure-pod-configmap/

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: dapi-test-pod
spec:
  containers:
    - name: test-container
      image: registry.k8s.io/busybox
      command: [ "/bin/sh", "-c", "env" ]
      envFrom:
      - configMapRef:
          name: special-config
          # optional: false
  restartPolicy: Never
```

官网的 demo 中并没有提及 optional 这个字段， 但是在后面 **限制条件** 的时候做了详细解释。


> 限制
1. 在 Pod 规约中引用某个 ConfigMap 之前，必须先创建这个对象， 或者在 Pod 规约中将 ConfigMap 标记为 optional（请参阅可选的 ConfigMaps）。 如果所引用的 ConfigMap 不存在，并且没有将应用标记为 optional 则 Pod 将无法启动。 同样，引用 ConfigMap 中不存在的主键也会令 Pod 无法启动，除非你将 Configmap 标记为 optional。
2. 如果你使用 envFrom 来基于 ConfigMap 定义环境变量，那么无效的键将被忽略。 Pod 可以被启动，但无效名称将被记录在事件日志中（InvalidVariableNames）。 日志消息列出了每个被跳过的键。例如:

> 可选的 ConfigMap 
你可以在 Pod 规约中将对 ConfigMap 的引用标记为 可选（optional）。 如果 ConfigMap 不存在，那么它在 Pod 中为其提供数据的配置（例如环境变量、挂载的卷）将为空。 如果 ConfigMap 存在，但引用的键不存在，那么数据也是空的。

> optional 默认值为 `false` , 即配置文件必须存在，否则会报错。


## 编码

今天的编码很简单， 就几句话。

## `kustz.yml` 配置

在 `service.envs` 中增加两个字段 `configmaps` 和 `secrets`。 他们都是 **字符串切片**。

```yaml
# kustz.yml

service:
  name: nginx
  image: docker.io/library/nginx:alpine
  envs:
    configmaps:
      # - name:optional
      - srv-webapp-demo-envs:true
    secrets:
      - srv-webapp-demo-envs # default optional: false
```

字符串分两段 `name:optional`



## 解析字符串为 API 对象

代码很简单， 没什么好说的。

1. 在 `/pkg/tokube/container_env.go` 中， 函数 `ParseEnvFromSource` 解析字符串为 `corev1.EnvFromSource` 对象。
2. 在 `/pkg/kustz/k_container_env.go` 中， 循环遍历 configmaps 和 secrets 获取字符串。

> 补充说明

在 `corev1.EnvFromSource` 中， 有个字段叫 `Prefix` 是用于给变量 **加前缀** 的。 

```go
// EnvFromSource represents the source of a set of ConfigMaps
type EnvFromSource struct {
	// An optional identifier to prepend to each key in the ConfigMap. Must be a C_IDENTIFIER.
	// +optional
	Prefix string `json:"prefix,omitempty" protobuf:"bytes,1,opt,name=prefix"`
}
```

我更喜欢 **所见即所得**， 所以我并没有处理这个字段。  另外环境变量是服务定义的， 也用不着我们画蛇添足。


## 测试

执行命令查看结果

```bash
$ make test.deployment
```
