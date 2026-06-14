/**
 * Copyright (c) 2026 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import type { NextRequest } from 'next/server'
import { createRemoteJWKSet, jwtVerify, type JWTPayload } from 'jose'

import { auth } from '@/auth'
import { logger } from '@/lib/logger'

declare module 'jose' {
  interface JWTPayload {
    permissions?: string[]
  }
}

function tokenFromRequest(req: NextRequest): string {
  const header = req.headers.get('authorization')
  if (header?.split(' ')[0] !== 'Bearer') return ''
  return decodeURIComponent(header.split(' ')[1])
}

let jwks: ReturnType<typeof createRemoteJWKSet> | undefined

function getJwks() {
  if (!jwks) {
    jwks = createRemoteJWKSet(
      new URL(`${process.env.AUTH_AUTH0_ISSUER}/.well-known/jwks.json`)
    )
  }
  return jwks
}

export async function verifyAccessToken(token: string): Promise<JWTPayload | null> {
  if (!token) return null
  try {
    const { payload } = await jwtVerify(token, getJwks(), {
      audience: process.env.AUTH_AUTH0_AUDIENCE,
      issuer: `${process.env.AUTH_AUTH0_ISSUER}/`,
      requiredClaims: ['permissions'],
    })
    return payload
  } catch (err) {
    logger.error(`Error verifying token: ${err}`)
    return null
  }
}

function hasPermission(payload: JWTPayload, permission: string): boolean {
  return Array.isArray(payload.permissions) && payload.permissions.includes(permission)
}

/**
 * Resolves an authenticated principal for the request and checks a permission.
 * Tries a Bearer token first (M2M / API clients), then falls back to a
 * NextAuth session (interactive users). Both paths are verified against the
 * Auth0 JWKS and required to carry the requested permission.
 */
export async function authorize(
  request: NextRequest,
  permission: string
): Promise<JWTPayload | null> {
  const bearer = tokenFromRequest(request)
  if (bearer) {
    const payload = await verifyAccessToken(bearer)
    return payload && hasPermission(payload, permission) ? payload : null
  }

  const session = await auth()
  if (session?.accessToken) {
    const payload = await verifyAccessToken(session.accessToken)
    return payload && hasPermission(payload, permission) ? payload : null
  }

  return null
}

/**
 * Server-component helper: returns true when the current session carries a
 * verified Auth0 access token with the requested permission. Use for UI
 * gating (e.g., showing/hiding controls). Route handlers should use
 * `authorize()` instead — it also accepts Bearer tokens.
 */
export async function sessionHasPermission(permission: string): Promise<boolean> {
  const session = await auth()
  if (!session?.accessToken) return false
  const payload = await verifyAccessToken(session.accessToken)
  return payload ? hasPermission(payload, permission) : false
}
