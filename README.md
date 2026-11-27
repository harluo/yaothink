# `Gin`

[![编译状态](https://github.ruijc.com:20443/api/badges/pangum/pangu/status.svg)](https://github.ruijc.com:20443/pangum/pangu)
[![Golang质量](https://goreportcard.com/badge/github.com/pangum/pangu)](https://goreportcard.com/report/github.com/pangum/pangu)
![版本](https://img.shields.io/github/go-mod/go-version/pangum/pangu)
![仓库大小](https://img.shields.io/github/repo-size/pangum/pangu)
![最后提交](https://img.shields.io/github/last-commit/pangum/pangu)
![授权协议](https://img.shields.io/github/license/pangum/pangu)
![语言个数](https://img.shields.io/github/languages/count/pangum/pangu)
![最佳语言](https://img.shields.io/github/languages/top/pangum/pangu)
![星星个数](https://img.shields.io/github/stars/pangum/pangu?style=social)

基于`Gin`的Web服务器快速开发框架

## 查询汇率

在应用的 `config/rate.jsonc` 中通过 `yaothink.rate` 配置汇率服务（见[示例](config/rate.jsonc)）。
`yaothink.rate` 的 `hostname`、`port` 和 `prefix` 分别默认是 `rate-server.service`、`80` 和 `/api/private`。
客户端使用 HTTP，并在路径前缀后追加固定的 `/parities/singles`，默认请求地址为
`http://rate-server.service:80/api/private/parities/singles`。如已配置旧的 `endpoint`，请改用这三个字段。
`yaothink.rate.ttl` 默认是 `24h`（1 天），必须大于零。

在业务组件的构造函数中注入 `*yaothink.Rate`，然后调用 `rate.Value(ctx, yaothink.USD, yaothink.CNY)`。
返回值是以 `1_000_000` 为基数的整数，缓存按有方向的币种对保存成功结果。
