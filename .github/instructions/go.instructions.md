---
applyTo: '**/*.go'
---

# Go 编码规范

## 错误处理

### 错误的返回

返回错误的时候必须携带一定的上下文信息，方便排查问题，绝对禁止直接返回 err。

```go
if err != nil {
    return nil, err // Bad
}
if err != nil {
    return nil, fmt.Errorf("do somethin err: %w", err) // Good
}
```
