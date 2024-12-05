/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { createShortlink } from './tinyurl'

// import { createShortlink } from './bitly'

/**
 * The response from the shortlink service
 */
export interface ShortlinkResponse {
  url: string
  shortlink?: string
}

/**
 * @param url The URL to shorten
 * @returns A shortened version of the URL from the configured shortlink service
 */
export async function shorten(url: string): Promise<ShortlinkResponse> {
  return await createShortlink(url)
}
