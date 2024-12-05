/**
 * Copyright (c) 2024 Quadient Group AG
 * SPDX-License-Identifier: MIT
 */

import {
  FileArchiveIcon,
  FileJsonIcon,
  FileSlidersIcon as FilePresentationIcon,
  FileQuestionIcon,
  FileSpreadsheetIcon,
  FileTextIcon,
  FileUpIcon,
  FileCodeIcon as FileXmlIcon,
} from 'lucide-react'

interface FileIconProps {
  contentType: string
  [key: string]: any
}

export default function FileIcon({ contentType, ...props }: FileIconProps) {
  switch (contentType) {
    case 'application/pdf':
      return <FileUpIcon {...props} />
    case 'text/plain':
    case 'application/rtf':
    case 'application/msword':
    case 'application/vnd.openxmlformats-officedocument.wordprocessingml.document':
      return <FileTextIcon {...props} />
    case 'application/vnd.ms-excel':
    case 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet':
      return <FileSpreadsheetIcon {...props} />
    case 'application/vnd.ms-powerpoint':
    case 'application/vnd.openxmlformats-officedocument.presentationml.presentation':
      return <FilePresentationIcon {...props} />
    case 'application/json':
      return <FileJsonIcon {...props} />
    case 'application/xml':
      return <FileXmlIcon {...props} />
    case 'application/zip':
    case 'application/x-rar-compressed':
    case 'application/x-zip-compressed':
      return <FileArchiveIcon {...props} />
    default:
      return <FileQuestionIcon {...props} />
  }
}
