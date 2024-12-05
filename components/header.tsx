/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import Link from '@/components/custom-link'
import UserButton from '@/components/user-button'

export default async function Header() {
  return (
    <header className="bg-primary text-primary-foreground sticky flex justify-center py-4">
      <div className="mx-auto flex w-full items-center justify-between px-4">
        <Link href="/" className="text-2xl font-bold">
          <h1>SPT DocuHost</h1>
        </Link>
        <UserButton />
      </div>
    </header>
  )
}
