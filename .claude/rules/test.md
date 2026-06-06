---
paths: ["*_test.go"]
---

# Coding Guide for Go Tests

## Coding Style

- 対象ファイルのテストはカバレッジが100%になるように書いてください

## Table Driven Tests

- テストは Table Driven Tests の形式で書いてください

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        a, b, want int
    }{
        {1, 2, 3},
    }

    for _, tt := range tests {
        got := Add(tt.a, tt.b)
        if got != tt.want {
            t.Errorf("got %v want %v", got, tt.want)
        }
    }
}
```
