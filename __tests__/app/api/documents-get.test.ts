/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { NextRequest } from 'next/server'

vi.mock('@/lib/logger', () => ({
  logger: { debug: vi.fn(), error: vi.fn(), info: vi.fn(), warn: vi.fn() },
}))

const mockFetch = vi.fn()
const mockIsValidObjectId = vi.fn()
vi.mock('@/lib/api/documents', () => ({
  fetch: mockFetch,
  isValidObjectId: mockIsValidObjectId,
}))

function makeGetRequest(id: string): NextRequest {
  return new NextRequest(`http://localhost/api/v2/documents/${id}`)
}

const mockReadStream = {
  on: vi.fn((event: string, cb: (...args: unknown[]) => void) => {
    if (event === 'end') cb()
    return mockReadStream
  }),
  destroy: vi.fn(),
}

const sampleFile = {
  filename: 'report.pdf',
  content: mockReadStream,
  contentType: 'application/pdf',
  size: 12345,
}

beforeEach(() => {
  vi.clearAllMocks()
  mockIsValidObjectId.mockImplementation((id: string) => {
    if (id === 'invalid-id') throw new Error('Invalid ObjectId')
    return id
  })
  mockFetch.mockResolvedValue(sampleFile)
})

describe('GET /api/[version]/documents/[id]', () => {
  describe('version validation', () => {
    it('returns 400 for an unsupported version', async () => {
      const { GET } = await import('@/app/api/[version]/documents/[id]/route')
      const id = '64f1a2b3c4d5e6f7a8b9c0d1'
      const req = makeGetRequest(id)
      const res = await GET(req, { params: Promise.resolve({ version: 'v3', id }) })
      expect(res.status).toBe(400)
    })

    it('accepts v1', async () => {
      const { GET } = await import('@/app/api/[version]/documents/[id]/route')
      const id = '64f1a2b3c4d5e6f7a8b9c0d1'
      const req = makeGetRequest(id)
      const res = await GET(req, { params: Promise.resolve({ version: 'v1', id }) })
      expect(res.status).toBe(200)
    })

    it('accepts v2', async () => {
      const { GET } = await import('@/app/api/[version]/documents/[id]/route')
      const id = '64f1a2b3c4d5e6f7a8b9c0d1'
      const req = makeGetRequest(id)
      const res = await GET(req, { params: Promise.resolve({ version: 'v2', id }) })
      expect(res.status).toBe(200)
    })
  })

  describe('ObjectId validation', () => {
    it('returns 400 for an invalid ObjectId', async () => {
      const { GET } = await import('@/app/api/[version]/documents/[id]/route')
      const req = makeGetRequest('invalid-id')
      const res = await GET(req, { params: Promise.resolve({ version: 'v2', id: 'invalid-id' }) })
      expect(res.status).toBe(400)
    })
  })

  describe('document retrieval', () => {
    it('returns 404 when the document does not exist', async () => {
      const { GET } = await import('@/app/api/[version]/documents/[id]/route')
      mockFetch.mockResolvedValue(null)
      const id = '64f1a2b3c4d5e6f7a8b9c0d1'
      const req = makeGetRequest(id)
      const res = await GET(req, { params: Promise.resolve({ version: 'v2', id }) })
      expect(res.status).toBe(404)
    })

    it('returns 200 with correct content-type and content-length headers', async () => {
      const { GET } = await import('@/app/api/[version]/documents/[id]/route')
      const id = '64f1a2b3c4d5e6f7a8b9c0d1'
      const req = makeGetRequest(id)
      const res = await GET(req, { params: Promise.resolve({ version: 'v2', id }) })
      expect(res.status).toBe(200)
      expect(res.headers.get('content-type')).toBe('application/pdf')
      expect(res.headers.get('content-length')).toBe('12345')
      expect(res.headers.get('content-disposition')).toBe('inline')
    })
  })
})
