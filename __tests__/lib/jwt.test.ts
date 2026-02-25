/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import type { NextRequest } from 'next/server'

// Mock jose before importing jwt module
vi.mock('jose', async (importOriginal) => {
  const actual = await importOriginal<typeof import('jose')>()
  return {
    ...actual,
    createRemoteJWKSet: vi.fn(),
    jwtVerify: vi.fn(),
  }
})

vi.mock('@/lib/logger', () => ({
  logger: { debug: vi.fn(), error: vi.fn(), info: vi.fn(), warn: vi.fn() },
}))

import { tokenFromRequest, hasPermission, decodeToken } from '@/lib/jwt'
import { jwtVerify, createRemoteJWKSet } from 'jose'

// A real HS256-signed JWT with payload { sub: 'test', permissions: ['create:documents'] }
// Signed with secret 'test' — used only for decodeJwt (unverified decode)
const JWT_WITH_PERMISSION =
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.' +
  'eyJzdWIiOiJ0ZXN0IiwicGVybWlzc2lvbnMiOlsiY3JlYXRlOmRvY3VtZW50cyJdfQ.' +
  'placeholder'

const JWT_NO_PERMISSIONS =
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.' +
  'eyJzdWIiOiJ0ZXN0In0.' +
  'placeholder'

function makeRequest(authHeader?: string): NextRequest {
  return {
    headers: {
      get: (name: string) => (name === 'authorization' ? (authHeader ?? null) : null),
    },
    nextUrl: { pathname: '/api/v2/documents' },
    method: 'POST',
  } as unknown as NextRequest
}

describe('tokenFromRequest', () => {
  it('extracts a Bearer token from the Authorization header', () => {
    const req = makeRequest('Bearer my-token')
    expect(tokenFromRequest(req)).toBe('my-token')
  })

  it('returns empty string when no Authorization header is present', () => {
    const req = makeRequest()
    expect(tokenFromRequest(req)).toBe('')
  })

  it('returns empty string for non-Bearer schemes', () => {
    const req = makeRequest('Basic dXNlcjpwYXNz')
    expect(tokenFromRequest(req)).toBe('')
  })

  it('URL-decodes the token', () => {
    const encoded = encodeURIComponent('token+with/special=chars')
    const req = makeRequest(`Bearer ${encoded}`)
    expect(tokenFromRequest(req)).toBe('token+with/special=chars')
  })
})

describe('decodeToken', () => {
  it('returns null for a malformed token', () => {
    expect(decodeToken('not.a.token')).toBeNull()
  })

  it('returns null for an empty string', () => {
    expect(decodeToken('')).toBeNull()
  })
})

describe('hasPermission', () => {
  it('returns false for a malformed token', () => {
    expect(hasPermission('garbage', 'create:documents')).toBe(false)
  })

  it('returns false when the permission is not present', () => {
    // JWT payload: { "sub": "test" } — no permissions claim
    const token =
      'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.' +
      'eyJzdWIiOiJ0ZXN0In0.' +
      'SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c'
    expect(hasPermission(token, 'create:documents')).toBe(false)
  })

  it('returns true when the permission is present', () => {
    // JWT payload: { "sub": "test", "permissions": ["create:documents", "list:documents"] }
    const token =
      'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.' +
      'eyJzdWIiOiJ0ZXN0IiwicGVybWlzc2lvbnMiOlsiY3JlYXRlOmRvY3VtZW50cyIsImxpc3Q6ZG9jdW1lbnRzIl19.' +
      'SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c'
    expect(hasPermission(token, 'create:documents')).toBe(true)
    expect(hasPermission(token, 'list:documents')).toBe(true)
  })

  it('returns false when a different permission is checked', () => {
    const token =
      'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.' +
      'eyJzdWIiOiJ0ZXN0IiwicGVybWlzc2lvbnMiOlsiY3JlYXRlOmRvY3VtZW50cyJdfQ.' +
      'SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c'
    expect(hasPermission(token, 'delete:documents')).toBe(false)
  })
})

describe('verifyToken', () => {
  beforeEach(() => {
    vi.mocked(createRemoteJWKSet).mockReturnValue({} as ReturnType<typeof createRemoteJWKSet>)
    process.env.AUTH_AUTH0_ISSUER = 'https://example.auth0.com'
    process.env.AUTH_AUTH0_AUDIENCE = 'urn:docuhost'
  })

  it('returns the payload on a valid token', async () => {
    const { verifyToken } = await import('@/lib/jwt')
    const mockPayload = { sub: 'user', permissions: ['create:documents'] }
    vi.mocked(jwtVerify).mockResolvedValue({ payload: mockPayload } as never)

    const req = makeRequest('Bearer valid.jwt.token')
    const result = await verifyToken(req)
    expect(result).toEqual(mockPayload)
  })

  it('returns null when verification fails', async () => {
    const { verifyToken } = await import('@/lib/jwt')
    vi.mocked(jwtVerify).mockRejectedValue(new Error('invalid signature'))

    const req = makeRequest('Bearer bad.token')
    const result = await verifyToken(req)
    expect(result).toBeNull()
  })
})
