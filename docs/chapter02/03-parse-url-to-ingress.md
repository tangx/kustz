# 2.3. 解析 URL 为 Ingress


> 大家好， 我是老麦， 一个运维小学生。
> 今天我们处理 Ingress， 对外提供服务。

![logo](/docs/static/logo/kustz.jpg)

之前已经提到过， 在 `kustz.yml` 中的字段值， 要尽量做到 **见名知义**。

对于 Ingress 而言， 在发布之后， 我们访问的就是 URL 地址。 

```
http://api.example.com/v1
```

因此我们可以考虑 **从结果推导解析渲染 Ingress** 。

## Kubernetes Ingress

老规矩， 我们还是通过命令看看创建一个 ingress 需要提供哪些参数。

```bash
$ kubectl create ingress simple --rule="foo.com/bar=svc1:8080,tls=my-cert" -o
yaml --dry-run=client
```

在 rule 中， 提供了两组 k-v。 其中， `foo.com/bar` 就是一个不带协议的 URL。

再来看看输出结果。

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  creationTimestamp: null
  name: simple
spec:
  rules:
  - host: foo.com   # 多 host
    http:
      paths:
      - backend:    # 一个 host 多个后端服务
          service:
            name: svc1
            port:
              number: 8080
        path: /bar
        pathType: Exact
  tls:
  - hosts:         # 多个证书
    - foo.com
    secretName: my-cert
```

一个基本的 Ingress API 配置， 包含了

1. 主要由两个模块 `rules` 和 `tls` 构成。
2. rules 是个数组， 即 **一条或多条** URL。 
3. 每条规则可以有多个后端服务。
4. 规则路径有一个特殊的 `pathType` 参数， 表示规则是 **Prefix（前缀匹配）** 还是 **Excat（精确匹配）**

> 补充， 还有一个重要的模块 `annotations` 通过声明控制 `ingress-class` 行为。

```bash
$ kubectl create ingress annotated --class=default --rule="foo.com/bar=svc:port" --annotation ingress.annotation1=foo
```

## 编码

有了之前 Service 打样， Ingress 就容易很多了。


### `kustz.yml` 配置

之前已经提到过了， 我们希望 ingress rule 所见即所得。

```yaml
# ...省略
ingress:
  rules:
    - http://api.example.com/user/*?tls=my-cert&svc=srv-webapp-demo:8080
    - http://demo.example.com/login?tls=my-cert
  annotations:
    k1: v1
    k2: v2
```

1. rule 就是最终实际对外暴露的 URL
2. rule 通过 query 参数 `tls` 和 `svc` 传递后端证书名和服务。
3. 为了区分 `Prefix` 和 `Exact` 两种 PathType， 我使用了 **通配符 `*`**。

### Annotations

`Annotation` 还是很简单的， 本身就是 map 对象， 直接赋值就可以了。

```go
	ing := &netv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        kz.Name,
			Labels:      kz.CommonLabels(),
			Annotations: kz.Ingress.Annotations,
		},
		// ... 省略
	}
```


### Kubernetes 官方 IngressRule

在官方 `IngressRule` 结构体中规则描述还是挺多的。 截取部分

1. 需要满足 `RFC 3986` 规范
    1. host 字段不能是 IP 地址。 
    2. 不支持 80,443 之外的端口。 即只支持 `http` 和 `https` 协议。
2. host 可以带 `* 号`， 即支持泛域名。

> 文档中还专门强调了， 这些规则在以后可能会改变。 

### PathType

之前提到过， PathType 在控制匹配行为上非常重要。

1. Prefix： 前缀匹配， 也可以理解为模糊匹配。
2. Exact: 精确匹配， 只有100%匹配才生效。
3. **不建议使用** 还有一个值 `ImplementationSpecific`， 表示由 ingress-class 决定值是 **Prefix 或者 Excat**。 

PathType 其实虽说不难， 但是官网还是给了一个 Example 详细列举了匹配规则。 建议还是看一下。
https://kubernetes.io/docs/concepts/services-networking/ingress/#path-types

代码中

1. Type 默认类型是 `Excat`。 相对比较安全。
2. 通过判断 `path` 末尾最后一个字符是 `*` 则为 `Prefix` 规则。

```go
func NewIngressRuleFromString(value string) *IngressRule {
// ...省略
	// ex: /api/*
	path := ur.Path
	typ := netv1.PathTypeExact
	if strings.HasSuffix(path, "*") {
		path = strings.TrimSuffix(path, "*")
		typ = netv1.PathTypePrefix
	}
// ...省略
```

> 确认 Prefix 之后， 别忘记把 `*` 从 path 中去掉。

## 渲染 Ingress

在官方 Ingress 结构体中字段一层套一层， 少说五六层。

因此在代码中定义了一个 `IngressRuleString` 保存所需要的字段信息。

```go
type IngressRuleString struct {
	Host      string
	Path      string
	PathType  netv1.PathType
	TLSSecret string
	Service   string
}
```

1. 通过函数 `NewIngressRuleFromString` 解析字符串。

```go
func NewIngressRuleFromString(value string) *IngressRuleString {
}
```

2. 通过方法 `KubeIngressTLS`,`KubeIngressRule` 创建 tls 和 rule。

```go
func (ir *IngressRuleString) KubeIngressTLS() *netv1.IngressTLS {
}

func (ir *IngressRuleString) KubeIngressRule() netv1.IngressRule {
}
```


### IngressRuleString


> 注意， 在定义 IngressRuleString 的时候， 我偷了个懒。 

前面说过， ingress rule 是支持多个后端服务的, 所以 `Service` 应该是 **切片** 类型。

```go
type IngressRuleString struct {
	Service   []string
}
```

但在我的定义中是 **字符串** 类型。 

```go
// 渲染 Ingress
func (kz *Config) KubeIngress() *netv1.Ingress {
	rules, tlss := ParseIngreseRulesFromStrings(kz.Ingress.Rules, kz.Name)
	// ... 省略
}

// 解析字符串
func ParseIngreseRulesFromStrings(values []string, defaultService string) ([]netv1.IngressRule, []netv1.IngressTLS) {
	// ... 省略
	for _, value := range values {
		ing := NewIngressRuleFromString(value)
		if ing.Service == "" {
			ing.Service = defaultService
		}
	}

	// ... 省略
}
```

1. 在调用的时候传递了一个 **服务同名** 默认值。
2. **服务同名** 不包含端口， 即 `srv-webapp-demo`。 省略端口默认为 **`80`**

```go
func (ir *IngressRuleString) toKubeIngressBackend() netv1.IngressBackend {
	// srv-webapp-demo[:8080]
	svc := ir.Service
	port := int32(80)

	parts := strings.Split(svc, ":")
	if len(parts) == 2 {
		svc = parts[0]
		port = StringToInt32(parts[1])
	}
}
```

**这一处默认值** 反过来要求了 Service 的端口映射必须是 `80:xxxx`

```yaml
# kustz.yml
service:
  ports:
    - "80:8080" # cluster ip
```

> 为什么？
>> 因为当管理多个服务（尤其是多个团队或语言的服务）时， 统一使用 `80` 为 Service 入口也可以作为一个**强制规则约束**， 节省脑细胞。

## 测试

执行命令， 检查结果是不是和自己期待的一样。

```bash
$ go test -timeout 30s -run ^Test_KustzIngress$ ./pkg/kustz/ -v
```

如果不是， 就回去检查代码吧。

