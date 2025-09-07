## swagen-v2

### 1. Overview

swagen-v2 is a CLI tool that helps you define Swagger (OpenAPI) schemas interactively in your terminal. It adheres to standard Swagger properties and guides you to build models, request/response schemas, and API path definitions with minimal friction.

### 2. Why use this CLI?

- Generate Swagger schemas without manual editing — focus on selection/inputs, not boilerplate.
- Reduce typos in $ref relative paths — the tool resolves and inserts correct relative paths for you.

### 3. Installation

Install via `go install` from the GitHub repository (binaries are distributed as Go code for now):

```bash
go install github.com/Daaaai0809/swagen-v2@latest
```

Repository: https://github.com/Daaaai0809/swagen-v2

Ensure your `GOBIN` is on PATH so the `swagen-v2` command is available.

### 4. Usage

#### 4.1 `swagen-v2 model`
- Generate a model schema file.

#### 4.2 `swagen-v2 schema`
- Generate request/response schema files.
- You can reference model schema properties via `$ref` with interactive directory traversal and field selection.
- Alternatively, you can define properties inline without `$ref`.

#### 4.3 `swagen-v2 path`
- Generate API path definitions (endpoint specs).
- `$ref` referencing is supported from both `model` and `schema` directories.
- Where `$ref` is available: `parameters`, `requestBody`, and `responses.[status].content.[mediaType].schema`.
- You can also define these inline without `$ref`.

### 5. Bugs and Suggestions

- Please open an issue in this repository.
- Or send a DM to https://x.com/big_doge_
