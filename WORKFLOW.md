# Development Workflow

## The Loop

1. **Read the phase** — understand the goal, the concepts, and what needs to be in place
2. **Learn the concepts** — ask Claude to explain anything in the "Understand first" section before writing code
3. **Implement** — write the code and tests yourself
4. **Review** — paste your code and say `review phase <N>` — Claude will check it, explain what's wrong, and suggest what to improve
5. **Fix** — make the changes yourself based on the feedback
6. **Move on** — only when tests pass and Claude has signed off

---

## Rules

- No moving to the next phase until the current one is done and reviewed
- Tests are part of done — a phase is not complete without them
- Ask Claude to explain concepts, not to write code
- If stuck, ask for a nudge — not the answer

---

## Who does what

| Task | Who |
|------|-----|
| Implementation code | Link |
| Tests | Link |
| Code review and feedback | Claude |
| Teaching concepts | Claude |
| Documentation and comments | Claude (can edit directly) |

---

## Phrases to use

| Say this | What happens |
|----------|-------------|
| `starting phase N` | Claude walks you through the "Understand first" concepts |
| `review phase N` | Claude checks your code and tests, gives feedback |
| `nudge` | Claude gives a hint — not the answer |

---

## Build Plan

### Phase 0 — Module setup and project wiring
**Goal:** A runnable binary that responds to `term --help` before any real logic is written.

**Understand first:**
- What a Go module is and what `go.mod` tracks
- How `go get` adds a dependency and what `go.sum` is for
- The minimum structure of a cobra root command

**What needs to be in place:**
- A Go module initialised at the repo root
- Cobra added as a dependency
- A `main.go` with a root cobra command that can be run and responds to `--help`

**Tests:**
- No unit tests this phase — the verify step is sufficient
- `go build ./src/` and `go vet ./src/...` must pass clean

**Verify:** `go run src/main.go --help` shows a cobra help page with your description

---

### Phase 1 — Data model
**Goal:** Define the `Term` and `Glossary` types in the store package, no file I/O yet.

**Understand first:**
- How Go packages work and how one package imports another using the module path
- What struct field tags (`json:"word"`) do and why they matter for encoding
- The difference between exported and unexported identifiers in Go

**What needs to be in place:**
- `src/store/store.go` with the `Term` and `Glossary` type definitions
- The store package importable from `main.go` without errors

**Tests:**
- A test that instantiates a `Term` and a `Glossary` and confirms the types are correct
- Confirm JSON tags round-trip correctly: marshal a `Term` to JSON and check the keys are lowercase

**Verify:** `go build ./src/` and `go test ./...` pass with no errors

---

### Phase 2 — Store: file I/O
**Goal:** Persist and retrieve a `Glossary` from disk at a path that is always passed in — never hardcoded inside the function.

**Understand first:**
- How to resolve the user's home directory at runtime
- Reading and writing whole files with the standard library
- `json.Unmarshal` and `json.MarshalIndent`
- How to tell the difference between "file doesn't exist yet" and a real I/O error
- Why functions that accept a `path` parameter are easier to test than functions that decide the path themselves

**What needs to be in place:**
- A `DefaultPath()` function that resolves the production storage path
- A `Load(path string)` function — a missing file is not an error, just an empty glossary
- A `Save(path string, g Glossary)` function — creates the directory if needed, writes indented JSON

**Tests:**
- Load from a path that doesn't exist → empty glossary, no error
- Save a glossary, then load it back → contents match exactly
- Save to a path inside a directory that doesn't exist yet → directory and file are both created

**Verify:** After a `Save`, the file exists on disk with valid, human-readable JSON

---

### Phase 3 — `term add`
**Goal:** Add a new term to the glossary via command-line flags.

**Understand first:**
- How cobra subcommands are structured and registered
- The difference between `Run` and `RunE` on a cobra command
- How to bind flags to variables in cobra
- How to write to stderr in Go (errors go there, normal output goes to stdout)

**What needs to be in place:**
- A `root.go` in `src/cmd/` that owns and exports the root command
- An `add.go` command with flags for word, definition, analogy, and tags
- Validation that required fields are present before anything is saved
- `main.go` updated to delegate to the cmd package

**Tests:**
- Adding a valid term results in it appearing in the saved glossary
- Missing a required flag produces an error and does not modify the file
- Tags are stored as a slice, not a single string

