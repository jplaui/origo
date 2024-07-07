# ORIGO CMD Toolkit


#### Contents
- [Introduction](#introduction)
- [Disclaimer](#disclaimer)
- [Structure and Development](#structure-and-development)
- [Command Line Toolkit](#command-line-toolkit)
- [Evaluation](#evaluation)
- [Limitations](#limitations)
- [Citing](#citing)


## Introduction
_ORIGO_ is a decentralized oracle that doesn't rely on the trusted hardware and two-party computation. It allows users to prove the *provenance* of TLS records from a specific API and, in the meanwhile, prove the dynamic content fulfills a value constraint as pre-defined in a public policy in zero-knowledge. It achieves confidentiality and integrity under certain honesty assumptions (e.g. clients cannot mount a MITM attack between the proxy and the server) which we define in the associated paper [ORIGO](https://eprint.iacr.org/2024/447.pdf). _ORIGO_ relies on a proxy sitting between the client user and the server hosting the API. The TLS stack running on the server does not require any modifications, making _ORIGO_ extendable to various existent APIs and web resources.


## Disclaimer
_ORIGO_ is a research project and its implementation has not undergone any security audit. Do not use _ORIGO_ in a production system or use it to process sensitive confidential data.


## Structure and Development
The structure of the project is as follows:
```
    ├── certs		# self-signed certificates for local development
    ├── client		# logic of http client to perform an API requests and prove policy compliance of response value
    ├── docs		# documentation files
    ├── proxy		# tcp tunneling service for tls, record layer traffic capturing/parsing, and policy verification
    ├── server		# local https instance to serve sample traffic
```

We quickly outline the key ideas behind the struture of the repository. For local development, developers can run the script `certs/cert.sh` to automatically generate self-signed root, server, and client certificates. The `server` folder and server-commands are only required _ORIGO_ is executed entirely locally.
The `proxy` folder uses TLS SNI to establish TCP connections to the correct destination server by parsing the ClientHello of incoming TLS sessions. Any TLS traffic between the client and the server is forwarded/reverse tunneled. The `proxy` verifies the correctness of public TLS session parameters of the client (communication of these parameters must be added to the repo if the repo is deployed in a WAN setting! we fixed it in our [opex-repositories](https://github.com/opex-research/tls-oracle-demo)). Next, the `proxy` runs the ZKP setup proceedure and calls ZKP.Verify if the proof bytes exist.
The `client` connects via TLS 1.3 to the server and uses the destination in SNI to help the proxy in determining the right connection. Next the client sends a request and receives a response in the TLS record phase. Afterwards the client runs the witness generation in two phases: 1. the key derivation circuit (KDC) phase and 2. the record phase. Last, the client computes the proof once setup parameters of the ZKP circuit are available.


## Command Line Toolkit
_ORIGO_ is a command-line toolkit which allows users to (i) generate input data for policy compliant data provenance proofs from private APIs, and (ii) prove policy-compliant proofs in zero-knowledge. The high-level workflow of the _ORIGO_ command-line toolkit is as follows:

### How to run the repository
Make sure to follow the [installation instructions](./docs/00_installation.md) to correctly clone and setup all required packages. Details on different deployments can be found in our [tutorial](./docs/tutorials) guidelines.

#### Running the protocol locally:
Notice that you can add the `-debug` flag to see further outputs.

1. (cd server) start the server service `go run main.go`, you can optionally use the `-largeresponse` flag to exchange multiple records
2. (cd proxy) start the proxy service `go run main.go -listen`
3. (cd client) set up a TLS connection with the server by tunneling through the proxy and perform a request `go run main.go -request`
4. (cd client) extract public inputs for the zkp computation of the tls handshake phase `go run main.go -postprocess-kdc`
5. (cd client) extract witness and public inputs for the zkp computation of the tls record phase `go run main.go -postprocess-record`
6. (cd proxy) postprocess captured traffic transcript and verify parameters from client to confirm zkp public input `go run main.go -postprocess`
7. (cd proxy) compute zkp setup parameters (proving key, verifier key, compiled constraint system of circuit, etc..) `go run main.go -debug -setup`
8. (cd client) run the zkp proof computation with the full witness `go run main.go -prove`
9. (cd proxy) run the zkp verification algorithm with the confirmed public inputs `go run main.go -debug -verify`
10. (cd proxy) output some benchmarks `go run main.go -debug -stats`


## Evaluations
The benchmarks found in the [research paper](https://eprint.iacr.org/2024/447.pdf) have been collected with a WAN deployment of this _ORIGO_ repository. You can find the WAN deployment of this repository in our repository [tls-oracle-demo](https://github.com/opex-research/tls-oracle-demo).


## Limitations
* _ORIGO_ currently supports off-chain zkp verification only and we plan to address on-chain zkp verification in the future. Take a look at the [portal repository](https://github.com/jplaui/portal) to access Groth16 and PLONK on-chain proof verifications. 
* We suggest to rely in the policy-to-circuit transpiler [zkGen](https://github.com/jplaui/zkGen) to automatically generate _ORIGO_ ZKP circuit according to flexibly configured policies (e.g. local or paypal).
* This sample implementation of the _ORIGO_ protocol works with TLS 1.3 configured with the cipher suite `TLS_AES_128_GCM_SHA256`. Additional cipher suite support must be added.


## Citing
We welcome you to cite our [research paper](https://eprint.iacr.org/2024/447.pdf) if you are using _ORIGO_ in academic research.
```
@inproceedings{ernstberger2024origo,
    author = {Ernstberger, Jens and Lauinger, Jan and Wu, Yinnan and Gervais, Arthur and Steinhorst, Sebastian},
    title = {ORIGO: Proving Provenance of Sensitive Data},
    year = {2024},
    publisher = {Cryptology ePrint Archive}
}
```

