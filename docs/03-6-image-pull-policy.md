# 3.6. 镜像拉取鉴权和策略

> 大家好， 我是老麦。 
> 今天我们解决镜像拉取鉴权和策略

![logo](./img/kustz-logo.jpg)


## 镜像拉取鉴权

拉取私有镜像或私有仓库镜像的时候， 需要提供鉴权信息。 

在 Kubernets 中， 通过 Secret 管理账号这些账号信息。 Secret 类型分为两种， 

1. `kubernetes.io/dockerconfigjson`: 如果有linux安装了 docker， 就是 `~/.docker/config.json` 这个文件。

2. `kubernetes.io/dockercfg`: 不太熟。

## 镜像拉去策略

镜像拉去策略分为三种，  `Never, Always, IfNotPresent`

在 `/pkg/tokube/container.go` 中， 可以看到 `ImagePullPolicy` 的处理方法。

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
