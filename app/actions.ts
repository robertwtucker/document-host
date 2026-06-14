/**
 * Copyright (c) 2026 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

'use server'

import { revalidatePath } from 'next/cache'

import { remove } from '@/lib/api/documents'
import { sessionHasPermission } from '@/lib/authorize'
import { logger } from '@/lib/logger'

export async function deleteDocument(id: string): Promise<void> {
  if (!(await sessionHasPermission('delete:documents'))) {
    logger.warn(`Unauthorized deleteDocument attempt for ${id}`)
    throw new Error('Forbidden')
  }
  await remove(id)
  revalidatePath('/')
}
