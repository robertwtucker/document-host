/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'

vi.mock('@/lib/logger', () => ({
  logger: { debug: vi.fn(), error: vi.fn(), info: vi.fn(), warn: vi.fn() },
}))

const mockFetch = vi.fn()
vi.stubGlobal('fetch', mockFetch)

beforeEach(() => {
  vi.clearAllMocks()
  process.env.SHORTLINK_SERVICE_URL = 'https://api.tinyurl.com'
  process.env.SHORTLINK_API_KEY = 'test-key'
  process.env.SHORTLINK_DOMAIN = 'tinyurl.com'
})

describe('TinyURL createShortlink', () => {
  it('returns a shortened URL on success', async () => {
    const { createShortlink } = await import('@/lib/shortlink/tinyurl')
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({
        data: {
          domain: 'tinyurl.com',
          alias: 'abc123',
          deleted: false,
          archived: false,
          tags: [],
          analytics: [],
          tiny_url: 'https://tinyurl.com/abc123',
          url: 'https://example.com/api/v2/documents/123',
        },
        code: 0,
        errors: [],
      }),
    })

    const result = await createShortlink('https://example.com/api/v2/documents/123')
    expect(result.shortlink).toBe('https://tinyurl.com/abc123')
    expect(result.url).toBe('https://example.com/api/v2/documents/123')
  })

  it('returns the original URL (no shortlink) when the API call fails', async () => {
    const { createShortlink } = await import('@/lib/shortlink/tinyurl')
    mockFetch.mockResolvedValue({
      ok: false,
      text: async () => 'Unauthorized',
    })

    const result = await createShortlink('https://example.com/doc/456')
    expect(result.shortlink).toBeUndefined()
    expect(result.url).toBe('https://example.com/doc/456')
  })

  it('returns the original URL when fetch throws', async () => {
    const { createShortlink } = await import('@/lib/shortlink/tinyurl')
    mockFetch.mockRejectedValue(new Error('network error'))

    const result = await createShortlink('https://example.com/doc/789')
    expect(result.shortlink).toBeUndefined()
    expect(result.url).toBe('https://example.com/doc/789')
  })
})

describe('Bitly createShortlink', () => {
  beforeEach(() => {
    process.env.SHORTLINK_SERVICE_URL = 'https://api-ssl.bitly.com'
    process.env.SHORTLINK_DOMAIN = 'bit.ly'
  })

  it('returns a shortened URL on success', async () => {
    const { createShortlink } = await import('@/lib/shortlink/bitly')
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({
        created_at: '2024-01-01T00:00:00+0000',
        id: 'bit.ly/abc',
        link: 'https://bit.ly/abc',
        custom_bitlinks: [],
        long_url: 'https://example.com/api/v2/documents/123',
        archived: false,
        tags: [],
        deeplinks: [],
        references: { group: 'Bg1' },
      }),
    })

    const result = await createShortlink('https://example.com/api/v2/documents/123')
    expect(result.shortlink).toBe('https://bit.ly/abc')
    expect(result.url).toBe('https://example.com/api/v2/documents/123')
  })

  it('returns the original URL when the API call fails', async () => {
    const { createShortlink } = await import('@/lib/shortlink/bitly')
    mockFetch.mockResolvedValue({
      ok: false,
      text: async () => 'Rate limit exceeded',
    })

    const result = await createShortlink('https://example.com/doc/456')
    expect(result.shortlink).toBeUndefined()
    expect(result.url).toBe('https://example.com/doc/456')
  })
})
