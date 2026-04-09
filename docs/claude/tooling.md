# Tooling Notes

## Linting and formatting

- Use Biome for linting and formatting, not ESLint.
- Run `pnpm lint` for workspace-wide linting.
- For targeted changes, run `biome lint` or `biome format` on specific files.

## Development notes

- Web app runs with HTTPS in dev mode (`--experimental-https` flag).
- PocketBase types are generated via `npm run typegen` in `/apps/base`.
- Next.js transpiles `@thom/ui` for hot reload.
- 3D pet animations are stored in `/public/static/` as GIF files.
- Route definitions are centralized in `/src/constants/routes.ts`.
- Build dependencies are managed in `turbo.json` (for example, dev depends on `^ui`).
