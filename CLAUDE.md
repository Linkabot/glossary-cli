# glossary-cli (term)

## Project Overview
A personal CLI tool called `term` for storing and retrieving CS/engineering terms.
The user adds their own definitions, analogies, and tags to terms they encounter at work
or while reading — helping them decode jargon faster and communicate more confidently.

## Tech Stack
- **Language:** Go
- **CLI Framework:** cobra (command routing)
- **Storage:** Local JSON file (`~/.term/glossary.json`)
- **No external dependencies** beyond cobra and the standard library

## Data Model

```go
type Term struct {
    Word       string   `json:"word"`
    Definition string   `json:"definition"`
    Analogy    string   `json:"analogy"`
    Tags       []string `json:"tags"`
}
```

## Commands
| Command | Description |
|---------|-------------|
| `term add` | Add a new term interactively |
| `term list` | List all stored terms |
| `term search <query>` | Search terms by word or tag |
| `term delete <word>` | Remove a term by name |
| `term random` | Return a random term (for self-quizzing) |

## Project Structure
```
glossary-cli/
├── src/
│   ├── main.go           — entrypoint
│   ├── cmd/              — cobra command definitions (one file per command)
│   └── store/            — JSON read/write logic
├── tests/                — unit and integration tests
├── scripts/              — build or install helpers
└── docs/                 — usage notes, design decisions
```

## Conventions
- One file per cobra command under `src/cmd/`
- Storage logic lives in `src/store/` and is separate from command logic
- All file I/O goes through the store package — commands never touch the JSON directly
- Errors are printed to stderr; normal output to stdout
- Match on `strings.ToLower` for search and delete (case-insensitive)

## Build Plan
See `WORKFLOW.md` for the full phase-by-phase plan.

## Out of Scope (for now)
- Remote sync or cloud backup
- Import/export
- Edit command (delete + re-add for now)
