/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { NextRequest, NextResponse } from 'next/server'
import { auth } from '@/auth'

import { insert } from '@/lib/api/documents'
import { hasPermission, tokenFromRequest, verifyToken } from '@/lib/jwt'
import { logger } from '@/lib/logger'

export async function POST(request: NextRequest, context: { params: { version: string } }) {
  return auth(async (req: any & { auth?: { accessToken?: string } }) => {
    const { version } = context.params
    const requestInfo = `${request.method} ${request.nextUrl.pathname}`

    let authorized = false
    if (req?.auth?.accessToken) {
      logger.debug('User authenticated', req.auth)
      authorized = hasPermission(req.auth.accessToken, 'create:documents')
    } else {
      logger.debug('User not authenticated, checking for token')
      const verifiedToken = await verifyToken(req)
      if (verifiedToken) {
        authorized = hasPermission(tokenFromRequest(req), 'create:documents')
      }
    }

    if (authorized) {
      logger.debug('User has permission to create documents')
      if (version && version.match(new RegExp('^v[1-2]'))) {
        const payload = await req.json()
        const document = await insert(payload)
        if (document) {
          let versionedResponse = {}
          if (version === 'v1') {
            versionedResponse = document
          } else {
            versionedResponse = { document }
          }

          logger.info(requestInfo, { status: 201 })
          return NextResponse.json(versionedResponse, {
            headers: {
              'Content-Type': 'application/json',
              Location: document.url,
            },
            status: 201,
          })
        } else {
          logger.error('Failed to insert document into the database')
          logger.info(requestInfo, { status: 500 })
          return NextResponse.json({ error: 'Document insert failed' }, { status: 500 })
        }
      } else {
        logger.info(`${req.method} ${req.nextUrl.pathname}`, { status: 400 })
        return NextResponse.json({ error: 'Bad Request' }, { status: 400 })
      }
    } else {
      logger.info(requestInfo, { status: 401 })
      return NextResponse.json({ error: 'Unauthorized' }, { status: 401 })
    }
  })(request, context) as Promise<Response>
}
