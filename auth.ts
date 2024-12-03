/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import NextAuth from 'next-auth'

import 'next-auth/jwt'

import Auth0 from 'next-auth/providers/auth0'

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [
    Auth0({
      authorization: {
        params: {
          audience: process.env.AUTH_AUTH0_AUDIENCE,
        },
      },
    }),
  ],
  basePath: '/auth',
  session: { strategy: 'jwt' },
  callbacks: {
    async jwt({ token, trigger, session, account }) {
      if (account?.token_type === 'bearer') {
        return { ...token, accessToken: account.access_token }
      }
      return token
    },
    async session({ session, token }) {
      if (token?.accessToken) {
        session.accessToken = token.accessToken
      }
      return session
    },
  },
  debug: process.env.AUTH_DEBUG === 'true',
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
