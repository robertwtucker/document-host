/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import clientPromise from '@/lib/mongodb'
import { GridFSBucket, ObjectId } from 'mongodb'
import { Readable } from 'stream'
import { HostedDocument, HostedFile } from './types'
import { logger } from '@/lib/logger'

export * from './types'

export async function insert(document: HostedDocument): Promise<HostedDocument | null> {
  const client = await clientPromise
  const bucket = new GridFSBucket(client.db())
  const buffer = Buffer.from(document.fileBase64, 'base64')
  const writeStream = bucket.openUploadStream(document.filename, {
    metadata: { contentType: document.contentType },
  })
  Readable.from(buffer).pipe(writeStream)

  document.id = writeStream.id.toHexString()
  document.fileBase64 = '[stored]' // don't pass this back now that it's in the DB

  return document
}

export async function find(id: string): Promise<HostedFile | null> {
  const client = await clientPromise
  const bucket = new GridFSBucket(client.db())
  const objectId = ObjectId.createFromHexString(id)
  const files = await bucket.find({ _id: objectId }).toArray()
  if (files.length !== 1) {
    return null
  }

  const contentStream = bucket.openDownloadStream(objectId)
  if (contentStream) {
    const hostedFile: HostedFile = {
      filename: files[0].filename,
      content: contentStream,
      contentType: files[0].metadata?.contentType ?? 'application/octet-stream',
      size: files[0].length,
    }

    return hostedFile
  }

  logger.warn(`No content stream found for ${id}`)
  return null
}

export function isValidObjectId(id: string): string {
  return ObjectId.createFromHexString(id).toHexString()
}
