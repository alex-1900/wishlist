# wishlist
The wishlist SNS

# 开发规范
## 依赖项
依赖项包括：
- 自定义的 Service
- 配置、文件等资源
- 外部依赖、插件等

依赖项的使用:
- 依赖项在 `app/providers.go` 中定义，在 `app/dependency.go` 中实例化。
- 其他文件从 `app/dependency.go` 中导入依赖项。
