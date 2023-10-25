/**
 * SPDX-License-Identifier: MIT
 */

import { NextRequest, NextResponse } from 'next/server'
import { verifyToken } from '@/lib/jwt'
import { insert } from '@/lib/api/document'
import { shorten } from '@/lib/shortlink'

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
        document.url = `${process.env.APP_URL}/${document.id}`
        const shortened = await shorten(document.url)
        if (shortened && shortened.shortlink) {
          document.shortLink = shortened.shortlink
        }

        let versionedResponse = {}
        if (version === 'v1') {
          versionedResponse = document
        } else {
          versionedResponse = { document }
        }

        return NextResponse.json(versionedResponse, {
          headers: {
            'Content-Type': 'application/json',
            Location: document.url,
          },
          status: 201,
        })
      } else {
        return new NextResponse(null, { status: 500 })
      }
    } else {
      return new NextResponse(null, { status: 401 })
    }
  } else {
    return new NextResponse(null, { status: 400 })
  }
}
