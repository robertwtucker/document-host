/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { ShortlinkResponse } from './'

type BitlyResponse = {
  created_at: string
  id: string
  link: string
  custom_bitlinks: string[]
  long_url: string
  archived: boolean
  tags: string[]
  deeplinks: string[]
  references: {
    group: string
  }
}

export async function createShortlink(url: string): Promise<ShortlinkResponse> {
  const serviceUrl = `${process.env.SHORTLINK_SERVICE_URL}/v4/shorten`
  const emptyShortlinkResponse = { url: url }

  const headers = new Headers()
  headers.append('Accept', 'application/json')
  headers.append('Authorization', `Bearer ${process.env.SHORTLINK_API_KEY}`)
  headers.append('Content-Type', 'application/json')

  try {
    const response = await fetch(serviceUrl, {
      method: 'POST',
      headers: headers,
      body: JSON.stringify({
        long_url: url,
        domain: process.env.SHORTLINK_DOMAIN,
      }),
    })

    if (!response.ok) {
      console.error(`Failed to create shortlink: ${await response.text()}`)
      return emptyShortlinkResponse
    }

    const json: BitlyResponse = await response.json()
    return { url: json.long_url, shortlink: json.link }
  } catch (err) {
    console.error(`Failed to create shortlink: ${err}`)
    return emptyShortlinkResponse
  }
}
