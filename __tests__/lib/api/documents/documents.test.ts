/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ObjectId } from 'mongodb'

vi.mock('@/lib/logger', () => ({
  logger: { debug: vi.fn(), error: vi.fn(), info: vi.fn(), warn: vi.fn() },
}))

// Mock the shortlink service
vi.mock('@/lib/shortlink', () => ({
  shorten: vi.fn().mockResolvedValue({ url: 'http://example.com/doc/123', shortlink: 'https://tinyurl.com/abc' }),
}))

// Build a mock MongoDB client/bucket
const mockToArray = vi.fn()
const mockFind = vi.fn(() => ({ toArray: mockToArray }))
const mockOpenUploadStream = vi.fn()
const mockOpenDownloadStream = vi.fn()

const mockBucket = {
  find: mockFind,
  openUploadStream: mockOpenUploadStream,
  openDownloadStream: mockOpenDownloadStream,
}

vi.mock('mongodb', async (importOriginal) => {
  const actual = await importOriginal<typeof import('mongodb')>()
  return {
    ...actual,
    // Must use a regular function (not arrow) so it can be called with `new`
    GridFSBucket: vi.fn(function () { return mockBucket }),
  }
})

vi.mock('@/lib/mongodb', () => ({
  default: Promise.resolve({ db: () => ({}) }),
}))

beforeEach(() => {
  vi.clearAllMocks()
  process.env.APP_URL = 'http://localhost/api/v2/documents'
})

describe('isValidObjectId', () => {
  it('returns the hex string for a valid ObjectId', async () => {
    const { isValidObjectId } = await import('@/lib/api/documents')
    const id = new ObjectId().toHexString()
    expect(isValidObjectId(id)).toBe(id)
  })

  it('throws for an invalid ObjectId', async () => {
    const { isValidObjectId } = await import('@/lib/api/documents')
    expect(() => isValidObjectId('not-an-object-id')).toThrow()
  })

  it('throws for an empty string', async () => {
    const { isValidObjectId } = await import('@/lib/api/documents')
    expect(() => isValidObjectId('')).toThrow()
  })
})

describe('findAll', () => {
  it('maps GridFS files to HostedDocument shape', async () => {
    const { findAll } = await import('@/lib/api/documents')
    const id = new ObjectId()
    mockToArray.mockResolvedValue([
      {
        _id: id,
        filename: 'test.pdf',
        metadata: { contentType: 'application/pdf' },
        uploadDate: new Date('2024-01-01'),
      },
    ])

    const docs = await findAll()
    expect(docs).toHaveLength(1)
    expect(docs[0]).toMatchObject({
      id: id.toHexString(),
      filename: 'test.pdf',
      contentType: 'application/pdf',
      fileBase64: '[stored]',
      url: `http://localhost/api/v2/documents/${id.toHexString()}`,
    })
  })

  it('falls back to application/octet-stream when contentType is missing', async () => {
    const { findAll } = await import('@/lib/api/documents')
    mockToArray.mockResolvedValue([
      { _id: new ObjectId(), filename: 'unknown.bin', metadata: null, uploadDate: new Date() },
    ])

    const docs = await findAll()
    expect(docs[0].contentType).toBe('application/octet-stream')
  })
})

describe('findOne', () => {
  it('returns null when no file is found', async () => {
    const { findOne } = await import('@/lib/api/documents')
    mockToArray.mockResolvedValue([])
    const result = await findOne(new ObjectId().toHexString())
    expect(result).toBeNull()
  })

  it('returns a HostedDocument when the file exists', async () => {
    const { findOne } = await import('@/lib/api/documents')
    const id = new ObjectId()
    mockToArray.mockResolvedValue([
      {
        _id: id,
        filename: 'report.pdf',
        metadata: { contentType: 'application/pdf' },
        uploadDate: new Date('2024-06-01'),
      },
    ])

    const doc = await findOne(id.toHexString())
    expect(doc).not.toBeNull()
    expect(doc?.filename).toBe('report.pdf')
    expect(doc?.fileBase64).toBe('[stored]')
  })
})

describe('insert', () => {
  it('stores the document, clears fileBase64, and sets url + shortLink', async () => {
    const { insert } = await import('@/lib/api/documents')
    const { Writable } = await import('stream')
    const fakeId = new ObjectId()
    // Use a real Writable so Readable.pipe() works correctly
    const mockWritable = new Writable({ write(_chunk, _enc, cb) { cb() } }) as ReturnType<typeof mockOpenUploadStream>
    mockWritable.id = fakeId
    mockOpenUploadStream.mockReturnValue(mockWritable)

    const doc = {
      id: '',
      filename: 'sample.pdf',
      contentType: 'application/pdf',
      fileBase64: Buffer.from('fake pdf content').toString('base64'),
      url: '',
    }

    const result = await insert(doc)
    expect(result?.id).toBe(fakeId.toHexString())
    expect(result?.fileBase64).toBe('[stored]')
    expect(result?.url).toBe(`http://localhost/api/v2/documents/${fakeId.toHexString()}`)
    expect(result?.shortLink).toBe('https://tinyurl.com/abc')
  })
})
