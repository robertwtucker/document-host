/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { MongoClient, ServerApiVersion } from 'mongodb'

import { logger } from '@/lib/logger'

const config = {
  protocol: process.env.MONGODB_PROTOCOL || 'mongodb',
  host: process.env.MONGODB_HOST || 'localhost',
  port: process.env.MONGODB_PORT || '',
  database: process.env.MONGODB_DATABASE || 'documents',
  username: encodeURIComponent(process.env.MONGODB_USERNAME || 'docuhost'),
  password: encodeURIComponent(process.env.MONGODB_PASSWORD || ''),
  options: process.env.MONGODB_OPTIONS || '',
}

let uri = `${config.protocol}://${config.username}:${config.password}@${config.host}`
uri += `${config.port.length > 0 ? ':' + config.port : ''}/${config.database}`
uri += `${config.options.length > 0 ? '?' + config.options : ''}`

const options = {
  serverApi: {
    version: ServerApiVersion.v1,
    strict: true,
    deprecationErrors: true,
  },
}

let client
let clientPromise: Promise<MongoClient>

if (process.env.NODE_ENV === 'development') {
  // In development mode, use a global variable so that the value
  // is preserved across module reloads caused by HMR (Hot Module Replacement).
  const globalWithMongo = global as typeof globalThis & {
    _mongoClientPromise?: Promise<MongoClient>
  }
  if (
    typeof globalWithMongo._mongoClientPromise === undefined ||
    !globalWithMongo._mongoClientPromise
  ) {
    logger.info('Creating a global MongoClient instance...')
    client = new MongoClient(uri, options)
    logger.debug('Connecting to the database...')
    globalWithMongo._mongoClientPromise = client.connect()
  }
  clientPromise = globalWithMongo._mongoClientPromise
} else {
  logger.debug('Creating a new MongoClient instance...')
  client = new MongoClient(uri, options)
  logger.debug('Connecting to the database...')
  clientPromise = client.connect()
}

// Export a module-scoped MongoClient promise. By doing this in a
// separate module, the client can be shared across functions.
export default clientPromise
