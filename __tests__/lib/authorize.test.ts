/**
 * Copyright (c) 2026 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import type { NextRequest } from 'next/server'

vi.mock('jose', async (importOriginal) => {
  const actual = await importOriginal<typeof import('jose')>()
  return {
    ...actual,
    createRemoteJWKSet: vi.fn(() => ({}) as never),
    jwtVerify: vi.fn(),
  }
})

vi.mock('@/lib/logger', () => ({
  logger: { debug: vi.fn(), error: vi.fn(), info: vi.fn(), warn: vi.fn() },
}))

const { mockAuth } = vi.hoisted(() => ({ mockAuth: vi.fn() }))
vi.mock('@/auth', () => ({ auth: mockAuth }))

import { jwtVerify } from 'jose'
import { authorize, sessionHasPermission, verifyAccessToken } from '@/lib/authorize'

function makeRequest(authHeader?: string): NextRequest {
  return {
    headers: {
      get: (name: string) => (name === 'authorization' ? (authHeader ?? null) : null),
    },
  } as unknown as NextRequest
}

beforeEach(() => {
  vi.clearAllMocks()
  process.env.AUTH_AUTH0_ISSUER = 'https://example.auth0.com'
  process.env.AUTH_AUTH0_AUDIENCE = 'urn:docuhost'
  mockAuth.mockResolvedValue(null)
})

describe('Bearer header parsing (via authorize)', () => {
  it('rejects a non-Bearer scheme without calling jose', async () => {
    expect(await authorize(makeRequest('Basic dXNlcjpwYXNz'), 'create:documents')).toBeNull()
    expect(jwtVerify).not.toHaveBeenCalled()
    expect(mockAuth).toHaveBeenCalled()
  })

  it('URL-decodes the Bearer token before verifying', async () => {
    vi.mocked(jwtVerify).mockResolvedValue({
      payload: { permissions: ['create:documents'] },
    } as never)
    const encoded = encodeURIComponent('token+with/special=chars')
    await authorize(makeRequest(`Bearer ${encoded}`), 'create:documents')
    expect(jwtVerify).toHaveBeenCalledWith('token+with/special=chars', expect.anything(), expect.anything())
  })
})

describe('verifyAccessToken', () => {
  it('returns the payload when jose verifies the token', async () => {
    const payload = { sub: 'user', permissions: ['create:documents'] }
    vi.mocked(jwtVerify).mockResolvedValue({ payload } as never)
    expect(await verifyAccessToken('good.token')).toEqual(payload)
  })

  it('returns null when jose throws', async () => {
    vi.mocked(jwtVerify).mockRejectedValue(new Error('bad sig'))
    expect(await verifyAccessToken('bad.token')).toBeNull()
  })

  it('returns null for an empty token without calling jose', async () => {
    expect(await verifyAccessToken('')).toBeNull()
    expect(jwtVerify).not.toHaveBeenCalled()
  })

})

describe('authorize', () => {
  it('authorizes a bearer token carrying the required permission', async () => {
    vi.mocked(jwtVerify).mockResolvedValue({
      payload: { sub: 'm2m', permissions: ['create:documents'] },
    } as never)
    const result = await authorize(makeRequest('Bearer good.token'), 'create:documents')
    expect(result).toMatchObject({ sub: 'm2m' })
    expect(mockAuth).not.toHaveBeenCalled()
  })

  it('rejects a bearer token missing the required permission', async () => {
    vi.mocked(jwtVerify).mockResolvedValue({
      payload: { sub: 'm2m', permissions: ['list:documents'] },
    } as never)
    expect(
      await authorize(makeRequest('Bearer good.token'), 'create:documents')
    ).toBeNull()
  })

  it('rejects a bearer token that fails verification', async () => {
    vi.mocked(jwtVerify).mockRejectedValue(new Error('invalid'))
    expect(
      await authorize(makeRequest('Bearer bad.token'), 'create:documents')
    ).toBeNull()
  })

  it('falls back to session when no bearer header is present', async () => {
    mockAuth.mockResolvedValue({ accessToken: 'session.token' })
    vi.mocked(jwtVerify).mockResolvedValue({
      payload: { sub: 'user', permissions: ['create:documents'] },
    } as never)
    const result = await authorize(makeRequest(), 'create:documents')
    expect(result).toMatchObject({ sub: 'user' })
  })

  it('returns null when neither bearer nor session is available', async () => {
    expect(await authorize(makeRequest(), 'create:documents')).toBeNull()
    expect(jwtVerify).not.toHaveBeenCalled()
  })

  it('verifies the session access token (not unverified decode)', async () => {
    mockAuth.mockResolvedValue({ accessToken: 'session.token' })
    vi.mocked(jwtVerify).mockRejectedValue(new Error('tampered'))
    expect(await authorize(makeRequest(), 'create:documents')).toBeNull()
  })
})

describe('sessionHasPermission', () => {
  it('returns false when there is no session', async () => {
    expect(await sessionHasPermission('list:documents')).toBe(false)
  })

  it('returns false when the session token fails verification', async () => {
    mockAuth.mockResolvedValue({ accessToken: 'session.token' })
    vi.mocked(jwtVerify).mockRejectedValue(new Error('bad sig'))
    expect(await sessionHasPermission('list:documents')).toBe(false)
  })

  it('returns false when the verified payload lacks the permission', async () => {
    mockAuth.mockResolvedValue({ accessToken: 'session.token' })
    vi.mocked(jwtVerify).mockResolvedValue({
      payload: { permissions: ['create:documents'] },
    } as never)
    expect(await sessionHasPermission('list:documents')).toBe(false)
  })

  it('returns true when the verified payload includes the permission', async () => {
    mockAuth.mockResolvedValue({ accessToken: 'session.token' })
    vi.mocked(jwtVerify).mockResolvedValue({
      payload: { permissions: ['list:documents'] },
    } as never)
    expect(await sessionHasPermission('list:documents')).toBe(true)
  })
})
