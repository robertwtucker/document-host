import './globals.css'

import type { Metadata } from 'next'
import { Inter } from 'next/font/google'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Document Host',
  description: 'View and manage hosted documents',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <header className="bg-primary text-primary-foreground py-4">
          <div className="container mx-auto px-4">
            <h1 className="text-2xl font-bold">Document Host</h1>
          </div>
        </header>
        <main>{children}</main>
      </body>
    </html>
  )
}
