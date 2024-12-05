/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { NextRequest, NextResponse } from 'next/server'

import { fetch, HostedFile, isValidObjectId } from '@/lib/api/documents'
import { logger } from '@/lib/logger'

type Params = {
  version: string
  id: string
}

export async function GET(req: NextRequest, context: { params: Promise<Params> }) {
  const { version, id } = await context.params
  const requestInfo = `${req.method} ${req.nextUrl.pathname}`

  if (version && version.match(new RegExp('^v[1-2]'))) {
    try {
      isValidObjectId(id)
    } catch (err) {
      logger.error(err)
      logger.info(requestInfo, { status: 400 })
      return new NextResponse(null, {
        status: 400,
      })
    }

    const hostedFile = await fetch(id)
    if (typeof hostedFile == undefined || !hostedFile) {
      logger.info(requestInfo, { status: 404 })
      return new NextResponse(null, { status: 404 })
    } else {
      logger.info(requestInfo, { status: 200 })
      return new NextResponse(streamData(hostedFile), {
        status: 200,
        headers: {
          'Content-Disposition': 'inline',
          'Content-Type': hostedFile.contentType,
          'Content-Length': hostedFile.size.toString(),
        },
      })
    }
  } else {
    logger.info(requestInfo, { status: 400 })
    return new NextResponse(null, { status: 400 })
  }
}

function streamData(file: HostedFile): ReadableStream<Uint8Array> {
  const downloadStream = file.content
  return new ReadableStream({
    start(controller) {
      downloadStream.on('data', (chunk) => {
        controller.enqueue(new Uint8Array(chunk))
      })
      downloadStream.on('end', () => {
        controller.close()
      })
    },
    cancel() {
      downloadStream.destroy()
    },
  })
}
