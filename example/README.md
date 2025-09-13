## swagen-v2 Example

### 1. Running with Makefile

You can build and prepare the swagen-v2 example environment in one shot using `make` at the repository root.

- `make` : Performs the full initial setup in a single command.
- `make build-example` : Builds the example environment (creates `example/swagen-v2.local`).
- `make build-binary` : Rebuilds `example/swagen-v2.local` after you change the source code.

If you keep editing the code, just run `make build-binary` again to refresh the local binary.

### 2. Running manually

1. Install swagen-v2 (see the root README Installation section):
	- ../README.md#3-installation
2. Copy the environment file template:
	```bash
	cp .env.example .env
	```
	Adjust `SWAGEN_MODEL_PATH`, `SWAGEN_SCHEMA_PATH`, and `SWAGEN_API_PATH` as needed.
3. Go back to the repository root (if you are not already there) and build the binary:
	```bash
	go build -o ./example/swagen-v2.local
	```

That’s it — the environment is ready.

To run swagen-v2 inside the example environment, move into the `example` directory and execute:

```bash
./swagen-v2 <command>
```

If you prefer using the explicitly built file name, you can also run:

```bash
./swagen-v2.local <command>
```

