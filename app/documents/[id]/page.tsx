/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import { notFound } from 'next/navigation'
import { auth } from '@/auth'

import { findOne } from '@/lib/api/documents'
import { hasPermission } from '@/lib/jwt'
import { Button } from '@/components/ui/button'
import Link from '@/components/custom-link'
import FileIcon from '@/components/file-icon'

export default async function DocumentPage(props: { params: Promise<{ id: string }> }) {
  const params = await props.params
  let canListDocuments = false
  // let canDeleteDocuments = false

  const session = await auth()
  if (session?.accessToken) {
    canListDocuments = hasPermission(session.accessToken, 'list:documents')
    // canDeleteDocuments = hasPermission(session.accessToken, 'delete:documents')
  }

  if (canListDocuments) {
    const document = await findOne(params.id)
    if (!document) {
      notFound()
    }

    return (
      <div className="container mx-auto px-4 py-8">
        <Link href="/">
          <Button variant="outline" className="mb-4">
            Back to Documents
          </Button>
        </Link>
        <div className="rounded-lg bg-white p-6 shadow-md">
          <div className="mb-4 flex items-center gap-4">
            <FileIcon contentType={document.contentType} className="size-8" />
            <h1 className="text-3xl font-bold">{document.filename}</h1>
          </div>
          <Link href={document.url} target="_blank">
            <p className="mb-4 text-gray-600">{document.url}</p>
          </Link>
          <div className="grid grid-cols-2 gap-4 text-sm text-gray-500">
            <p>Content Type: {document.contentType}</p>
            <p>
              Uploaded:{' '}
              {document?.uploadedAt ? new Date(document.uploadedAt).toLocaleString() : 'Unknown'}
            </p>
          </div>
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
