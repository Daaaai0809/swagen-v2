# swagen-v2 サンプル

### 1. 動作手順（Makefileを使用する場合）

リポジトリ直下で `make` を実行するだけで、swagen-v2 の動作環境を一度に構築できます。

- `make` : 初期セットアップをまとめて実行します。
- `make build-example` : `example` ディレクトリ内で動作させるための `swagen-v2.local` バイナリを生成します。
- `make build-binary` : コードを変更したあとに再度 `example/swagen-v2.local` を生成します。

コードを更新したら `make build-binary` を再実行するだけで反映できます。

### 2. 動作手順（手動の場合）

1. swagen-v2 のインストール方法は `../README_JP.md` の「3. Install 手順」を参照してください。
2. `.env.example` をコピーして `.env` を作成:
   ```bash
   cp .env.example .env
   ```
   必要に応じて `MODEL_PATH`, `SCHEMA_PATH`, `API_PATH` を編集します。
3. リポジトリのルートディレクトリで以下を実行:
   ```bash
   go build -o ./example/swagen-v2.local
   ```

これで環境構築は終わりです。`example` ディレクトリに移動し、以下を実行すると swagen-v2 を動かせます。

```bash
./swagen-v2 <command>
```

もしくはビルドされたローカルバイナリ名で実行:

```bash
./swagen-v2.local <command>
```
