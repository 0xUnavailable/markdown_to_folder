Below is a **fixed and up-to-date README** that accurately reflects the **current behavior of the tool**, including:

* **Both supported input formats** (`layered` markdown lists **and** `tree` / `tree-like` output)
* The **two-phase tree parsing fix** (tree prefix stripping + filename extraction)
* Correct rules for filenames with `_` and normal names (`text.py`)
* No misleading claims about “nested markdown lists only”

You can **replace your README entirely with this**.

---

# Markdown to Repository Structure Tool

This is a Go command-line tool that creates a directory structure and empty files from a textual description of a repository layout.

It supports **two input formats**:

1. **Layered Markdown lists** (indentation-based)
2. **Tree-style layouts** (e.g. `tree` command output using `├──`, `└──`, `│`)

The tool deterministically converts these formats into real directories and files on disk.

---

## Features

* Supports **layered markdown lists** (`-` / `*` with indentation)
* Supports **tree-style structures** (`├──`, `└──`, `│`)
* Correctly strips tree-drawing characters before file creation
* Handles files that:

  * start with `_` (e.g. `__init__.py`)
  * do **not** start with `_` (e.g. `main.go`, `text.py`)
* Creates:

  * **directories** for names without dots
  * **files** for names containing dots
* Continues processing even if individual operations fail
* Prints all created files and directories

---

## Prerequisites

* Go **1.16+**

No external dependencies — standard library only.

---

## Installation

1. Save the program as:

```text
markdown_to_structure.go
```

2. Ensure Go is installed:

```bash
go version
```

---

## Usage

### Basic Usage

```bash
go run markdown_to_structure.go -input structure.md
```

By default, the tool assumes **tree format**.

---

### Command-Line Flags

| Flag      | Description                                         |
| --------- | --------------------------------------------------- |
| `-input`  | Path to input file (default: `structure.md`)        |
| `-format` | Input format: `tree` or `layered` (default: `tree`) |

#### Examples

```bash
# Tree format (default)
go run markdown_to_structure.go -input structure.md

# Explicit tree format
go run markdown_to_structure.go -input structure.md -format tree

# Layered markdown format
go run markdown_to_structure.go -input structure.md -format layered
```

---

## Supported Input Formats

---

## 1️⃣ Tree Format (Recommended)

This format matches the output of the Unix `tree` command.

### Rules

* First non-empty line = **root directory**
* Tree drawing characters are allowed:

  * `├──`
  * `└──`
  * `│`
* Tree characters are **automatically stripped**
* Filenames are extracted safely and deterministically

### Example Input

```text
tx_batch_system/
├── __init__.py
├── types.py
├── journal.py
├── queue.py
├── execution.py
├── batch.py
├── persistence.py
├── verification.py
└── query.py
```

### Output

```text
tx_batch_system/
├── __init__.py
├── types.py
├── journal.py
├── queue.py
├── execution.py
├── batch.py
├── persistence.py
├── verification.py
└── query.py
```

✔ Tree characters never appear in filenames
✔ Works for `_files` and normal filenames
✔ Depth is preserved

---

## 2️⃣ Layered Markdown Format

### Rules (STRICT)

* First non-empty line = **root directory**
* Use `- ` or `* ` for list items
* **Exactly 2 spaces per indentation level**
* **Spaces only** — tabs are not supported

### Files vs Directories

* **Files** → names containing `.`
  (`main.go`, `.env`, `README.md`)
* **Directories** → names without `.`
  (`src`, `utils`, `internal`)

---

### Example Input

```markdown
  - src
    - bot
      - index.ts
      - commands
      - middlewares
    - services
      - wallet.service.ts
      - aave.service.ts
      - moonpay.service.ts
      - transak.service.ts
      - gasless.service.ts
    - db
      - models
      - index.ts
    - config
      - env.ts
      - constants.ts
    - utils
      - logger.ts
      - error.ts
    - app.ts
  - package.json
  - .env.example
  - README.md
```

---

## Common Issues & Fixes

---

### Tree characters appear in filenames

**Cause**: Old versions did not strip tree prefixes.

**Fix**:
Use the updated version — tree prefixes are now stripped **before** filename parsing.

---

### Files created at wrong depth (layered format)

**Cause**: Incorrect indentation.

**Fix**:
Indentation must be **exactly 2 spaces per level**.

```markdown
# WRONG
   - file.ts   # 3 spaces

# CORRECT
    - file.ts  # 4 spaces
```

---

### Files created in current directory instead of root

**Cause**: Missing root directory line.

**Fix**:
Ensure the first line is a directory name with **no indentation or markers**.

---

## Design Notes (Important)

* Tree parsing is **two-phase**:

  1. Strip tree-drawing characters
  2. Extract the filename
* This prevents Unicode tree symbols from leaking into filenames
* Filename grammar supports:

  * `_file.py`
  * `text.py`
  * `README.md`
  * `.env`

---

## Limitations

* Files are created **empty**
* No content templating
* No validation of markdown correctness
* Layered format requires strict spacing
* Does not parse markdown headers or code blocks

---

## Tips

* Prefer **tree format** if possible — it is more robust
* Generate input using:

  ```bash
  tree > structure.md
  ```
* Enable “show whitespace” in your editor when using layered format
* Start small, then expand

---


## License

MIT License
