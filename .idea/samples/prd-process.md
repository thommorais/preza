# PRD (Product Requirements Document)

## IMPORTANT: Using bd CLI for Task Management

**The PRD workflow now uses the `bd` CLI for task management instead of prd.json.**

**You MUST maintain:**
- `prd/progress.txt` - Progress tracking log (MANDATORY)

**Task management is done via `bd` commands:**
- Create tasks: `bd create --title="..." --type=task|feature --priority=2`
- Find work: `bd ready`
- Claim work: `bd update <id> --status=in_progress`
- Complete work: `bd close <id>`
- Sync: `bd sync`

---

## File: progress.txt

**Location:** `prd/progress.txt`

This is a **mandatory file** that serves as a memory log for completed work.

### When to Update

- **After completing ANY task** - Append a brief summary of what was done
- **When encountering blockers** - Document the issue and why work stopped
- **When encountering errors** - Describe any error that occurred during implementation
- **When making important decisions** - Log architectural or implementation choices

### Initial Structure

```txt
# Progress log

```

---

## bd CLI Workflow

### Creating Tasks

Use `bd create` to create new tasks with implementation steps in the description:

```bash
bd create --title="Add zen mode toggle button" \
  --type=feature \
  --priority=2 \
  --description="Implementation steps:
- Create ZenModeToggle component
- Add toggle state to store
- Import and render in header"
```

**Priority levels:**
- 0 or P0: Critical
- 1 or P1: High
- 2 or P2: Medium (default)
- 3 or P3: Low
- 4 or P4: Backlog

**Types:**
- `task`: General task or bug fix
- `feature`: New functionality or UI component
- `bug`: Bug fix

### Working on Tasks

```bash
# 1. Find available work (no blockers)
bd ready

# 2. View task details
bd show <id>

# 3. Claim the task
bd update <id> --status=in_progress

# 4. [Implement the task]

# 5. Close the task when complete
bd close <id>

# 6. Append summary to prd/progress.txt
echo "$(date): Completed <id> - Brief summary" >> prd/progress.txt

# 7. Sync with remote
bd sync
```

### Managing Dependencies

If tasks have dependencies:

```bash
# Create tasks
bd create --title="Implement feature X" --type=feature
bd create --title="Write tests for X" --type=task

# Add dependency (tests depend on feature)
bd dep add beads-yyy beads-xxx

# View blocked tasks
bd blocked
```

### Project Health

```bash
# View statistics
bd stats

# Check for issues
bd doctor

# List all open tasks
bd list --status=open

# List tasks in progress
bd list --status=in_progress
```

---

## Migration from prd.json

If you have an existing `prd.json` file, use the migration script:

```bash
./.idea/scripts/prd_migrate.sh prd/prd.json
```

This will:
1. Read all tasks from prd.json
2. Create corresponding bd issues
3. Preserve task completion status
4. Include implementation steps in issue descriptions

After migration:
```bash
# Verify migrated tasks
bd list --status=open

# Sync to remote
bd sync

# Archive old prd.json
mv prd/prd.json prd/prd.json.backup
```

---

## Common Mistakes to Avoid

❌ **Not using bd CLI for task management**
✅ Always use `bd create`, `bd ready`, `bd close`, etc.

❌ **Not updating progress.txt**
✅ Always append to prd/progress.txt after completing work

❌ **Not syncing changes**
✅ Run `bd sync` at the end of each session

❌ **Creating subtasks in bd for large tasks**
✅ Log blockers in progress.txt instead, keep bd for main tasks

---

## Script Usage

### prd.sh - Automated Loop Mode

Run the PRD workflow multiple times:

```bash
# Run 5 iterations
./.idea/scripts/prd.sh --loop 5

# Run once with live output
./.idea/scripts/prd.sh
```

### prd_once.sh - Single Task

Run the PRD workflow once:

```bash
./.idea/scripts/prd_once.sh
```

Both scripts will:
1. Find next available task via `bd ready`
2. Claim and implement the task
3. Close the task with `bd close`
4. Update prd/progress.txt
5. Sync with `bd sync`
6. Run type checks and linting
