/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import {
  FileIcon,
  FileIcon as FilePresentationIcon,
  FileSpreadsheetIcon,
  FileTextIcon,
} from 'lucide-react'

import { HostedDocument } from '@/lib/api/document'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import Link from '@/components/custom-link'

interface DocumentListProps {
  documents: HostedDocument[]
}

export default function DocumentList({ documents }: DocumentListProps) {
  const getFileIcon = (fileType: string) => {
    switch (fileType) {
      case 'pdf':
        return <FileIcon className="w-6 h-6" />
      case 'docx':
        return <FileTextIcon className="w-6 h-6" />
      case 'xlsx':
        return <FileSpreadsheetIcon className="w-6 h-6" />
      case 'pptx':
        return <FilePresentationIcon className="w-6 h-6" />
      default:
        return <FileIcon className="w-6 h-6" />
    }
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {documents.map((doc) => (
        <Link href={`/documents/${doc.id}`} key={doc.id}>
          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader className="flex flex-row items-center gap-4">
              {getFileIcon(doc.contentType)}
              <CardTitle>{doc.filename}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-gray-500 mb-2">{doc.url}</p>
              {/* <p className="text-xs text-gray-400">
                Updated: {new Date(doc.updatedAt).toLocaleDateString()}
              </p> */}
            </CardContent>
          </Card>
        </Link>
      ))}
    </div>
  )
}
