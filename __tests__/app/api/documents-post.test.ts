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

const mockAuthorize = vi.fn()
vi.mock('@/lib/authorize', () => ({ authorize: mockAuthorize }))

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
  mockAuthorize.mockResolvedValue(null)
  mockInsert.mockResolvedValue(sampleDoc)
})

describe('POST /api/[version]/documents', () => {
  describe('authorization', () => {
    it('returns 401 when authorize returns null', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      const req = makePostRequest({ filename: 'test.pdf' })
      const res = await POST(req, { params: Promise.resolve({ version: 'v2' }) })
      expect(res.status).toBe(401)
    })

    it('returns 201 when authorize returns a principal', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      mockAuthorize.mockResolvedValue({ sub: 'user', permissions: ['create:documents'] })

      const req = makePostRequest({ filename: 'test.pdf', contentType: 'application/pdf', fileBase64: 'abc' })
      const res = await POST(req, { params: Promise.resolve({ version: 'v2' }) })
      expect(res.status).toBe(201)
    })

    it('passes the required permission to authorize', async () => {
      const { POST } = await import('@/app/api/[version]/documents/route')
      mockAuthorize.mockResolvedValue({ sub: 'user', permissions: ['create:documents'] })

      const req = makePostRequest({ filename: 'test.pdf', contentType: 'application/pdf', fileBase64: 'abc' })
      await POST(req, { params: Promise.resolve({ version: 'v2' }) })
      expect(mockAuthorize).toHaveBeenCalledWith(req, 'create:documents')
    })
  })

  describe('version handling', () => {
    beforeEach(() => {
      mockAuthorize.mockResolvedValue({ sub: 'user', permissions: ['create:documents'] })
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
      mockAuthorize.mockResolvedValue({ sub: 'user', permissions: ['create:documents'] })
      mockInsert.mockResolvedValue(null)

      const req = makePostRequest({ filename: 'test.pdf', contentType: 'application/pdf', fileBase64: 'abc' })
      const res = await POST(req, { params: Promise.resolve({ version: 'v2' }) })
      expect(res.status).toBe(500)
    })
  })
})
