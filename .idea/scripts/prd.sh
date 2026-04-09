#!/bin/bash
set -e

# Define the PRD prompt using bd CLI
PROMPT="MANDATORY FILES TO READ AND UPDATE:
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

CRITICAL REQUIREMENTS:
- You MUST close the task with 'bd close <id>' when complete
- You MUST append to prd/progress.txt with your progress notes be succinct, describe any error that occurred while implementing the task
- If the task is too big to complete, append it in prd/progress.txt as a blocker and output: <promise>ERROR</promise>
- ONLY WORK ON A SINGLE TASK
- When 'bd ready' shows no available tasks, output: <promise>COMPLETE</promise>

WORKFLOW EXAMPLE:
1. bd ready                              # Find next task
2. bd update <id> --status=in_progress   # Claim the task
3. [implement the task]
4. bd close <id>                         # Mark complete
5. Append summary to prd/progress.txt
6. bd sync                               # Sync to remote"

# Parse arguments
LOOP_COUNT=""
while [[ $# -gt 0 ]]; do
  case $1 in
    --loop)
      LOOP_COUNT="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--loop <iterations>]"
      echo "  --loop <iterations>  Run the PRD workflow N times"
      echo "  (no flag)            Run once with live output"
      exit 1
      ;;
  esac
done

# If no loop flag, execute once with live output (interactive mode)
if [ -z "$LOOP_COUNT" ]; then
  claude --permission-mode acceptEdits "$PROMPT"
  exit 0
fi

# Validate loop count
if ! [[ "$LOOP_COUNT" =~ ^[0-9]+$ ]] || [ "$LOOP_COUNT" -lt 1 ]; then
  echo "Error: --loop requires a positive integer"
  exit 1
fi

# Run the PRD workflow in a loop
for ((i=1; i<=LOOP_COUNT; i++)); do
  echo "Iteration $i"
  echo "----------------"

  result=$(claude --permission-mode acceptEdits -p "$PROMPT")
  echo "$result"

  if [[ "$result" == *"<promise>ERROR</promise>"* ]]; then
    echo "Error updating files, stopping."
    exit 1
  fi

  if [[ "$result" == *"<promise>COMPLETE</promise>"* ]]; then
    echo "PRD complete, exiting."
    exit 0
  fi
done
