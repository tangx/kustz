# 2.2. 定义字符串创建 Service

![](./img/kustz-logo.jpg)

> 大家好， 我是老麦， 一个小运维。
> 今天我们为 kustz 增加 service 解析功能。

通过 `kubectl create service` 命令可以看到， service 的模式还是挺多的。

```bash
$ kubectl create service -h
Create a service using a specified subcommand.

Aliases:
service, svc

Available Commands:
  clusterip      Create a ClusterIP service
  externalname   Create an ExternalName service
  loadbalancer   Create a LoadBalancer service
  nodeport       Create a NodePort service
```

除了以上列出来的四种之外， 还用一种 `Headless Service`(https://kubernetes.io/docs/concepts/services-networking/service/#headless-services)。

`Headless Service` 是当 **类型** 为 ClusterIP，且 IP 值为 `None`。 所以是 `Cluster IP` 的一种特殊情况。

```bash
# 创建一个新的 ClusterIP service
$ kubectl create service clusterip my-cs --tcp=5678:8080

# 创建一个新的 ClusterIP service (headless 模式)
$ kubectl create service clusterip my-cs --clusterip="None"
```

## Service API 中的 Port

如果你之前留意过 Service API， 你就应该会发现 API 中带有 port 的字段就有 3 个。 弄清楚他们分别对应什么**这点很重要**

先来看看一个 `NodePort` 的 API。

```yaml
# $ kubectl create service nodeport my-cs --tcp=80:8080 --dry-run=client -o yaml
apiVersion: v1
kind: Service
  # ... 省略
spec:
  ports:
  - name: 80-8080
    nodePort: 18080
    port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    app: my-cs
  type: NodePort
```

对于 `nodePort, port, targetPort` 这三个名词如果还不能直接回答的话， 认真记住下面这张图。

![nodeport](./img/nodeport-port-targetport.png)

1. 集群外部用户通过 `nodePort -> port -> targetPort`。 
2. 集群内部用户通过 `port -> targetPort`。

## kubez 中的 service 规则

在看随后的规则的时候， 再回头看看之前的 API， 能否找到几个关键要素。

1. 端口映射规则。
2. Service 类型。
3. 端口协议。 通常不太关注。
4. 明确且有意义区分符号。 这是一个隐藏点。

+ clusterip: `80:8080`。 无任何特殊符号。
+ headless: `#80:8080`。 **注释符号 #** 顾名思义， 就是把某些 Head **隐藏** 之后成为 Headless。
+ nodeport: `!18080:80:800` / `!80:8080`。 **叹号 !** 表示重要， 因为我们要暴露端口到外部， 这是一个风险点。 NodePort 有两个规则， 是因为如果 **不指定** 端口就使用 **随机** 端口。
+ externalname: `@example.com`。 **AT @**， 指定或指向某个地方。
+ loadbalancer: `%80:8080`。 **分号 %**， 看起来就像个天平， 一人一半， 分流。

上面的规则中， 并没有提到 **端口协议**。

+ 含有协议的规则 `tcp://80:8080` / `udp://!18080:80:8080`。 就用协议地址的写法， 使用 `://` 分割 **协议** 和 **转发规则**。

之所以协议拿出来单独说， 其原因是协议通常不太关注， 默认都用 TCP。因此在 kustsz 中也作为可省略字段。

> 以上规则都是个人习惯。 仅供参考。

## 代码实现


凑 API 还是很简单的， 不过需要注意当出现**多个规则**的时候**怎样**确定**最终**的 `Type`。

例如规则中存在 `ClusterIP, NodePort 和 ExternalName`， 那在最终呈现的时候， Type 值是顺序第一？逆序第一？还是冲突报错？ 一切看自己需求和习惯。

```yaml
ports:
  - "8099:80"
  - "!28080:80"
  - "@example.com"
```

本章代码就先只实现两个， ClusterIP 和 NodePort。 完整配置查看 `/pkg/kustz/kustz.yml`

```yaml
# ... 省略
service:
# ... 省略
  ports:
    - "80:8080" # cluster ip
    - "udp://!9998:8889" # 随机 nodeport
    # - "!20080:80:8080" # 指定 nodeport
```

### 先说一个变更

在 `/pkg/kustz/common.go` 中， 之前的 `func CommonLabels(kz *Config)` 函数变为了 `func (kz *Config) CommonLabels()` 方法。

```go
func (kz *Config) CommonLabels() map[string]string {
	return map[string]string{
		"app": kz.Name,
	}
}
```

### 生成 Service API

在创建 Service API 的方法中还是中规中矩的在凑 API。

```go

func (kz *Config) KubeService() *corev1.Service {

	ports, typ := ParsePortStrings(kz.Service.Ports)

	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   kz.Name,
			Labels: kz.CommonLabels(),
			// Namespace: kz.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: kz.CommonLabels(),
			Type:     typ,
			Ports:    ports,
		},
	}

	return svc
}
```

在 ObjectMeta 中可以注意到， `Namespace` 字段是被注释了的。 

1. 因为 namespace 是**环境限定**， 而非服务本身特性。 即部署在 ns-a 和 ns-b 中， 服务配置本身并没有改变。
2. 对于 namespace 值提供， 我想放在 `kustomization.yml` 中, 通过 `kubectl apply -n ns-demo` 提供。

> 同样的， 在 `/pkg/kustz/k_deployment.go` 中的 Deployment API 的 ObjectMeta 也删除了 namespace 字段。


### PortString

从 `corev1.ServiceSpec` 中可以看到， Type 和 Ports 是并列关系， 且只有一个 Type

```go
		Spec: corev1.ServiceSpec{
			Selector: kz.CommonLabels(),
			Type:     typ,
			Ports:    ports,
		},
```

而在我们 kustz.yml 文件中, `!9998:8889`,  Type **符号** 存在于每一条规则中。 

因此， 通过函数 `ParsePortStrings` 解析所有规则， 并返回一个并列关系。

```go
func ParsePortStrings(values []string) ([]corev1.ServicePort, corev1.ServiceType) {
// ...
	for _, value := range values {
		port := NewPortFromString(value)
		if port.Type != corev1.ServiceTypeClusterIP {
			typ = port.Type
		}
		sps = append(sps, port.KubeServicePort())
	}

	return sps, typ
}
```

在 for 循环中， 处理了当 `ports` 出现多个 Type 规则的时候以谁为准的问题， 正如我们之前提到的。

这里我使用了 **Type 默认为 ClusterIP， 且最后出现的 `非ClusterIP` 为准（如有）** 。


为此， 我定义一个命令 `PortString` 的结构体， 

1. 通过 `NewPortFromString` 函数从 **字符串** 中提取关键信息。
2. 并通过 `PortString KubeServicePort` 方法将结构体转换成的 `corev1.ServicePort`


```go
type PortString struct {
	Port       int32
	TargetPort int32
	NodePort   int32
	Protocol   corev1.Protocol
	Type       corev1.ServiceType
}

// NewPortFromString parse port from string to PortString
func NewPortFromString(value string) PortString {}

// KubeServicePort return a corev1.ServicePort
func (p *PortString) KubeServicePort() corev1.ServicePort {}
```

`PortString` 提供了相应的方法按规则从 `string` 提取关键信息。

```go
// toServiceClusterIP parse value from for ClusterIP
func (p *PortString) toServiceClusterIP(value string) {

	parts := strings.Split(value, ":")
	switch len(parts) {
	case 1:
		n := p.StringToInt32(parts[0])
		p.Port = n
		p.TargetPort = n
	case 2:
		p.Port = p.StringToInt32(parts[0])
		p.TargetPort = p.StringToInt32(parts[1])
	}

	p.Type = corev1.ServiceTypeClusterIP
}

// toServiceNodePort parse value from for NodePort
func (p *PortString) toServiceNodePort(value string) {

	value = strings.TrimPrefix(value, "!")
	parts := strings.Split(value, ":")
	switch len(parts) {
	case 1, 2:
		p.toServiceClusterIP(value)
	case 3:
		p.NodePort = p.StringToInt32(parts[0])
		p.Port = p.StringToInt32(parts[1])
		p.TargetPort = p.StringToInt32(parts[2])
	}

	p.Type = corev1.ServiceTypeNodePort
}
```

## 测试

在 `/pkg/kustz/kustz_test.go` 中， 拆分了 test 规则。

执行命令， 检查结果是不是和自己期待的一样。

```bash
$ go test -timeout 30s -run ^Test_KustzService$ ./pkg/kustz/ -v
```

如果不是， 就回去检查代码吧。

