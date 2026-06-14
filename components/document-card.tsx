/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { Trash2 } from 'lucide-react'

import { deleteDocument } from '@/app/actions'
import { HostedDocument } from '@/lib/api/documents'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import FileIcon from '@/components/file-icon'

interface DocumentCardProps {
  document: HostedDocument
  canDelete?: boolean
}

function formatSize(bytes?: number): string | null {
  if (bytes === undefined || bytes === null) return null
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`
}

export default function DocumentCard({ document, canDelete = false }: DocumentCardProps) {
  const fileHref = `/api/v2/documents/${document.id}`
  const formattedSize = formatSize(document.size)

  return (
    <Card className="transition-shadow hover:shadow-lg">
      <CardHeader className="flex flex-row items-center gap-4">
        <FileIcon contentType={document.contentType} className="size-6 shrink-0" />
        <CardTitle className="grow">
          <a
            href={fileHref}
            target="_blank"
            rel="noopener noreferrer"
            className="hover:underline focus:underline focus:outline-none"
          >
            {document.filename}
          </a>
        </CardTitle>
        {canDelete && (
          <form action={deleteDocument.bind(null, document.id)}>
            <Button
              type="submit"
              variant="ghost"
              size="icon"
              aria-label={`Delete ${document.filename}`}
              title="Delete"
            >
              <Trash2 className="size-4" />
            </Button>
          </form>
        )}
      </CardHeader>
      <CardContent className="space-y-1">
        <p className="text-sm text-gray-500">{document.contentType}</p>
        <p className="text-xs text-gray-400">
          {formattedSize && <>Size: {formattedSize} · </>}
          Uploaded:{' '}
          {document.uploadedAt
            ? new Date(document.uploadedAt).toLocaleDateString()
            : 'Unknown'}
        </p>
        {document.shortLink && (
          <p className="text-xs">
            <a
              href={document.shortLink}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 hover:underline"
            >
              {document.shortLink}
            </a>
          </p>
        )}
      </CardContent>
    </Card>
  )
}
