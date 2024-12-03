/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import Link from 'next/link'
import { notFound } from 'next/navigation'
import {
  FileIcon,
  FileIcon as FilePresentationIcon,
  FileSpreadsheetIcon,
  FileTextIcon,
} from 'lucide-react'

import { documents } from '@/lib/mock-data'
import { Button } from '@/components/ui/button'

export default function DocumentPage({ params }: { params: { id: string } }) {
  const document = documents.find((doc) => doc.id === params.id)

  if (!document) {
    notFound()
  }

  const getFileIcon = (fileType: string) => {
    switch (fileType) {
      case 'pdf':
        return <FileIcon className="w-8 h-8" />
      case 'docx':
        return <FileTextIcon className="w-8 h-8" />
      case 'xlsx':
        return <FileSpreadsheetIcon className="w-8 h-8" />
      case 'pptx':
        return <FilePresentationIcon className="w-8 h-8" />
      default:
        return <FileIcon className="w-8 h-8" />
    }
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <Link href="/">
        <Button variant="outline" className="mb-4">
          Back to Documents
        </Button>
      </Link>
      <div className="bg-white shadow-md rounded-lg p-6">
        <div className="flex items-center gap-4 mb-4">
          {getFileIcon(document.contentType)}
          <h1 className="text-3xl font-bold">{document.filename}</h1>
        </div>
        <p className="text-gray-600 mb-4">{document.url}</p>
        {/* <div className="grid grid-cols-2 gap-4 text-sm text-gray-500">
          <p>Created: {new Date(document.createdAt).toLocaleString()}</p>
          <p>Updated: {new Date(document.updatedAt).toLocaleString()}</p>
          <p>File Type: {document.fileType.toUpperCase()}</p>
        </div> */}
      </div>
    </div>
  )
}
