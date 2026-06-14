/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { HostedDocument } from '@/lib/api/documents'
import DocumentCard from '@/components/document-card'

interface DocumentListProps {
  documents: HostedDocument[]
  canDelete?: boolean
}

export default function DocumentList({ documents, canDelete = false }: DocumentListProps) {
  return (
    <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
      {documents.map((doc) => (
        <DocumentCard key={doc.id} document={doc} canDelete={canDelete} />
      ))}
    </div>
  )
}
