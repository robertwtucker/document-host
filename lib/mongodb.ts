/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { MongoClient, ServerApiVersion } from 'mongodb'
import { logger } from '@/lib/logger'

const username = encodeURIComponent(process.env.MONGODB_USERNAME as string)
const password = encodeURIComponent(process.env.MONGODB_PASSWORD as string)
let uri = `mongodb://${username}:${password}@${process.env.MONGODB_HOST}:${process.env.MONGODB_PORT}/${process.env.MONGODB_DATABASE}`
username.toLowerCase() === 'root'
  ? (uri += `?authSource=admin`)
  : (uri += `?authSource=${process.env.MONGODB_DATABASE}`)
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
  let globalWithMongo = global as typeof globalThis & {
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
