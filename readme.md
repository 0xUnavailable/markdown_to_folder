# Markdown to Repository Structure Tool

This is a Go command-line tool that creates a directory structure and files based on a markdown file specifying the repository layout. It parses a nested markdown list to generate folders (items without extensions) and empty files (items with extensions).

## Features

- Reads a markdown file with a nested list format to define the repository structure
- Creates directories and empty files according to the specified hierarchy
- Supports custom input markdown files via a command-line flag
- Handles errors gracefully, continuing processing even if some operations fail
- Provides detailed output showing which files and directories are created

## Prerequisites

- Go (version 1.16 or higher recommended) installed on your system

## Installation

1. Save the Go program as `markdown_to_structure.go`
2. Ensure you have Go installed and configured
3. No additional dependencies are required, as the tool uses the Go standard library

## Usage

### Basic Usage

1. Create a markdown file (e.g., `structure.md`) that defines the repository structure
2. Run the tool using the following command:

```bash
go run markdown_to_structure.go -input structure.md
```

3. The tool will create the directory structure and files in the current working directory

### Command-Line Flags

- `-input`: Path to the markdown file containing the repository structure (default: `structure.md`)

```bash
# Use default structure.md file
go run markdown_to_structure.go

# Use custom markdown file
go run markdown_to_structure.go -input my_structure.md
```

## Markdown File Format

### Important Rules

The markdown file **must** follow these specific rules for the tool to work correctly:

#### 1. Root Directory
- The **first non-empty line** specifies the root directory name
- This line should **not** have any leading spaces or list markers (`-` or `*`)
- Example: `my_project`

#### 2. List Items
- Use `-` (dash) or `*` (asterisk) followed by a **single space** to denote list items
- Format: `- item_name` or `* item_name`

#### 3. Indentation (CRITICAL)
- **Use exactly 2 spaces per indentation level**
- **Do NOT use tabs** - only spaces are supported
- Each nested level must be indented exactly 2 more spaces than its parent
- Inconsistent indentation will cause files/folders to be created in the wrong location

#### 4. Files vs Directories
- **Files**: Items containing a dot/period (`.`) in the name (e.g., `index.ts`, `package.json`, `.env.example`)
- **Directories**: Items without a dot/period in the name (e.g., `src`, `components`, `utils`)

### Example Markdown Structure

Here's a properly formatted `structure.md` file:

```markdown
my_project
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

### Indentation Breakdown

Let's analyze the indentation in the example above:

```
my_project          ← No indentation (0 spaces) - Root directory
  - src             ← 2 spaces - Level 1 (inside my_project)
    - bot           ← 4 spaces - Level 2 (inside src)
      - index.ts    ← 6 spaces - Level 3 (inside bot)
      - commands    ← 6 spaces - Level 3 (inside bot)
    - services      ← 4 spaces - Level 2 (inside src)
      - wallet.service.ts  ← 6 spaces - Level 3 (inside services)
```

### Output Structure

Running the tool with the example above will create:

```
my_project/
├── src/
│   ├── bot/
│   │   ├── index.ts
│   │   ├── commands/
│   │   └── middlewares/
│   ├── services/
│   │   ├── wallet.service.ts
│   │   ├── aave.service.ts
│   │   ├── moonpay.service.ts
│   │   ├── transak.service.ts
│   │   └── gasless.service.ts
│   ├── db/
│   │   ├── models/
│   │   └── index.ts
│   ├── config/
│   │   ├── env.ts
│   │   └── constants.ts
│   ├── utils/
│   │   ├── logger.ts
│   │   └── error.ts
│   └── app.ts
├── package.json
├── .env.example
└── README.md
```

## Common Issues and Troubleshooting

### Files Created in Wrong Directory

**Problem**: Files appear in the wrong folder or at the wrong nesting level.

**Solution**: Check your indentation carefully. Each level must be exactly 2 spaces more than its parent.

```markdown
# ❌ WRONG - Inconsistent spacing
my_project
  - src
   - file.ts          ← 3 spaces (should be 4)
     - another.ts     ← 5 spaces (should be 6)

# ✅ CORRECT - Consistent 2-space indentation
my_project
  - src
    - file.ts         ← 4 spaces
      - nested
        - another.ts  ← 8 spaces
```

### Mixed Tabs and Spaces

**Problem**: Tool doesn't recognize the structure properly.

**Solution**: Ensure you're using only spaces, not tabs. Configure your text editor to insert spaces when you press Tab.

### Missing Root Directory

**Problem**: Tool creates files in the current directory instead of a new folder.

**Solution**: Ensure the first line is the root directory name with no indentation or list markers.

## Advanced Examples

### Web Application Structure

```markdown
webapp
  - public
    - index.html
    - favicon.ico
  - src
    - components
      - Header.jsx
      - Footer.jsx
    - pages
      - Home.jsx
      - About.jsx
    - styles
      - main.css
    - App.jsx
    - index.js
  - package.json
  - .gitignore
```

### Go Project Structure

```markdown
go-api
  - cmd
    - api
      - main.go
  - internal
    - handlers
      - user.go
      - auth.go
    - models
      - user.go
    - database
      - db.go
  - pkg
    - utils
      - validator.go
  - go.mod
  - go.sum
  - README.md
```

## Tips for Creating Your Structure File

1. **Plan your structure first**: Sketch out your directory hierarchy before writing the markdown
2. **Use a code editor**: Use editors like VS Code, Sublime Text, or Vim that show space/tab characters
3. **Enable visible whitespace**: Turn on "Show Whitespace" in your editor to see indentation clearly
4. **Be consistent**: Stick to 2 spaces per level throughout the entire file
5. **Test incrementally**: Start with a small structure to verify it works, then expand
6. **Validate indentation**: Count spaces carefully for each level (2, 4, 6, 8, etc.)

## Notes

- The tool creates empty files; you need to populate them with content manually
- Indentation **must** be consistent (exactly 2 spaces per level)
- Errors during file or directory creation are logged, but the tool continues processing
- Hidden files (starting with `.`) are supported (e.g., `.env`, `.gitignore`)

## Limitations

- The tool does not populate files with content
- Only supports markdown lists with `-` or `*` prefixes
- Requires exactly 2 spaces per indentation level
- Does not support mixing tabs and spaces
- Cannot parse other markdown formats (code blocks, headers, etc.)

## Contributing

Feel free to submit issues or pull requests to improve the tool. Suggestions for features like:
- Adding default file content templates
- Supporting other markdown formats
- Auto-detecting indentation style
- Validating markdown structure before processing


## License

This project is licensed under the MIT License.