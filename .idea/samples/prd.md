# prd.json

**Location:** `.idea/to-work-on`

**Name:** name the file with the name of the epic and append -prd.json example: tests-cleanup-prd.json

This file contains an array of tasks. Each task object MUST have these fields:

### Required Fields

- **category** (string): Either `"functional"` or `"ui"`
- **description** (string): Simple description, max 240 characters
- **steps** (array): Array of step strings, max 10 items
- **passes** (boolean): Starts at `false`, set to `true` when task is complete

### Example Structure

```json
[
  {
    "category": "functional",
    "description": "Update hooks/index.ts to export keyboard shortcuts hook",
    "steps": [
      "Open app/javascript/admins/SmartReports/components/EditSmartReport/hooks/index.ts",
      "Add export { useSmartReportKeyboardShortcuts } from './useSmartReportKeyboardShortcuts'",
      "Import and call useSmartReportKeyboardShortcuts() in EditSmartReportComponent"
    ],
    "passes": false
  },
  {
    "category": "ui",
    "description": "Add zen mode toggle button to SmartReports header",
    "steps": [
      "Create ZenModeToggle component in components/",
      "Add toggle state to useSmartReportStore",
      "Import and render toggle in ReportHeader component",
      "Add onClick handler to toggle zen mode state"
    ],
    "passes": false
  }
]
```
