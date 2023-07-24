/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

/**
 * Default contents of the Auth0 access token.
 */
export interface DefaultJWT extends Record<string, unknown> {
  iss?: string | null
  sub?: string | null
  aud?: string | string[] | null
  iat?: number | null
  exp?: number | null
  azp?: string | null
  scope?: string | string[] | null
  gty?: string | null
}

/**
 * Returned by the `verifyToken` function as the decoded payload.
 */
export interface JWT extends Record<string, unknown>, DefaultJWT {}
