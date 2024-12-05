/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import type { NextRequest } from 'next/server'
import { createRemoteJWKSet, decodeJwt, JWTPayload, jwtVerify } from 'jose'
import { JWT } from 'next-auth/jwt'

import { logger } from '@/lib/logger'

/**
 * @param req The Next.js request object
 * @returns An encrypted string (JWT) or an empty string if not present
 */
export function tokenFromRequest(req: NextRequest): string {
  let token = ''
  const authorizationHeader = req.headers.get('authorization')
  if (authorizationHeader?.split(' ')[0] === 'Bearer') {
    const urlEncodedToken = authorizationHeader.split(' ')[1]
    token = decodeURIComponent(urlEncodedToken)
  }
  return token
}

/**
 * @param req The Next.js request object
 * @returns A decoded JWT object or null if token is invalid/not present
 */
export async function verifyToken(req: NextRequest): Promise<JWT | null> {
  const jwks = createRemoteJWKSet(new URL(`${process.env.AUTH_AUTH0_ISSUER}/.well-known/jwks.json`))
  let token = tokenFromRequest(req)

  try {
    const { payload } = await jwtVerify(token, jwks, {
      audience: process.env.AUTH_AUTH0_AUDIENCE,
      issuer: `${process.env.AUTH_AUTH0_ISSUER}/`, // jose expects a trailing slash
      requiredClaims: ['permissions'],
    })
    return payload
  } catch (err) {
    logger.error(`Error verifying token: ${err}`)
    return null
  }
}

/**
 * @param token An encrypted JWT token
 * @returns A decoded JWT payload (unverified) or null if token is invalid
 */
export function decodeToken(token: string): JWTPayload | null {
  try {
    const payload = decodeJwt(token)
    return payload
  } catch (err) {
    logger.error(`Error decoding token: ${err}`)
    return null
  }
}

/**
 * @param token An encrypted JWT token
 * @param permission A string representing a particular permission
 * @returns True if the token has the specified permission, false otherwise
 */
export function hasPermission(token: string, permission: string): boolean {
  const payload = decodeToken(token)
  if (!payload) return false

  if (payload?.permissions) {
    return payload.permissions.includes(permission)
  }
  return false
}

declare module 'jose' {
  interface JWTPayload {
    permissions?: string[]
  }
}
