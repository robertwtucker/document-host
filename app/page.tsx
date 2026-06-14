/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { findAll } from '@/lib/api/documents'
import { sessionHasPermission } from '@/lib/authorize'
import DocumentList from '@/components/document-list'

export default async function Home() {
  const [canListDocuments, canDeleteDocuments] = await Promise.all([
    sessionHasPermission('list:documents'),
    sessionHasPermission('delete:documents'),
  ])

  if (canListDocuments) {
    const documents = await findAll()
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="mb-8 text-3xl font-bold">Documents</h1>
        {documents.length === 0 ? (
          <p className="text-muted-foreground text-center">No documents have been uploaded.</p>
        ) : (
          <DocumentList documents={documents} canDelete={canDeleteDocuments} />
        )}
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
