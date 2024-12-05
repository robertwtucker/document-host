/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { HostedDocument } from '@/lib/api/documents'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import FileIcon from '@/components/file-icon'

interface DocumentCardProps {
  document: HostedDocument
}

export default function DocumentCard({ document }: DocumentCardProps) {
  return (
    <Card className="hover:shadow-lg transition-shadow">
      <CardHeader className="flex flex-row items-center gap-4">
        <FileIcon contentType={document.contentType} className="w-6 h-6" />
        <CardTitle>{document.filename}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-gray-500 mb-2">{document.contentType}</p>
        <p className="text-xs text-gray-400">
          Uploaded:{' '}
          {document.uploadedAt ? new Date(document.uploadedAt).toLocaleDateString() : 'Unknown'}
        </p>
      </CardContent>
    </Card>
  )
}
