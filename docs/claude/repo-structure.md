# Repo Structure

## Monorepo structure

- `/apps/web-app` - Next.js 16 application with React 19
- `/apps/base` - Go-based PocketBase backend
- `/packages/ui` - Shared UI component library (Tailwind + Radix UI)
- `/packages/typescript-config` - Shared TypeScript configuration

## Key directories (web app)

```
/apps/web-app/src/
├── app/[locale]/        # Next.js App Router with i18n
├── features/            # Feature modules (hexagonal architecture)
├── components/          # Global reusable components
├── hooks/               # Global React hooks
├── lib/                 # Utilities (result.ts, db utilities)
├── constants/           # App constants (routes.ts for centralized routes)
├── types/               # Global TypeScript types
│   └── pocketbase-types.ts  # Auto-generated from PocketBase
├── data/mocks/          # Mock data for development/testing
├── i18n/                # Internationalization
└── scripts/             # Build-time scripts
```

Path alias: `_/*` maps to `/src/*` (configured in `tsconfig.json`).
