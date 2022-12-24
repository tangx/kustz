# 3.6. 镜像拉取鉴权和策略

> 大家好， 我是老麦。 
> 今天我们解决镜像拉取鉴权和策略

![logo](/docs/img/kustz-logo.jpg)


## 镜像拉取鉴权

拉取私有镜像或私有仓库镜像的时候， 需要提供鉴权信息。 

在 Kubernets 中， 通过 Secret 管理账号这些账号信息。 Secret 类型分为两种， 

1. `kubernetes.io/dockerconfigjson`: 如果有linux安装了 docker， 就是 `~/.docker/config.json` 这个文件。

2. `kubernetes.io/dockercfg`: 不太熟。

在 `/pkg/tokube/pod.go` 中， 可以看到 `ImagePullSecrets` 的处理方法。 就是将字符串转为 kubernetes 的引用对象， 其它没什么好说的。

```go
func ImagePullSecrets(secrets []string) []corev1.LocalObjectReference {
	if len(secrets) == 0 {
		return nil
	}
	objs := []corev1.LocalObjectReference{}
	for _, s := range secrets {
		objs = append(objs, corev1.LocalObjectReference{
			Name: s,
		})
	}
	return objs
}
```

## 镜像拉去策略

镜像拉去策略分为三种，  `Never, Always, IfNotPresent`

在 `/pkg/tokube/container.go` 中， 可以看到 `ImagePullPolicy` 的处理方法。

```go
func ImagePullPolicy(s string) corev1.PullPolicy {
	switch strings.ToLower(s) {
	case "always":
		return corev1.PullAlways
	case "never":
		return corev1.PullNever
	case "ifnotpresent":
		return corev1.PullIfNotPresent
	}
	return ""
}
```

1. 在 `kustz.yml` 不再大小写敏感， 因为我们将值全部转为小写。
2. 当不指定配置策略的时候， 使用默认策略。

## 使用

如果在 `kustz.yml` 配置中， 通过如下配置。

假设配置文件名为 `docker-config.json`， 支持多个账号， 参考如下。

```json5
// docker-config.json
{
	"auths": {
		"ghcr.io": {
			"auth": "Abcdefg=="
		},
		"https://index.docker.io/v1/": {
			"auth": "Abcdefg="
		}
	}
}
```

`auth` 值是 `user:password` 的 base64 编码。 如果不知道怎么弄 `docker login` 生成

```bash
$ docker login -u myUser -p myPass
```

在 `kustz.yml` 中， 通过 `docker-config.json` 创建 Secret 并引用。

```yaml
service:
  imagePullSecrets:
    - aliyun-repo

secrets:
  files:
    - name: aliyun-repo
      files:
        - .dockerconfigjson=docker-config.json
      type: kubernetes.io/dockerconfigjson
```
