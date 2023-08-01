/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { ShortlinkResponse } from './types'
import { createShortlink } from './tinyurl'
// import { createShortlink } from './bitly'

export * from './types'

export async function shorten(url: string): Promise<ShortlinkResponse> {
  return await createShortlink(url)
}
