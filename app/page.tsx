/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { auth } from '@/auth'

import { documents } from '@/lib/mock-data'
import DocumentList from '@/components/document-list'
import SearchBar from '@/components/search-bar'

export default async function Home() {
  const session = await auth()

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Documents</h1>
      {session ? (
        <div>
          <SearchBar />
          <DocumentList documents={documents} />
        </div>
      ) : (
        <div className="text-center text-muted-foreground">
          <p>
            Please <em>Sign In</em> to view documents.
          </p>
        </div>
      )}
    </div>
  )
}
