# Konfetti 🎉

Your config archaeology sidekick. Point it at a directory (or pipe stuff in) and it flattens JSON / YAML / XML / key=value blobs so you can actually see what you're running in prod (or blindly copied from Stack Overflow in 2014). It doesn’t judge your configs. Much.

## Features (now, not aspirational)
* Scan current dir or explicit path (falls back to OS defaults only if CWD somehow explodes)
* Auto-detect & parse: JSON, YAML, XML, .conf/.ini/.properties/.txt key=value, raw text fallback
* Case-insensitive filtering: filename (-filter), key (-key), value (-value)
* Flatten nested structures (dot notation)
* Output: text (default), json, table
* Profiles & defaults via `~/.konfetti.yaml`
* Pipe stdin into `scan` or `explain`
* Warnings for unreadable paths (silence with `-no-warn`)

## Install / Run
```bash
go build -o konfetti   # build binary
./konfetti scan -key debug
```
Or just: `go run main.go scan -key debug`

Cross-compile example: `GOOS=linux GOARCH=amd64 go build -o konfetti-linux`.

## Commands
* `scan`    – find & parse config files
* `explain` – very lightweight heuristic summary (stdin or file)
* `explore` – placeholder for future TUI
* `tips`    – show a small helpful nudge

## Core Flags
| Flag | Purpose |
|------|---------|
| `-path` | Directory to scan (single or comma-separated list) |
| `-filter` | Filename substring filter |
| `-key` | Match setting key (case-insensitive substring) |
| `-value` | Match setting value (case-insensitive substring) |
| `-output` | `text` (default) | `json` | `table` |
| `-no-warn` | Suppress unreadable path warnings |
| `-profile` | Use profile from config file |
| `-interactive` | Prompt before scanning when path empty |

Defaults if no path given: scan current working directory. If that fails, OS fallbacks:
* macOS / Linux: `/etc`, `~/.config` (+ `~/Library/Application Support` on macOS)
* Windows: `%ProgramData%`, `%APPDATA%`, `%LOCALAPPDATA%`

## Stdin Tricks
```bash
cat app.yaml | ./konfetti scan -key log
curl -s https://example.com/config.json | ./konfetti explain
echo -e 'MODE=prod\nLOG_LEVEL=info' | ./konfetti scan -value prod -output table
```

## Profiles & Config
`~/.konfetti.yaml` structure:
```yaml
defaults:
  output: text
profiles:
  debug:
    key: debug
    output: table
  prod:
    value: prod
    output: json
```
Precedence: defaults < profile < CLI flag.

## Example Output (text)
```
File: /etc/app/config.yaml [yaml]
  server.port = 8080
  logging.level = error
  logging.level = warn
  logging.level = info
---
File: ./app/settings.json [json]
  database.host = db.example.com
  database.port = 5432
  feature.enabled = true
---
```

## FAQ (micro-dose)
* 0 matches? Loosen filters or point at a richer path (`-path ~/.config`).
* Multiple paths? `-path "./cfg,/etc,/opt/app"`.
* Too noisy? Add `-no-warn`.
* Want structured scripting? `-output json`.
* Stdin parse failed? It probably wasn’t valid JSON/YAML/XML; fell back to raw key=value or plain text.
* Windows Defender flagged the exe? New unsigned Go binaries sometimes trigger generic ML detections (e.g. `Win32/Sabsik.FL.A!ml`). See `AV-SUBMISSION.md` for false positive submission steps or use a signed release when available.

## Roadmap (short list)
* Real explain engine
* TUI explorer
* Sensitive key detector

## License
MIT (unless entropy wins). Built with coffee & sarcasm.

---
Enjoy responsibly. Spray confetti, not secrets.