import { documents } from '@/lib/mock-data'
import DocumentList from '@/components/document-list'
import SearchBar from '@/components/search-bar'

export default function Home() {
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Document Repository</h1>
      <SearchBar />
      <DocumentList documents={documents} />
    </div>
  )
}
