# Copilot Instructions

## Commands

```bash
pnpm run dev        # Start dev server (Turbopack)
pnpm run build      # Production build
pnpm run lint       # ESLint
pnpm run test       # Vitest (run once)
pnpm run test:watch # Vitest (watch mode)
```

Run a single test file: `pnpm run test __tests__/lib/jwt.test.ts`

## Architecture

Docuhost is a Next.js 16 App Router application that provides:
- A **REST API** (`/api/[version]/documents`) for uploading documents (base64-encoded) and retrieving them by ID
- A **browser UI** for browsing/viewing documents (requires Auth0 login with `list:documents` permission)

Documents are stored in **MongoDB GridFS** (binary content) via a singleton `clientPromise` exported from `lib/mongodb.ts`. After upload, a short link is generated via a pluggable shortlink service (TinyURL is active; Bitly is implemented but commented out in `lib/shortlink/index.ts`).

### API versioning

All routes are under `/api/[version]/` and accept `v1` or `v2`. The only behavioral difference between versions is the POST response shape:
- `v1`: returns `HostedDocument` directly
- `v2`: wraps it as `{ document: HostedDocument }`

### Auth flow

Two parallel authorization paths exist for `POST /api/[version]/documents`:
1. **Browser session** — Auth.js v5 (next-auth beta) with Auth0 provider. The session `accessToken` is the raw Auth0 access token stored in the JWT callback.
2. **Bearer token** — Direct API access. Token is verified against Auth0 JWKS (`lib/jwt.ts`), then decoded to check the `permissions` array.

Permission strings used: `create:documents`, `list:documents`.

Auth.js proxy (`proxy.ts`) protects all routes **except** `/api/*`, `/_next/*`, and `favicon.ico`.

### Environment variables

Copy `example.env` to `.env.local`. Key variables:
- `APP_URL` — base URL for document links; should be the full path to `/api/v2/documents`
- `AUTH_AUTH0_ISSUER`, `AUTH_AUTH0_AUDIENCE` (`urn:docuhost`), `AUTH_AUTH0_ID`, `AUTH_AUTH0_SECRET`
- `AUTH_SECRET` — generate with `pnpm exec auth secret`
- `MONGODB_*` — connection config; `MONGODB_PROTOCOL` supports `mongodb+srv` for Atlas
- `SHORTLINK_SERVICE_URL`, `SHORTLINK_API_KEY`, `SHORTLINK_DOMAIN`

## Conventions

- **Tailwind config**: No `tailwind.config.ts` — Tailwind v4 theme (shadcn/ui colors, border radius) lives in the `@theme` block in `app/globals.css`
- **Tests**: Vitest test files live in `__tests__/` mirroring the source tree. Mock modules with `vi.mock()`; use regular `function` (not arrow) when mocking constructors (e.g. `GridFSBucket`)
- **Copyright header**: Every source file begins with a JSDoc comment block: `Copyright (c) [year] Quadient Group AG` + `SPDX-License-Identifier: MIT`
- **Path aliases**: Use `@/` for all internal imports (`@/lib/...`, `@/components/...`)
- **Import order**: Enforced by `@ianvs/prettier-plugin-sort-imports` — react → next → third-party → `@/lib` → `@/components/ui` → `@/components` → relative
- **Prettier**: `semi: false`, `singleQuote: true`, `printWidth: 100`, trailing commas (`es5`)
- **Next.js params**: Route segment params are `Promise`s in App Router — always `await context.params`
- **Logging**: Use `logger` from `@/lib/logger` (Winston, JSON format). Log all API responses with `logger.info(requestInfo, { status: N })`
- **UI components**: shadcn/ui primitives live in `components/ui/`; `cn()` from `@/lib/utils` merges Tailwind classes
- **Docker**: `make docker` builds multi-arch (arm64 + amd64) and pushes to `registry.sptcloud.com/spt/docuhost`
