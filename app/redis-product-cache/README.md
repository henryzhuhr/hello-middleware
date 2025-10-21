# Redis 电商商品缓存微服务测试

## 生成框架

```bash
goctl api go --api app/redis-product-cache/main.api --dir app/redis-product-cache --style gozero
```

## 启动代码

```bash
go run app/redis-product-cache/product.go -f app/redis-product-cache/etc/product.yaml
```
