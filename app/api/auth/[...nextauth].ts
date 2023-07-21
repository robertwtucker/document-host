/**
 * Copyright (c) 2023 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import NextAuth, { NextAuthOptions } from 'next-auth'
import Auth0Provider from 'next-auth/providers/auth0'

export const authOptions: NextAuthOptions = {
  providers: [
    Auth0Provider({
      clientId: process.env.AUTH0_ID ?? '',
      clientSecret: process.env.AUTH0_SECRET ?? '',
      issuer: process.env.AUTH0_ISSUER ?? '',
    }),
  ],
}

export default NextAuth(authOptions)
