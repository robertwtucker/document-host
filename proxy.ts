/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

export { auth as proxy } from '@/auth'

// Read more: https://nextjs.org/docs/app/guides/upgrading/version-16#middleware-renamed-to-proxy
export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
}
