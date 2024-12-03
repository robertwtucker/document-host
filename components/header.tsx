/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import Link from '@/components/custom-link'
import UserButton from '@/components/user-button'

export default async function Header() {
  return (
    <header className="sticky flex justify-center bg-primary text-primary-foreground py-4">
      <div className="flex w-full mx-auto items-center justify-between px-4">
        <Link href="/" className="text-2xl font-bold">
          <h1>SPT DocuHost</h1>
        </Link>
        <UserButton />
      </div>
    </header>
  )
}
