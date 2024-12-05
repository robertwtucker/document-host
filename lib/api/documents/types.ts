/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { GridFSBucketReadStream } from 'mongodb'

/**
 * Represents a file and its metadata stored in the database.
 */
export interface HostedDocument {
  id: string
  filename: string
  contentType: string
  fileBase64: string
  url: string
  shortLink?: string
  uploadedAt?: Date
}

/**
 * Represents the base properties of a file stored in the database.
 */
export interface HostedFile {
  filename: string
  content: GridFSBucketReadStream
  contentType: string
  size: number
}
