# Kustz

使用 `kustomize` 简化 kubernetes 服务部署和配置

![logo](./img/kustz-logo.jpg)

## 目录

### 第一章 简介

+ [01. 简介](./01-introduce.md)

### 第二章 基础结构

+ [2.1. 模仿 kubectl create 创建 Deployment 样例](./02-1-sample-deployment.md)
+ [2.2. 定义字符串创建 Service](./02-2-define-strings-to-service.md)
+ [2.3. 解析 URL 为 Ingress](./02-3-parse-url-to-ingress.md)
+ [2.4. kustomize 流水线](./02-4-kustomize.md)
+ [2.5. 使用 cobra 实现 kustz 命令](./02-5-kustz-cli.md)

### 第三章 高级扩展

+ [3.1. 为 Container 添加环境变量](./03-1-container-env-var.md)
+ [3.2. ConfigMap 和 Secret 的生成器](./03-2-configmap-secret-generator.md)
+ [3.3. 注入 ConfigMap 和 Secrets 到容器环境变量](./03-3-container-env-from.md)
+ [3.4. 用字符串定义容器申请资源上下限](./03-4-container-resources.md)
+ [3.5. 为 Container 添加健康检查方法](./03-5-container-probe.md)
+ [3.6. 3.6. 镜像拉取鉴权和策略](./03-6-image-pull-policy.md)

### 第四章 投入使用
+ [4.1. 使用 cobrautils 为命令添加更实用的命令参数](./04-1-kustz-flags.md)

# 请一杯咖啡

![mp-weixin](./img/mp-qrcode.png)

如果你觉得这个项目还不错， 请我一杯咖啡 ☕️ 吧。

![wxpay](./img/pay/wxpay.png)
![alipay](./img/pay/alipay.jpg)
