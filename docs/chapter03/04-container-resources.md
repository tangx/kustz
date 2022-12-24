# 3.4. 用字符串定义容器申请资源上下限

> 大家好， 我是老麦， 一个运维小学生。
> 今天我们来给 kustz 添加资源申请字段。

![logo](/docs/img/kustz-logo.jpg)

Pod 的资源申请， 在调度策略中， 是一个重要的参数数据。 因此其重要性自然不必多说


## 容器资源申请

在官网中， 对于资源的申请和管理有详细的描述。 
https://kubernetes.io/zh-cn/docs/concepts/configuration/manage-resources-containers/

和 **服务质量 QoS** 息息相关， https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/quality-service-pod/

这里简单的归类， 可以速记， 按照服务质量高到低

1. Guaranteed: request = limit
2. Burstable: request < limit
3. BestEffort: 没有 request 和 limit


## `kustz.yml` 配置

还是先来看看 `kustz.yml` 配置中， 资源的抽象 `name: request/limit`。

```yaml
# kustz.yml
service:
  resources:
    cpu: 10m/20m
    memory: 10Mi/20Mi
  #   nvidia.com/gpu: 1/1
```

对应的， 在 `/pkg/kustz/kustz.go` 中， 增加字段， 如下。

```go
type Service struct {
	Resources map[string]string `json:"resources,omitempty"`
}
```


## 编码

这部分的编码还是还是很简单的。 因为在 `k8s.io/apimachinery/pkg/api/resource` 库中， 已经为我们提供了 **数据及单位** 的解析封装， 直接调用即可。

所以大部分编码工作还是在字符串的解析上面。

第一步， 在 `/pkg/kustz/k_container.go` 中， 为容器添加资源字段

```go
func (kz *Config) kubeContainerResources() corev1.ResourceRequirements {
	return tokube.ContainerResources(kz.Service.Resources)
}
```

第二步， 在 `/pkg/tokube/container_resources.go` 中， 解析 map[string]string 对象， 获取资源名和值， 并返回 corev1.Resources

```go
func ContainerResources(res map[string]string) corev1.ResourceRequirements {
// ...省略
}
```


第三步， 在 `/pkg/tokube/container_resources.go` 中分割我们定义的字符串。

```go
func toResourceQuantity(value string) (request resource.Quantity, limit resource.Quantity) {

	re, li := "", ""
	parts := strings.Split(value, "/")
// ... 省略
	request = resource.MustParse(re)
	limit = resource.MustParse(li)

	return
}
```

可以看到， 这里使用了官方的函数 `resource.MustParse(value)` ， 直接返回结论， 省了很多事情。


## 测试

执行命令， 查看结果。

```bash
$ make test.deployment
```

除了默认的 `cpu, memory` 之外， 我们还添加了 `nvidia/gpu` 这种 CRD 资源。

GPU 结论可以参考， https://help.aliyun.com/document_detail/94800.html
