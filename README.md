```go
mws := []gin.HandlerFunc{ginmetric.Metrics(), gintrace.Trace(), ginlog.Log(true), ginrecovery.Recovery()}
r := gin.New()
r.Use(mws...)
```