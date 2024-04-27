/**
 * SPDX-License-Identifier: MIT
 */

import { NextRequest, NextResponse } from 'next/server'
import { pingDb } from '@/lib/api/healthcheck'
import { logger } from '@/lib/logger'

type Params = {
  version: string
}

export async function GET(req: NextRequest, context: { params: Params }) {
  const { version } = context.params
  const requestInfo = `${req.method} ${req.nextUrl.pathname}`

  if (version && version.match(new RegExp('^v[1-2]'))) {
    try {
      await pingDb()
      logger.info(requestInfo, { status: 200 })
      return new NextResponse('OK', { status: 200 })
    } catch (err) {
      logger.error(err)
      logger.info(requestInfo, { status: 500 })
      return new NextResponse(null, { status: 500 })
    }
  } else {
    logger.info(requestInfo, { status: 400 })
    return new NextResponse(null, { status: 400 })
  }
}
