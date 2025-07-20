# GEI
参考GIN格式的Web框架

---
# 使用
```go
func main() {
	r := gei.New()
	r.GET("/test", func(c *gei.Context) {
		c.JSON(http.StatusOK, gei.H{
			"hello": "world",
		})
	})

	_ = r.Run(":8080")

}
```
---

# 更新

## 动态路由
- 由Trie树实现


