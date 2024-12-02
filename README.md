# Document Host

![Release](https://img.shields.io/github/v/tag/robertwtucker/document-host?label=release)
![License](https://img.shields.io/github/license/robertwtucker/document-host)
![Open issues](https://img.shields.io/github/issues-raw/robertwtucker/document-host)
![Open pull requests](https://img.shields.io/github/issues-pr-raw/robertwtucker/document-host)

The Document Host (Docuhost) service provides a REST endpoint to upload demo-generated
documents for temporary storage. Documents can be retrieved via the short link
returned in the upload response.

## Getting Started

### Prerequisites

- Kubernetes 1.21+
- Helm 3.1+
- MongoDB 6.0+
- [Auth0](https://auth0.com) account (v0.4.0+)

The application is designed to be installed in [Kubernetes](https://kubernetes.io)
using a [Helm](https://helm.sh) chart. Files are stored in [MongoDB](https://www.mongodb.com)
using [GridFS](https://mongodb.com/docs/manual/core/gridfs/) so, prior to deployment,
an instance of MongoDB 6.0.x is required (tested with v6.0.13).

### Installation

A chart is available via the [SPT Chart Library](https://github.com/robertwtucker/spt-charts/docuhost).
See the [instructions](https://github.com/robertwtucker/spt-charts#getting-started)
in that repository for steps to clone the library.

To install the chart with the release name `docuhost`:

```bash
$ cd charts
$ helm upgrade --install docuhost ./docuhost \
    --namespace=spt-prod \
    --set db.user=admin \
    --set db.password=s3cr3t \
    ... # See chart for other configuration values
```

These commands deploy Docuhost to the Kubernetes cluster in the `spt-prod`
namespace. The parameters for the database username and password are set to
`admin` and `s3cr3t`, respectively. See the [Parameters](https://github.com/robertwtucker/spt-charts/tree/master/docuhost#parameters)
section of the [README](https://github.com/robertwtucker/spt-charts/tree/master/docuhost)
for a list of the parameters that can/need to be configured for a successful installation.

### Usage

Since v0.4.0, adding documents to Docuhost is secured using a Bearer token in the
HTTP Authorization header. The current implementation uses [Auth0](https://auth0.com)
as the sole authentication provider. A future version may use [Auth.js](https://authjs.dev)
(once their v5 library has been released) to allow for configurable authentication
providers.

To set up authentication in Auth0:

1. Under the _Applications_ menu, create a new _API_ with the name of your
   choosing using the identifier `urn:docuhost`.
2. In the _Permissions_ tab, add a `create:documents` permission.
3. Create a new Machine to Machine type _Application_ for your API making sure
   to select the permission created in the previous step.
4. Follow the instructions in the _QuickStart_ tab to obtain and use the token.

Sample requests and responses are provided below.

#### v0.4+

```json
# Request

Authorization: Bearer <API-Token>
HTTP POST /api/v2/documents
{
  "filename": "simple.pdf",
  "contentType": "application/pdf",
  "fileBase64": "JVBERi0xLjc[...]zd="
}

# Response

{
  "document": {
    "id": "61f0023ee260d827b7156c55",
    "filename": "simple.pdf",
    "contentType": "application/pdf",
    "url": "http://docuhost.localdev/api/v2/documents/61f0023ee260d827b7156c55",
    "shortLink": "https://tinyurl.com/k4ius98"
  }
}
```

#### v0.3.21 and lower

```json
# Request

HTTP POST /v1/documents
{
  "filename": "simple.pdf",
  "contentType": "application/pdf",
  "fileBase64": "JVBERi0xLjc[...]zd="
}

# Response

{
  "id": "61f0023ee260d827b7156c55",
  "filename": "simple.pdf",
  "contentType": "application/pdf",
  "url": "http://docuhost.localdev/v1/documents/61f0023ee260d827b7156c55",
  "shortLink": "https://tiny.one/yckaxkhv"
}
```

## Roadmap

See the [open issues](https://github.com/robertwtucker/document-host/issues) for
a list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to
learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Copyright (c) 2022 Quadient Group AG and distributed under the MIT License. See
`LICENSE` for more information.

## Contact

Robert Tucker - [@robertwtucker](https://x.com/robertwtucker)

Project Link: [https://github.com/robertwtucker/document-host](https://github.com/robertwtucker/document-host)
