/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { HostedDocument } from '@/lib/api/documents'
import Link from '@/components/custom-link'
import DocumentCard from '@/components/document-card'

interface DocumentListProps {
  documents: HostedDocument[]
}

export default function DocumentList({ documents }: DocumentListProps) {
  return (
    <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
      {documents.map((doc) => (
        <Link href={`/documents/${doc.id}`} key={doc.id}>
          <DocumentCard document={doc} />
        </Link>
      ))}
    </div>
  )
}
