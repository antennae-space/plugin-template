# Plugin template of Web IDE

Start an external plugin of [Web IDE](https://github.com/antennae-space/web-ide) here: a tool, a tab view, a command, an assistant tool and a backend, each small, built on the [plugin SDK](https://github.com/antennae-space/plugin-sdk). It lists the entries at the root of the project, as a list or a table, opens one in a tab, and summarizes the project for you and for the AI.

## Start

1. *Use this template* on GitHub (or `gh repo create you/plugin-mine --template antennae-space/plugin-template --public --clone`).
2. Rename it: `./scripts/rename.sh mine "Mine" you/plugin-mine` (its id, its name, its repository). The id is the prefix of everything it owns: methods `mine.*`, events `mine.*`, the tool `mine_summary`.
3. `npm ci && npm run build` (the page module), `go test ./...` (the backend).
4. Change it: the IDE's how-to [Write a plugin](https://github.com/antennae-space/web-ide/blob/main/docs/howto/write-a-plugin.md) explains every part.

On Linux with GNU sed; npm 10 needs the `.npmrc` here (`legacy-peer-deps`, for the optional peers of Vite 8).

## What is in it

| Path | What |
|---|---|
| `plugin.json` | The manifest (checked by the SDK's JSON Schema): id, version, the IDE versions it runs on, the address of its archives, its backend, its permissions, its tool, tab view, command and assistant tool |
| `web/src/index.tsx` | The page module: what it registers, by the ids of `plugin.json` |
| `web/src/Tool.tsx` | The tool: the entries from the backend, a list or a table (the SDK's kit), reloaded on the backend's event |
| `web/src/Entry.tsx` | The tab view: an entry, and a note kept with the tab |
| `web/src/example.css` | Its own styles, on the theme's variables |
| `i18n/fr.json`, `i18n/es.json` | Its texts in French and Spanish; English is the key (`ide.t`) |
| `backend/plugin.go` | The backend: its methods, its event, its assistant tool, through the IDE's `fs.list` |
| `backend/plugin_test.go` | Its tests, the IDE played by the test |
| `scripts/pack.sh` | Builds `dist/<id>-linux-amd64.tar.gz` and `dist/<id>-linux-arm64.tar.gz` |
| `scripts/rename.sh` | Makes the template yours, once |
| `.github/workflows/` | `ci.yml` builds and tests each push; `release.yml` publishes a tag `v*` with its archives |

## Permissions

It asks `files:read` (its backend may read the project) and `core:fs.list` (it lists the project through the IDE, which also works for a project over SSH). Ask for the least your plugin needs, and handle a refusal: the user grants each one, and may revoke it.

## Release

1. Set `version` in `plugin.json` (`0.2.0`), commit.
2. Tag `v0.2.0` and push the tag: the workflow builds and publishes the archives.
3. In the IDE: *Settings → Plugins → Install a plugin*, the address of the repository.

Never move a tag once released: the IDE keeps the hash of each archive it installed and refuses another one under the same version. To try a version first, release a pre-release (`v0.2.0-rc1`, with that version in `plugin.json`) and install it with the constraint `=0.2.0-rc1`.

## License

MIT: yours to change, for your own plugin under any license.
