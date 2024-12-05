/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { auth } from '@/auth'

import { SignIn, SignOut } from '@/components/auth-buttons'

export default async function UserButton() {
  const session = await auth()
  if (!session?.user) return <SignIn />
  return <SignOut />
}
