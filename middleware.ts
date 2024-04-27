import { NextRequest } from 'next/server'
import { verifyToken } from '@/lib/auth'

// Read more: https://nextjs.org/docs/app/building-your-application/routing/middleware#matcher
export const config = {
  matcher: ['/(.*)'],
}

export async function middleware(req: NextRequest) {
  const requestInfo = `${req.method} ${req.nextUrl.pathname}`
  const verifiedToken = await verifyToken(req).catch((err) => {
    console.error(`Failed to verify token: ${err}`)
  })

  // Are we POSTing to /api/v[1-2]/documents?
  if (req.method === 'POST' && req.nextUrl.pathname.match(new RegExp('/api/v[1-2]/documents'))) {
    if (verifiedToken && verifiedToken.scope?.includes('create:documents')) {
      return
    } else {
      // Logging not supported in Next.js Edge Runtime (middleware)
      console.log(
        JSON.stringify({
          level: 'info',
          message: requestInfo,
          status: 401,
          timestamp: new Date().toISOString(),
        })
      )
      return new Response(null, { status: 401 })
    }
  }

  // Forbid access to the root path for now (suppress default app rendering)
  if (req.method === 'GET' && req.nextUrl.pathname === '/') {
    // Logging not supported in Next.js Edge Runtime (middleware)
    console.log(
      JSON.stringify({
        level: 'info',
        message: requestInfo,
        status: 403,
        timestamp: new Date().toISOString(),
      })
    )
    return new Response(null, { status: 403 })
  }
}
