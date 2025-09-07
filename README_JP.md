# swagen-v2

## 1. 概要

swagen-v2 は、ターミナル上で対話的に Swagger（OpenAPI）スキーマを定義できる CLI アプリです。Swagger の標準プロパティに準拠し、モデル、リクエスト／レスポンスのスキーマ、API パス定義をスムーズに作成できます。

## 2. この CLI を導入することによってもたらされる価値

- ファイルを直接編集せずに、対話的な操作だけで Swagger スキーマを生成できます。
- `$ref` の相対パスを自動で解決・挿入するため、タイポの心配が減ります。

## 3. Install 手順

現状は GitHub 上で Go コードのバイナリを配布しているため、`go install` でインストールします。

```bash
go install github.com/DAAAai0809/swagen-v2@latest
```

リポジトリ: https://github.com/Daaaai0809/swagen-v2

`GOBIN` が PATH に通っていることを確認し、`swagen-v2` コマンドを実行できるようにしてください。

## 4. 事前準備

`swagen-v2` を使用するプロジェクトのルートディレクトリに、`.env.example` を参考に `.env` を作成してください。

- `MODEL_PATH`: モデルスキーマを生成するディレクトリ
- `SCHEMA_PATH`: request/response スキーマを生成するディレクトリ
- `API_PATH`: path スキーマを生成するディレクトリ

いずれか 1 つでも欠けると、正常に動作しません。

## 5. 各コマンドの使い方

### 5.1 `swagen-v2 model`
- モデルスキーマ生成コマンド

### 5.2 `swagen-v2 schema`
- リクエスト／レスポンスのスキーマ生成コマンド
- `$ref` により model スキーマのプロパティを参照可能
- `$ref` を使用せず、その場でプロパティを定義することも可能

### 5.3 `swagen-v2 path`
- API 定義（エンドポイント）生成コマンド
- `$ref` による `model`／`schema` からの参照が可能
- 参照は `parameters`, `requestBody`, `responses.[status].content.[mediaType].schema` で使用可能
- `$ref` を使用しない場合は、その場で定義することも可能

## 6. バグや提案など

- このリポジトリに Issue を作成してください。
- または https://x.com/big_doge_ まで DM をお送りください。

