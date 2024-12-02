import { HostedDocument } from './api/document/types'

export const documents: HostedDocument[] = [
  {
    id: '1',
    filename: 'Project Proposal',
    url: 'Detailed proposal for the new client project',
    fileBase64: '',
    // createdAt: "2023-05-01T10:00:00Z",
    // updatedAt: "2023-05-02T14:30:00Z",
    contentType: 'pdf',
  },
  {
    id: '2',
    filename: 'Financial Report Q2',
    url: 'Quarterly financial report for Q2 2023',
    fileBase64: '',
    // createdAt: "2023-07-15T09:00:00Z",
    // updatedAt: "2023-07-16T11:45:00Z",
    contentType: 'xlsx',
  },
  {
    id: '3',
    filename: 'Meeting Minutes',
    url: 'Minutes from the last team meeting',
    fileBase64: '',
    // createdAt: "2023-08-01T15:00:00Z",
    // updatedAt: "2023-08-01T16:30:00Z",
    contentType: 'docx',
  },
  {
    id: '4',
    filename: 'Product Roadmap',
    url: 'Product development roadmap for the next 12 months',
    fileBase64: '',
    // createdAt: "2023-06-10T11:00:00Z",
    // updatedAt: "2023-06-12T09:15:00Z",
    contentType: 'pptx',
  },
  {
    id: '5',
    filename: 'User Research Results',
    url: 'Results and analysis from the recent user research study',
    fileBase64: '',
    // createdAt: "2023-07-20T13:00:00Z",
    // updatedAt: "2023-07-22T10:30:00Z",
    contentType: 'pdf',
  },
]
