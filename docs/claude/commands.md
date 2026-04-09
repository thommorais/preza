# Commands

## Development

- `pnpm dev` - Run all apps in dev mode
- `pnpm web:dev` - Run only web-app in dev mode

## Build

- `pnpm build` - Build all apps
- `pnpm web:build` - Build only web-app

## Lint/Typecheck

- `pnpm lint` - Lint all workspaces with Biome
- `pnpm typecheck` - Type check all workspaces

## Icon generation

- `pnpm build:icons` - Build icons for UI package
- `pnpm web:icons` - Build icons for web-app

## Workspace-specific

### /apps/web-app

- `npm run dev` - Next.js dev server with HTTPS and turbopack
- `npm run lint` - Biome lint ./src
- `npm run format` - Biome format ./src

### /apps/base

- `npm run dev` - Start PocketBase server
- `npm run typegen` - Generate TypeScript types from PocketBase schema
