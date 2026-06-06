---
name: review-pr
description: GitHub Pull Requestをレビューする。コード品質、バグ、セキュリティ、保守性、テスト観点からレビューを実施し、指摘事項を優先度付きでまとめる。
user-invocable: true
disable-model-invocation: true
---

# Review PR Skill

- このSkillはGitHub Pull Requestを日本語でレビューする。
- 必ず日本語でレビューすること

## Preparation

レビュー対象のPRを確認する。

ユーザーがPR番号を指定していない場合は質問する。

例:

- PR #123 をレビューして
- 現在のブランチのPRをレビューして

## Procedure

1. PR詳細を取得する。

```bash
gh pr view <PR_NUMBER>
```

2. 変更ファイル一覧を取得する。

```bash
gh pr diff <PR_NUMBER>
```

3. 必要に応じて個別ファイルを読む。

4. 変更内容を分析し、以下の観点でレビューを実施する。

5. レビュー内容を `outputs/pr/<PR_NUMBER>/review.md` に保存する。

## レビュー観点

以下の順で確認する。

### 1. 正確性

- ロジックバグ
- 境界値問題
- null/undefined考慮漏れ
- 例外処理不足

### 2. セキュリティ

- 認可漏れ
- 認証漏れ
- 秘密情報漏洩
- インジェクションリスク

### 3. パフォーマンス

- N+1問題
- 不要なループ
- メモリ効率
- DBアクセス

### 4. 保守性

- 命名
- 責務分離
- 重複コード
- 複雑度

### 5. テスト

- テスト不足
- テストケース漏れ
- 回帰リスク

## 出力フォーマット

### Summary

PRの概要を3〜5行で要約する。

### Findings

優先度付きで指摘する。

- Critical
- High
- Medium
- Low

各指摘について:

- 問題
- 影響
- 修正案

を記載する。

### Positive Feedback

良い実装や改善点も記載する。

## 重要ルール

- 推測でバグ扱いしない
- 差分から根拠を示す
- 指摘がなければ「重大な問題は見当たらない」と明示する
- スタイル指摘よりバグ・セキュリティを優先する
