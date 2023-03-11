# Gin 中间件集合

## 功能

- gincors
  - CORS 支持
- gingzip
  - gzip 压缩支持
- ginlog
  - 替换 Gin 自带的 log middleware；加入了 requestID 字段，并支持链路追踪
- ginmetric
  - 导出 Gin 的 prometheus 监控数据
- ginrecovery
  - 替换 Gin 自带的 recovery middleware；优化了错误内容
- gintrace
  - 使 Gin 支持链路追踪，基于 opentelemetry

## 使用

```go
mws := []gin.HandlerFunc{ginmetric.Metrics(), gintrace.Trace(), ginlog.Log(true), ginrecovery.Recovery()}
r := gin.New()
r.Use(mws...)
```