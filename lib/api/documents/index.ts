/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { Readable } from 'stream'
import { GridFSBucket, ObjectId } from 'mongodb'

import { logger } from '@/lib/logger'
import clientPromise from '@/lib/mongodb'
import { shorten } from '@/lib/shortlink'

import { HostedDocument, HostedFile } from './types'

export * from './types'

/**
 * @returns All documents stored in the database
 */
export async function findAll(): Promise<HostedDocument[]> {
  const client = await clientPromise
  const bucket = new GridFSBucket(client.db())
  const files = await bucket.find().toArray()
  return files.map((file) => {
    const id = file._id.toHexString()
    console.log('id:', id)
    return {
      id: id,
      filename: file.filename,
      contentType: file.metadata?.contentType ?? 'application/octet-stream',
      fileBase64: '[stored]',
      url: `${process.env.APP_URL}/${id}`,
      uploadedAt: file.uploadDate,
    } as HostedDocument
  })
}

/**
 * @param id The ID of the document to find
 * @returns The document with the specified ID, or null if not found
 */
export async function findOne(id: string): Promise<HostedDocument | null> {
  const client = await clientPromise
  const bucket = new GridFSBucket(client.db())
  const objectId = ObjectId.createFromHexString(id)
  const files = await bucket.find({ _id: objectId }).toArray()
  if (files.length !== 1) {
    return null
  }

  return {
    id: id,
    filename: files[0].filename,
    contentType: files[0].metadata?.contentType ?? 'application/octet-stream',
    fileBase64: '[stored]',
    url: `${process.env.APP_URL}/${id}`,
    uploadedAt: files[0].uploadDate,
  } as HostedDocument
}

/**
 * @param document The document to insert
 * @returns The inserted document with its ID and URL(s)
 */
export async function insert(document: HostedDocument): Promise<HostedDocument | null> {
  const client = await clientPromise
  const bucket = new GridFSBucket(client.db())
  const buffer = Buffer.from(document.fileBase64, 'base64')
  const writeStream = bucket.openUploadStream(document.filename, {
    metadata: { contentType: document.contentType.toLowerCase() },
  })
  Readable.from(buffer).pipe(writeStream)

  document.id = writeStream.id.toHexString()
  document.fileBase64 = '[stored]' // don't pass this back now that it's in the DB
  document.url = `${process.env.APP_URL}/${document.id}`
  const shortened = await shorten(document.url)
  if (shortened && shortened.shortlink) {
    document.shortLink = shortened.shortlink
  }

  return document
}

/**
 * @param id The ID of the file to fetch
 * @returns The file with the specified ID, or null if not found
 */
export async function fetch(id: string): Promise<HostedFile | null> {
  const client = await clientPromise
  const bucket = new GridFSBucket(client.db())
  const objectId = ObjectId.createFromHexString(id)
  const files = await bucket.find({ _id: objectId }).toArray()
  if (files.length !== 1) {
    return null
  }

  const contentStream = bucket.openDownloadStream(objectId)
  if (contentStream) {
    return {
      filename: files[0].filename,
      content: contentStream,
      contentType: files[0].metadata?.contentType ?? 'application/octet-stream',
      size: files[0].length,
    } as HostedFile
  }

  logger.warn(`No content stream found for ${id}`)
  return null
}

/**
 * @param id The ID value to test
 * @returns A stringified MongoDB ObjectId, if valid, otherwise it throws an error
 */
export function isValidObjectId(id: string): string {
  return ObjectId.createFromHexString(id).toHexString()
}
