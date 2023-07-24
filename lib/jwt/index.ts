/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { createRemoteJWKSet, jwtVerify } from 'jose'
import type { NextRequest } from 'next/server'
import type { JWT } from './types'

export * from './types'

/**
 *
 * @param req The Next.js request object
 * @returns A decoded JWT object or null if token is invalid/not present
 */
export async function verifyToken(req: NextRequest): Promise<JWT | null> {
  const jwks = createRemoteJWKSet(
    new URL(`${process.env.AUTH0_DOMAIN}/.well-known/jwks.json`)
  )
  const requiredClaims = ['scope']

  let token = ''
  const authorizationHeader = req.headers.get('authorization')
  if (authorizationHeader?.split(' ')[0] === 'Bearer') {
    const urlEncodedToken = authorizationHeader.split(' ')[1]
    token = decodeURIComponent(urlEncodedToken)
  }

  if (!token) {
    return null
  }

  try {
    const { payload } = await jwtVerify(token, jwks, {
      audience: process.env.AUTH0_AUDIENCE,
      issuer: `${process.env.AUTH0_DOMAIN}/`,
      requiredClaims: requiredClaims,
    })
    return payload
  } catch (err) {
    console.error(err)
    return null
  }
}
