/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { NextRequest, NextResponse } from 'next/server'
import { HostedFile, find, isValidObjectId } from '@/lib/api/document'

type Params = {
  version: string
  id: string
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

export async function GET(req: NextRequest, context: { params: Params }) {
  const { version, id } = context.params
  if (
    version.startsWith('v') &&
    (version.endsWith('1') || version.endsWith('2'))
  ) {
    try {
      isValidObjectId(id)
    } catch (err) {
      return new NextResponse(null, {
        status: 400,
      })
    }

    const hostedFile = await find(id)
    if (typeof hostedFile == undefined || !hostedFile) {
      return new NextResponse(null, { status: 404 })
    } else {
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
    return new NextResponse(null, { status: 400 })
  }
}
