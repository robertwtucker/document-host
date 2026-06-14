/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { findAll, HostedDocument } from '@/lib/api/documents'
import { sessionHasPermission } from '@/lib/authorize'
import DocumentList from '@/components/document-list'
import SearchBar from '@/components/search-bar'

export default async function Home() {
  let documents: HostedDocument[] = []

  const canListDocuments = await sessionHasPermission('list:documents')

  if (canListDocuments) {
    documents = await findAll()
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="mb-8 text-3xl font-bold">Documents</h1>
        <div>
          <SearchBar />
          <DocumentList documents={documents} />
        </div>
      </div>
    )
  } else {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-muted-foreground text-center">
          <p>
            Please <em>Sign In</em> to view documents.
          </p>
        </div>
      </div>
    )
  }
}
