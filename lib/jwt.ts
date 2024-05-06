/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { createRemoteJWKSet, jwtVerify } from 'jose'
import type { NextRequest } from 'next/server'
import { logger } from '@/lib/logger'

/**
 * Default contents of the Auth0 access token.
 */
export interface DefaultJWT extends Record<string, unknown> {
  iss?: string | null
  sub?: string | null
  aud?: string | string[] | null
  scope?: string | string[] | null
}

/**
 * Returned by the `verifyToken` function as the decoded payload.
 */
export interface JWT extends Record<string, unknown>, DefaultJWT {}

/**
 *
 * @param req The Next.js request object
 * @returns A decoded JWT object or null if token is invalid/not present
 */
export async function verifyToken(req: NextRequest): Promise<JWT | null> {
  const jwks = createRemoteJWKSet(new URL(`${process.env.AUTH_AUTH0_ISSUER}/.well-known/jwks.json`))

  let token = ''
  const authorizationHeader = req.headers.get('authorization')
  if (authorizationHeader?.split(' ')[0] === 'Bearer') {
    const urlEncodedToken = authorizationHeader.split(' ')[1]
    token = decodeURIComponent(urlEncodedToken)
  }

  try {
    const { payload } = await jwtVerify(token, jwks, {
      audience: process.env.AUTH_AUTH0_AUDIENCE,
      issuer: `${process.env.AUTH_AUTH0_ISSUER}/`, // jose expects a trailing slash
      requiredClaims: ['scope'],
    })
    return payload
  } catch (err) {
    logger.error(`Error verifying token: ${err}`)
    return null
  }
}
