---
name: make-pr
description: GitHub Pull Requestを作成する。変更内容を分析し、適切なタイトルと本文を生成してghコマンドでPRを作成する
---

# Make PR Skill

- このSkillはGitHub Pull Requestを作成するためのものです
- コミットされていない変更は無視して良いです
- マージ先は `main` や `develop` とは限らないので、ユーザーに必ず確認すること

## Procedure

1. 必ずユーザーからどのブランチをマージ先にするか確認する。`git branch -a` でそのマージ先ブランチが存在するか確認し、正しいブランチ名を確認する。

```bash
git branch -a
```

2. 現在のブランチを確認

```bash
git status
```

3. コミット済みの変更とベースブランチとの差分を確認

- コミットされていない変更は無視してください。必ずコミット済みのdiffを確認してください

```bash
git diff origin/<ベースブランチ>...HEAD
git log --oneline origin/<ベースブランチ>..HEAD
```

3. 変更内容を分析し、以下を作成

- PRタイトル
- 変更概要
- 変更理由
- テスト内容
- レビュー観点

4. PR本文を `outputs/pr/xxxx.md` に保存

```md
## Summary

(summary of changes)

## Changes

- xxx
- xxx

## Testing

- [ ] Unit Test
- [ ] Manual Test

## Notes

(any additional notes)
```

5. ユーザーにPRタイトルとPR本文の内容を確認してもらい、必要に応じて修正する。OKであれば、ユーザーにブランチをpushしてもらうように伝える

6. ユーザーのpushが確認できたら、ghコマンドでPR作成

```bash
gh pr create \
  --title "<生成したタイトル>" \
  --file outputs/pr/xxxx.md
  --base <ベースブランチ>
```

7. 作成されたPR URLを表示

```bash
gh pr view --json url -q .url
```

8. `outputs/pr/xxxx.md` を `outputs/pr/<PR番号>/xxxx.md` に`mv`で移動する

## PRタイトルルール

- [Feature]: 新機能
- [Fix]: バグ修正
- [Refactor]: リファクタリング
- [Docs]: ドキュメント
- [Chore]: 雑務

## 注意事項

- PRはタイトル･本文ともに英語で作成してください
- PRのタイトルは簡潔なもので良いです
- gh認証済みであることを確認
- PR作成前に未コミット変更がないか確認
- タイトルは50文字程度に収める
- 本文はレビューしやすい粒度で記載する
