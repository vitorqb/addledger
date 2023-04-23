# HledgerAdd

A helper tool to enter HLedger journal entries.

## Development

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

For running linting and formatting:

```
just format lint
```

## Configuration

All configuration variables can be set with command line flags or
environmental variables. Env vars must be prefixed with
`ADDLEDGER_`. For example: `--destfile=foo` is the same as `export
ADDLEDGER_DESTFILE=foo`.

To see all see the output of `--help`.
