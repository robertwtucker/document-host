FROM node:22 AS base
ENV PNPM_HOME=/pnpm
ENV PATH=$PNPM_HOME:$PATH
RUN corepack enable
COPY . /app
WORKDIR /app

FROM base AS prod-deps
RUN --mount=type=cache,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM base AS builder
COPY .env-docker.local ./
RUN --mount=type=cache,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM gcr.io/distroless/nodejs22-debian12 AS runner
WORKDIR /app

ENV NODE_ENV production
ENV HOSTNAME 0.0.0.0
ENV PORT 3000
ENV NEXT_TELEMETRY_DISABLED 1

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=prod-deps /app/node_modules ./node_modules

EXPOSE 3000

# server.js is created by next build from the standalone output
# https://nextjs.org/docs/pages/api-reference/next-config-js/output
CMD ["./server.js"]
