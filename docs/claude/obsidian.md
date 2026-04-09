# Obsidian Notes

Use the Obsidian CLI (`obsidian`) to log task notes to the `tirith` vault.
Notes live in the `work/project-name/` folder and are resolved by name with `file=`.
Keep it succinct and neutral.

## Creating a note

```bash
obsidian create vault=tirith path="work/project-name/XWWP-1234.md" content="## XWWP-1234 - Title\n\nShort description."
```

## Setting frontmatter properties

Always set properties immediately after creating a note, including the Jira ID when possible:

```bash
obsidian property:set vault=tirith file="XWWP-1234" name="task"  value="XWWP-1234"  type=list
```

## Appending content

```bash
obsidian append vault=tirith file="XWWP-1234" content="\nExtra content here."
```

## Reading a note

```bash
obsidian read vault=tirith file="XWWP-1234"
```

## Note content guidelines

- What was implemented and why
- Key decisions or trade-offs
- Files changed
- Anything non-obvious for future reference
- Date of the inclusition
- On updates add the date of the update

## Tips

- Always use `file=<name>` (not `path=`) on existing notes - resolves correctly after the note is moved
- Use `\n` for newlines in `content=` values
