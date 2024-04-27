/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { ShortlinkResponse } from './'

type TinyUrlResponse = {
  data: {
    domain: string
    alias: string
    deleted: boolean
    archived: boolean
    tags: string[]
    analytics: [{ key: boolean }]
    tiny_url: string
    url: string
  }
  code: number
  errors: string[]
}

export async function createShortlink(url: string): Promise<ShortlinkResponse> {
  const serviceUrl = `${process.env.SHORTLINK_SERVICE_URL}/create`
  const emptyShortlinkResponse = { url: url }

  const headers = new Headers()
  headers.append('Accept', 'application/json')
  headers.append('Authorization', `Bearer ${process.env.SHORTLINK_API_KEY}`)
  headers.append('Content-Type', 'application/json')

  try {
    const response = await fetch(serviceUrl, {
      method: 'POST',
      headers: headers,
      body: JSON.stringify({ url: url, domain: process.env.SHORTLINK_DOMAIN }),
    })

    if (!response.ok) {
      console.error(`Failed to create shortlink: ${await response.text()}`)
      return emptyShortlinkResponse
    }

    const json: TinyUrlResponse = await response.json()
    return { url: json.data.url, shortlink: json.data.tiny_url }
  } catch (err) {
    console.error(`Failed to create shortlink: ${err}`)
    return emptyShortlinkResponse
  }
}
