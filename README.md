# swagen-v2

## 1. Overview

swagen-v2 is a CLI that lets you define Swagger (OpenAPI) schemas right from your terminal. It follows Swagger’s standard properties and helps you create models, request/response schemas, and API path definitions with a smooth, interactive flow.

## 2. Why use this CLI?

- Generate Swagger schemas without manually editing files.
- Eliminate typos in $ref relative paths — the tool resolves and inserts correct paths for you.

## 3. Installation

Install via `go install` from the GitHub repository (binaries are distributed as Go code for now):

```bash
go install github.com/Daaaai0809/swagen-v2@latest
```

Repository: https://github.com/Daaaai0809/swagen-v2

Make sure your `GOBIN` is on your PATH so the `swagen-v2` command is available.

## 4. Before you start

Create a `.env` file at the root of the directory where you’ll use swagen-v2, based on `.env.example`.

- `SWAGEN_MODEL_PATH`: Directory where model schemas are generated.
- `SWAGEN_SCHEMA_PATH`: Directory where request/response schemas are generated.
- `SWAGEN_API_PATH`: Directory where path (API) schemas are generated.

All three variables are required. Missing any of them may cause the CLI to fail.

## 5. Commands

### 5.1 `swagen-v2 model`
- Generate a model schema.

### 5.2 `swagen-v2 schema`
- Generate request/response schemas.
- Reference model schema properties via `$ref` with interactive directory traversal and field selection.
- Or define properties inline without `$ref`.

### 5.3 `swagen-v2 path`
- Generate API definitions.
- `$ref` referencing is supported from both `model` and `schema`.
- Where `$ref` can be used: `parameters`, `requestBody`, and `responses.[status].content.[mediaType].schema`.
- You can also define these inline without `$ref`.

## 6. Bugs and suggestions

- Please open an issue in this repository.
- Or DM https://x.com/big_doge_
