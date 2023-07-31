/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import clientPromise from '@/lib/mongodb'

export async function pingDb(): Promise<Record<string, string> | null> {
  const client = await clientPromise
  return await client.db().admin().command({ ping: 1 })
}
