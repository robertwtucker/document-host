/**
 * SPDX-License-Identifier: MIT
 */

import { NextRequest, NextResponse } from 'next/server'
import { pingDb } from '@/lib/api/healthcheck'

type Params = {
  version: string
}

export async function GET(req: NextRequest, context: { params: Params }) {
  const { version } = context.params
  if (version && version.match(new RegExp('^v[1-2]'))) {
    try {
      await pingDb()
      return new NextResponse('OK', { status: 200 })
    } catch (err) {
      return new NextResponse(null, { status: 500 })
    }
  } else {
    return new NextResponse(null, { status: 400 })
  }
}
