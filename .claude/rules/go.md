---
paths: ["*.go"]
---

# Coding Guide for Go

## Coding Style

- 可能な限り型推論を使用してください。明示的に型を書くのは最小限にしなさい
- `interface{}` ではなく `any` を使用してください
- エラーは必ず返すようにしてください。エラーを無視することは禁止されています
- 変数名は短く、意味のある名前を使用してください

## Typed nil

- typed nilを返さない･格納しない
- interface に格納された nil ポインタ `((*T)(nil))` は nil ではないため、nil 判定のバグを招く
- interface を返す場合は明示的に nil を返すこと。

```go
 // NG
var err *MyError = nil
return err

// OK
return nil // OK
```

## Syntax

### newによるポインタ変数の初期化

- Go 1.26からは、`new` を使用してポインタ変数を初期化することができるので、以下のように書いてください

```go
// NG
x := int64(300)
p := &x

// OK
p := new(int64(300))
```

### rangeの編集キャプチャは警戒しない

- Go 1.22からは、`range` のループ変数のキャプチャが編集されないようになったため、以下のように書いても問題ありません

```go
// NG
for _, v := range values {
  v := v
  go func() {
      fmt.Println(v)
  }()
}

// OK
for _, v := range values {
  go func() {
      fmt.Println(v)
  }()
}
```

## slices

- スライスの操作には `slice` パッケージを使用してください

```go
import "slices"

// 比較
slices.Equal(a, b)
// 検索
idx := slices.Index(items, target)
// 存在確認
slices.Contains(users, "alice")
```

## maps

- マップの操作には `maps` パッケージを使用してください

```go
import "maps"

// コピー
cloned := maps.Clone(src)

// 全要素のコピー
maps.Copy(dst, src)

// 比較
maps.Equal(a, b)

// キー一覧の取得
keys := slices.Collect(maps.Keys(m))

// 値一覧の取得
values := slices.Collect(maps.Values(m))
```

## errgroup

- 複数の goroutine の実行・待機には `errgroup` を使用してください
- `sync.WaitGroup` はエラー伝搬が不要な場合のみ使用してください
- `context.Context` と組み合わせて利用してください

```go
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(ctx)

g.Go(func() error {
    return taskA(ctx)
})

g.Go(func() error {
    return taskB(ctx)
})

if err := g.Wait(); err != nil {
    return err
}
```

## Logging

- ログ出力には標準ライブラリの `log/slog` を使用してください
- `fmt.Printf` や `log.Printf` を使用しないでください
- 構造化ログ（Structured Logging）で出力してください
- メッセージに値を埋め込まず、属性として出力してください

```go
import "log/slog"

logger.Info(
  "user created",
  "user_id", userID,
  "email", email,
)
```
