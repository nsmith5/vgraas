# VGRaaS - Video Game Reviews as a Service

[![Build Status](https://cloud.drone.io/api/badges/nsmith5/vgraas/status.svg?branch=development)](https://cloud.drone.io/nsmith5/vgraas) [![Code Coverage](https://codecov.io/gh/nsmith5/vgraas/branch/development/graph/badge.svg)](https://codecov.io/gh/nsmith5/vgraas) [![Documentation](https://godoc.org/github.com/nsmith5/vgraas/pkg/vgraas?status.svg)](https://godoc.org/github.com/nsmith5/vgraas/pkg/vgraas)


Thats right, video game reviews as a service. When you're making a top quality
video game review blog you don't want to think about where you need to store
all those reviews and comments. Video game reviews as a service (vrgaas) is a
REST API for storing reviews and comments. You can host your own or use our 
hosted version at [https://vgraas.nfsmith.ca](https://vgraas.nfsmith.ca).

## Usage

The vgraas API specification is published as an 
[OpenAPI 3.0](https://swagger.io/specification/) specification. You can find
the specification in `specification.yaml` in the root of this repository. You
can also browse the specification on 
[SwaggerHub](https://app.swaggerhub.com/apis/nsmith5/vrgaas/0.1.1). Their UI
makes it especially easy to get a quick understanding of the API.

Users should note that the hosted version is rate limited at 5 requests per
second and limits uploads to 500 KiB.

## Installation & Self Hosting

If you'd like to host your own version of the API, or just want to test
it out without hitting the hosted version, deploying vgraas is a breeze.

**Docker**
```
# Runs vgraas on localhost port 8080
$ docker run -p 8080:8080 nsmith5/vgraas:0.1.1
```

**Release Binaries**
```
# Replace 'linux' with 'windows' or 'darwin' if that fits your needs
$ wget https://github.com/nsmith5/vgraas/releases/download/0.1.1/vgraas-linux-amd64.zip
$ unzip vgraas-linux-amd64.zip
$ ./vgraas -h
Usage of ./vgraas:
  -api string
        API listen address (default ":8080")
```

**Kubernetes**

So you want to deploy on Kubernetes? How very cool of you. There is a 
[Helm](https://helm.sh) chart in the `/chart` directory of the git repository
you can use to help out.

```
$ git clone https://github.com/nsmith5/vgraas.git
$ cd vgraas/chart
$ helm install --name vgraas .  # Modify values in values.yaml to configure install
```

If `ingress.expose` is false, the chart will deploy a ClusterIP service. This
works in a vanilla cluster. If you set `ingress.expose` to true, an ingress with
TLS termination will be created. For this to work you'll need to have an ingress
controller installed (I recommend 
[HAproxy](https://github.com/jcmoraisjr/haproxy-ingress)) and 
[cert-manager](https://github.com/jetstack/cert-manager).

## Hacking

Interested in contributing to vgraas? Awesome. To get started, clone the 
repository and install a [Go](https://golang.org) toolchain >=1.11.

```
$ git clone https://github.com/nsmith5/vgraas.git
$ cd vgraas
$ go run cmd/vgraas/*        # Build and run vgraas
$ go test -race -cover ./... # Run all units tests and check code coverage
```

The repo is laid out as follows:
```
$ tree
.
├── chart                       # Helm chart for kubernetes deployment
│   └── ...
├── cmd
│   └── vgraas                  # Main executable
│       └── ...
├── pkg                         # Where the libraries live (most of the code)
│   ├── middleware              # Misc middlewares for the API
│   │   └── ...
│   └── vgraas                  # Core logic (API and data model)
│       └── ...
├── README.md                   # You are here!
└── specification.yaml          # OpenAPI 3 specification

```

## Questions

If you have any questions about the repository email me at vgraas@nfsmith.ca!
