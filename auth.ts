/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import NextAuth from 'next-auth'

import 'next-auth/jwt'

import Auth0 from 'next-auth/providers/auth0'

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [Auth0({ authorization: { params: { audience: process.env.AUTH_AUTH0_AUDIENCE } } })],
  // providers: [Auth0],
  session: { strategy: 'jwt' },
  callbacks: {
    authorized({ request, auth }) {
      const { pathname } = request.nextUrl
      console.log('authorized for path:', pathname, ', session:', auth)
      return true
    },
    session({ session, token }) {
      // console.log('session callback', session, token)
      if (token?.accessToken) session.accessToken = token.accessToken
      return session
    },
  },
})

declare module 'next-auth' {
  interface Session {
    accessToken?: string
  }
}

declare module 'next-auth/jwt' {
  interface JWT {
    accessToken?: string
  }
}
