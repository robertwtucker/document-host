/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import NextAuth from 'next-auth'
import Auth0 from 'next-auth/providers/auth0'

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [Auth0({ authorization: { params: { audience: process.env.AUTH_AUTH0_AUDIENCE } } })],
})
