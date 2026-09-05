<!-- markdownlint-disable -->
# Example project files

These are ready-to-run HexTest project files. Run one with:

```sh
go run ./cmd/hextest execute examples/<file>.json
```

| File | Target | What it shows |
| --- | --- | --- |
| `demo-project.json` | bundled demo server (`:3443`) | Full project: account suite, status/time/size asserts, schema and regex checks. Run `go run ./cmd/exampleserver` first. |
| `editable-demo.json` | bundled demo server (`:3443`) | Smaller project focused on header and body-type assertions. Good starting point for editing in the web UI (`hextest front`). |
| `single-get.json` | `http://localhost:8080` | Minimal single GET request with `containsKey` / `containsKey -R` checks. Point it at any API by changing `Url`. |

## `legacy/`

Older hand-made snapshots kept for reference. They target `http://localhost:8080/base`, which no server in this repo provides, so they are examples of the on-disk format rather than something to run as-is.
