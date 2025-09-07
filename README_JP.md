# swagen-v2

## 1. 概要

swagen-v2 は、ターミナル上で対話的に Swagger（OpenAPI）スキーマを定義できる CLI ツールです。Swagger の標準プロパティに準拠し、モデル、リクエスト／レスポンススキーマ、API パス定義をスムーズに作成できます。

## 2. この CLI を導入する価値

- 手を動かしてファイルを編集せずに Swagger スキーマを生成できます。
- $ref の相対パスを自動解決・自動挿入し、タイポを防げます。

## 3. インストール手順

現状は GitHub で Go コードを配布しているため、`go install` を利用してインストールします。

```bash
go install github.com/Daaaai0809/swagen-v2@latest
```

リポジトリ: https://github.com/Daaaai0809/swagen-v2

`GOBIN` が PATH に通っていることを確認し、`swagen-v2` コマンドを実行できるようにしてください。

## 4. 各コマンドの使い方

### 4.1 `swagen-v2 model`
- モデルスキーマ生成コマンド

### 4.2 `swagen-v2 schema`
- リクエスト／レスポンススキーマ生成コマンド
- `$ref` により model スキーマのプロパティを参照可能
- `$ref` を使わず、その場でプロパティを定義することも可能

### 4.3 `swagen-v2 path`
- API 定義（エンドポイント）生成コマンド
- `$ref` による `model`／`schema` からの参照が可能
- 参照は `parameters`, `requestBody`, `responses.[status].content.[mediaType].schema` で使用可能
- `$ref` を使用しない場合はその場で定義することも可能

## 5. バグや提案など

- 本リポジトリに Issue を作成してください。
- もしくは https://x.com/big_doge_ まで DM をお送りください。

