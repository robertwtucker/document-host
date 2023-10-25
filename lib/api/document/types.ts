/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { GridFSBucketReadStream } from 'mongodb'

/**
 *
 */
export interface HostedDocument {
  id: string
  filename: string
  contentType: string
  fileBase64: string
  url: string
  shortLink?: string
}

/**
 *
 */
export interface HostedFile {
  filename: string
  content: GridFSBucketReadStream
  contentType: string
  size: number
}
