# HledgerAdd

A helper tool to enter HLedger journal entries.

## Development

### Setup

You will need to install and enabled https://github.com/asdf-vm/asdf.

Then for setting up development:

```
just setup
```

**REMEMBER** that you need to source the `.env` file for configuring
the app during development:

```
source .env
```

### Testing

```
# All
just test

# A single file
just test internal/utils/utils_test.go 
```

### Linting and Formatting

For running linting and formatting:

```
just format lint
```

### Useful commands

```
# Gets the transactions as a json and pipe to jq
hledger print --output-format=json | jq .
```

## Configuration

All configuration variables can be set with command line flags or
environmental variables. Env vars must be prefixed with
`ADDLEDGER_`. For example: `--destfile=foo` is the same as `export
ADDLEDGER_DESTFILE=foo`.

To see all see the output of `--help`.
