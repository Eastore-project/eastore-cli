# Eastore CLI

Eastore CLI is a command-line implementation that allows you to directly interact with the Eastore protocol trustlessly from your local device. It provides utilities for making storage deals, encrypting/decrypting files, and building on top of the protocol.

## Installation

```bash
go install github.com/eastore-project/eastore-cli@latest
```

## Configuration

The CLI can be configured either through command-line flags or environment variables:

- `PRIVATE_KEY` - Private key for signing transactions (required)
- `RPC_URL` - RPC URL for the network
- `EASTORE_CONTRACT_ADDRESS` - Address of the Eastore contract

## Commands

### version
Displays the current version of the Eastore CLI.

```bash
eastore version
```

### make-deal
Submit a new storage deal proposal. Supports optional file encryption and various deal parameters.

```bash
eastore make-deal --input <file-path> [options]
```

Key options:
- `--input` - Input file or folder path (required)
- `--outdir` - Output directory for CAR files (uses temp dir if not provided)
- `--duration` - Duration of the deal in epochs (default: 518400)
- `--encrypted` - Whether to encrypt the file before making the deal (default: false)
- `--encrypted-out-dir` - Output directory for encrypted files (uses temp dir if not provided)
- `--verified-deal` - Whether to use verified client data-cap (default: true)

Advanced options:
- `--buffer-type` - Buffer type: "lighthouse" or "local" (default: local)
- `--buffer-api-key` - API key for buffer service
- `--buffer-url` - Base URL for buffer service
- `--start-epoch-offset` - Offset from current chain head for deal start (default: 1000)
- `--start-epoch` - Explicit start epoch (overrides offset)
- `--storage-price` - Price in attoFIL per epoch per GiB (default: 0)
- `--provider-collateral` - Provider's collateral in attoFIL (default: 0)
- `--client-collateral` - Client's collateral in attoFIL (default: 0)
- `--skip-ipni` - Skip announcing deal to Network Indexer (default: false)
- `--remove-unsealed` - Don't keep unsealed copy for fast retrieval (default: false)

### encrypt
Encrypt a file using AES with a key derived from your wallet signature.It will give you key with which you can decrypt the file.

```bash
eastore encrypt --input <file-path> [--out-dir <directory>]
```

### decrypt
Decrypt a file that was previously encrypted using the encrypt command.

```bash
eastore decrypt --input <file-path> --key <hex-key> [--out-dir <directory>]
```

## Building from Source

```bash
git clone https://github.com/eastore-project/eastore-cli
cd eastore-cli
go build ./cmd/eastore
```
