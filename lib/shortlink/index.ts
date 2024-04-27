/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { createShortlink } from './tinyurl'
// import { createShortlink } from './bitly'

export interface ShortlinkResponse {
  url: string
  shortlink?: string
}

export async function shorten(url: string): Promise<ShortlinkResponse> {
  return await createShortlink(url)
}
