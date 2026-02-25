/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { NextRequest } from 'next/server'

vi.mock('@/lib/logger', () => ({
  logger: { debug: vi.fn(), error: vi.fn(), info: vi.fn(), warn: vi.fn() },
}))

const mockInsert = vi.fn()
vi.mock('@/lib/api/documents', () => ({ insert: mockInsert }))

const mockAuth = vi.fn()
vi.mock('@/auth', () => ({ auth: mockAuth }))

const mockVerifyToken = vi.fn()
const mockHasPermission = vi.fn()
const mockTokenFromRequest = vi.fn()
vi.mock('@/lib/jwt', () => ({
  verifyToken: mockVerifyToken,
  hasPermission: mockHasPermission,
  tokenFromRequest: mockTokenFromRequest,
}))

function makePostRequest(body: object, authHeader?: string): NextRequest {
  const headers = new Headers({ 'content-type': 'application/json' })
  if (authHeader) headers.set('authorization', authHeader)
  return new NextRequest('http://localhost/api/v2/documents', {
    method: 'POST',
    headers,
    body: JSON.stringify(body),
  })
}

const sampleDoc = {
  id: '64f1a2b3c4d5e6f7a8b9c0d1',
  filename: 'test.pdf',
  contentType: 'application/pdf',
  fileBase64: '[stored]',
  url: 'http://localhost/api/v2/documents/64f1a2b3c4d5e6f7a8b9c0d1',
  shortLink: 'https://tinyurl.com/test',
}

beforeEach(() => {
  vi.clearAllMocks()
  mockAuth.mockResolvedValue(null)
  mockVerifyToken.mockResolvedValue(null)
  mockHasPermission.mockReturnValue(false)
  mockTokenFromRequest.mockReturnValue('')
  mockInsert.mockResolvedValue(sampleDoc)
})

describe('POST /api/[version]/documents', () => {
  describe('authorization', () => {
    it('returns 401 when no session and no valid Bearer token', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      const req = makePostRequest({ filename: 'test.pdf' })
      const res = await POST(req, { params: Promise.resolve({ version: 'v2' }) })
      expect(res.status).toBe(401)
    })

    it('returns 401 when session exists but lacks create:documents permission', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      mockAuth.mockResolvedValue({ accessToken: 'some.token' })
      mockHasPermission.mockReturnValue(false)

      const req = makePostRequest({ filename: 'test.pdf' })
      const res = await POST(req, { params: Promise.resolve({ version: 'v2' }) })
      expect(res.status).toBe(401)
    })

    it('authorizes via session token with create:documents permission', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      mockAuth.mockResolvedValue({ accessToken: 'valid.session.token' })
      mockHasPermission.mockReturnValue(true)

      const req = makePostRequest({ filename: 'test.pdf', contentType: 'application/pdf', fileBase64: 'abc' })
      const res = await POST(req, { params: Promise.resolve({ version: 'v2' }) })
      expect(res.status).toBe(201)
    })

    it('authorizes via Bearer token when no session', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      mockAuth.mockResolvedValue(null)
      mockVerifyToken.mockResolvedValue({ sub: 'machine', permissions: ['create:documents'] })
      mockTokenFromRequest.mockReturnValue('valid.bearer.token')
      mockHasPermission.mockReturnValue(true)

      const req = makePostRequest(
        { filename: 'test.pdf', contentType: 'application/pdf', fileBase64: 'abc' },
        'Bearer valid.bearer.token'
      )
      const res = await POST(req, { params: Promise.resolve({ version: 'v2' }) })
      expect(res.status).toBe(201)
    })
  })

  describe('version handling', () => {
    beforeEach(() => {
      mockAuth.mockResolvedValue({ accessToken: 'valid.session.token' })
      mockHasPermission.mockReturnValue(true)
    })

    it('v2 wraps the document in a { document } envelope', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      const req = makePostRequest({ filename: 'test.pdf', contentType: 'application/pdf', fileBase64: 'abc' })
      const res = await POST(req, { params: Promise.resolve({ version: 'v2' }) })
      const body = await res.json()
      expect(body).toHaveProperty('document')
      expect(body.document.filename).toBe('test.pdf')
    })

    it('v1 returns the document directly (no envelope)', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      const req = makePostRequest({ filename: 'test.pdf', contentType: 'application/pdf', fileBase64: 'abc' })
      const res = await POST(req, { params: Promise.resolve({ version: 'v1' }) })
      const body = await res.json()
      expect(body).not.toHaveProperty('document')
      expect(body.filename).toBe('test.pdf')
    })

    it('returns 400 for an unsupported version', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      const req = makePostRequest({ filename: 'test.pdf' })
      const res = await POST(req, { params: Promise.resolve({ version: 'v3' }) })
      expect(res.status).toBe(400)
    })

    it('sets the Location header on 201', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      const req = makePostRequest({ filename: 'test.pdf', contentType: 'application/pdf', fileBase64: 'abc' })
      const res = await POST(req, { params: Promise.resolve({ version: 'v2' }) })
      expect(res.headers.get('location')).toBe(sampleDoc.url)
    })
  })

  describe('insert failure', () => {
    it('returns 500 when insert returns null', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      mockAuth.mockResolvedValue({ accessToken: 'valid.session.token' })
      mockHasPermission.mockReturnValue(true)
      mockInsert.mockResolvedValue(null)

      const req = makePostRequest({ filename: 'test.pdf', contentType: 'application/pdf', fileBase64: 'abc' })
      const res = await POST(req, { params: Promise.resolve({ version: 'v2' }) })
      expect(res.status).toBe(500)
    })
  })
})