**Verify:** `term add --word goroutine --definition "lightweight thread" --tags concurrency,go` adds the entry to `~/.term/glossary.json`

---

### Phase 4 — `term list`
**Goal:** Print all stored terms to stdout in a readable format.

**Understand first:**
- How to range over a slice in Go
- How to join a string slice for display
- What the right behaviour is when there's nothing to list

**What needs to be in place:**
- A `list.go` command that loads and prints every term
- Empty glossary handled gracefully (a message, not silence or an error)
- Registered on the root command

**Tests:**
- List with a populated glossary prints all terms
- List with an empty glossary prints the "nothing here" message without erroring

**Verify:** After adding terms with Phase 3, `term list` shows each one with all fields

---

### Phase 5 — `term search`
**Goal:** Return terms whose word or tags contain the query string, regardless of case.

**Understand first:**
- How cobra handles positional arguments (as opposed to flags)
- How to do case-insensitive substring matching with the standard library
- How to search inside a nested slice (tags inside a term)

**What needs to be in place:**
- A `search.go` command that takes one positional argument
- Matches on word or any tag using a case-insensitive substring check
- Prints matches or a "no results" message

**Tests:**
- Searching by word (partial, uppercase) returns the correct term
- Searching by tag returns the correct term
- A query with no matches returns the "no results" message without erroring
- Running `search` with no argument fails with a usage error

**Verify:** `term search GOROUT` returns any term whose word contains "gorout" (case-insensitive)

---

### Phase 6 — `term delete`
**Goal:** Remove a term by word, matched case-insensitively.

**Understand first:**
- The idiomatic Go way to remove an element from a slice (hint: think filtering, not splicing)
- How to track whether anything was actually matched and removed
- When to save vs when not to (only on a successful match)

**What needs to be in place:**
- A `delete.go` command that takes one positional argument
- Case-insensitive word match
- Prints a confirmation on success; an informative error to stderr if nothing matched

**Tests:**
- Deleting an existing term removes it from the saved glossary
- Deleting with a different case still matches (e.g. `DELETE GOROUTINE`)
- Deleting a term that doesn't exist returns a non-zero exit and an error message
- The glossary file is not rewritten if nothing was removed

**Verify:** After `term delete goroutine`, `term list` no longer shows that entry

---

### Phase 7 — `term random`
**Goal:** Return one randomly selected term for self-quizzing.

**Understand first:**
- How to generate a random integer within a range in Go
- What happens if you try to get a random element from an empty slice (and how to prevent it)
- Whether your Go version needs explicit seeding for `math/rand` (check `go.mod`)

**What needs to be in place:**
- A `random.go` command with no arguments
- An empty-glossary guard before any random selection
- Output in the same format as list and search

**Tests:**
- Random on a single-term glossary always returns that term
- Random on an empty glossary prints the "nothing here" message without panicking

**Verify:** Running `term random` several times with multiple terms in the glossary returns different results

---

### Phase 8 — Polish
**Goal:** Make the CLI feel complete — consistent output, full help text, and no rough edges.

**What needs to be in place:**
- `Short`, `Long`, and `Example` fields populated on every command
- A shared `printTerm` helper so list, search, and random all render terms identically
- A duplicate-word guard in `add` (same word already exists → warn, don't silently overwrite)
- All errors from Load/Save surfaced with enough context to understand what went wrong
- Install instructions added to `docs/`

**Tests:**
- Review all existing tests — do they cover the new duplicate guard?
- Add a test: adding a duplicate word does not result in two entries

**Verify:** Every `--help` page reads clearly; all five commands handle every edge case you can think of without panicking

---

## File map

| Phase | Files |
|-------|-------|
| 0 | `go.mod`, `go.sum`, `src/main.go` |
| 1 | `src/store/store.go` (types), `tests/store_test.go` |
| 2 | `src/store/store.go` (+ Load, Save, DefaultPath), `tests/store_test.go` |
| 3 | `src/cmd/root.go`, `src/cmd/add.go`, `src/main.go`, `tests/add_test.go` |
| 4 | `src/cmd/list.go`, `tests/list_test.go` |
| 5 | `src/cmd/search.go`, `tests/search_test.go` |
| 6 | `src/cmd/delete.go`, `tests/delete_test.go` |
| 7 | `src/cmd/random.go`, `tests/random_test.go` |
| 8 | `src/cmd/format.go`, edits to all cmd files, `docs/install.md` |
