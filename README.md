# Document Host

The Document Host (Docuhost) service provides a REST endpoint to upload demo-generated documents for temporary storage. Documents can be retrieved via the short link returned in the upload response.

## Getting Started

### Prerequisites

* Kubernetes 1.21+
* Helm 3.1+
* MongoDB 4.4+

The application is designed to be installed in [Kubernetes](https://kubernetes.io) using a [Helm](https://helm.sh) chart. Files are stored in [MongoDB](https://www.mongodb.com) using [GridFS](https://docs.mongodb.com/v4.4/core/gridfs) so, prior to deployment, an instance of MongoDB 4.4.x is required (tested with v4.4.12).

### Installation

The chart is available via the [SPT Chart Library](https://github.com/robertwtucker/spt-charts). See the [instructions](https://github.com/robertwtucker/spt-charts#getting-started) in that repository for steps to clone the library.

To install the chart with the release name `docuhost`:

``` bash
$ cd charts
$ helm upgrade --install docuhost ./docuhost --namespace=demo-prod \
    --set db.user=admin,db.password=s3cr3t
```

These commands deploy Docuhost to the Kubernetes cluster in the `demo-prod` namespace. The parameters for the database username and password are set to `admin` and `s3cr3t`, respectively. See the [Parameters](https://github.com/robertwtucker/spt-charts/tree/master/docuhost#parameters) section of the [README](https://github.com/robertwtucker/spt-charts/tree/master/docuhost) for a list of the parameters that can/need to be configured for a successful installation.

### Usage

A sample request and response are provided below.

#### Sample Request

``` json
HTTP POST /v1/documents
{
  "filename": "simple.pdf",
  "contentType": "application/pdf",
  "fileBase64": "JVBERi0xLjc[...]zd="
}
```

#### Sample Response

``` json
{
  "id": "61f0023ee260d827b7156c55",
  "filename": "simple.pdf",
  "contentType": "application/pdf",
  "url": "http://docuhost.localdev/v1/documents/61f0023ee260d827b7156c55",
  "shortLink": "https://tiny.one/yckaxkhv"
}
```

## Roadmap

See the [open issues](https://github.com/robertwtucker/document-host/issues) for a list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Copyright (c) 2022 Quadient Group AG and distributed under the MIT License. See `LICENSE` for more information.

## Contact

Robert Tucker - [@robertwtucker](https://twitter.com/robertwtucker)

Project Link: [https://github.com/robertwtucker/document-host](https://github.com/robertwtucker/document-host)
