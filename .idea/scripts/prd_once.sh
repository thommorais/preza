#!/bin/bash
set -e

claude --permission-mode acceptEdits "MANDATORY FILES TO READ AND UPDATE:
- prd/progress.txt

MANDATORY WORKFLOW:
Use 'bd' CLI for task management.

TASKS (ALL ARE MANDATORY):
1. Run 'bd ready' to find the next available task (no blockers)
2. Run 'bd show <id>' to get task details and implementation steps
3. Work ONLY on that single task
4. After implementing, run 'bd close <id>' to mark the task as complete
5. You MUST append a brief summary to prd/progress.txt describing what you did
6. Run 'bd sync' to sync changes with git remote
7. Run type checks and linting only for the changes made
8. Make a git commit of that feature

CRITICAL REQUIREMENTS:
- You MUST close the task with 'bd close <id>' when complete
- You MUST append to prd/progress.txt with your progress notes, describe any error that occurred while implementing the task
- ONLY WORK ON A SINGLE TASK
- If the task is too big to complete, log it in prd/progress.txt as a blocker and output: <promise>ERROR</promise>
- If you encounter ANY issues append to progress.txt, output: <promise>ERROR</promise> and explain the issue
- When 'bd ready' shows no available tasks, output: <promise>COMPLETE</promise>

WORKFLOW EXAMPLE:
1. bd ready                              # Find next task
2. bd update <id> --status=in_progress   # Claim the task
3. [implement the task]
4. bd close <id>                         # Mark complete
5. Append summary to prd/progress.txt
6. bd sync                               # Sync to remote"
