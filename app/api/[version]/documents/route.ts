/**
 * SPDX-License-Identifier: MIT
 */

import { NextRequest, NextResponse } from 'next/server'
import { verifyToken } from '@/lib/jwt'
import { insert } from '@/lib/api/document'

type Params = {
  version: string
}

export async function POST(req: NextRequest, context: { params: Params }) {
  const { version } = context.params
  if (version && version.match(new RegExp('^v[1-2]'))) {
    const token = await verifyToken(req)
    if (token && token.scope?.includes('create:documents')) {
      const payload = await req.json()
      const document = await insert(payload)
      if (document) {
        return NextResponse.json(
          { document },
          {
            headers: {
              'Content-Type': 'application/json',
              Location: `/api/v2/documents/${document.id}`,
            },
            status: 201,
          }
        )
      } else {
        return new NextResponse(null, { status: 400 })
      }
    } else {
      return new NextResponse(null, { status: 401 })
    }
  } else {
    return new NextResponse(null, { status: 400 })
  }
}
