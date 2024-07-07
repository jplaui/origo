## DOCS of the modified TLS client library as a module for data capturing/parsing

#### latest update to include tls lib as a module
- download [crypto/tls](https://github.com/golang/go/tree/master/src/crypto/tls) and unzip in `client` folder
- download [crypto/ecdh](https://github.com/golang/go/tree/master/src/crypto/ecdh) and unzip in `client/tls` folder
- download [crypto/internal](https://github.com/golang/go/tree/master/src/crypto/internal) and unzip in `client/tls` folder

follow up local adjustments before usage
- in `cipher_suites`, change `crypto/internal/boring` to `client/tls/internal/boring`, and uncomment all `internal/cpu` calls (set `hasAESGCMHardwareSupport = false`)
- in `handshake_client`, `handshake_client_tls13`, `kee_agreement`, `key_schedule`, `handshake_server_test`, change `crypto/ecdh` to `client/tls/ecdh`
- in `tls/ecdh/ecdh.go`, `tls/ecdh/x25519.go`, `tls/ecdh/nist.go`, `tls/internal/boring/notboring.go` change all `crypto/internal/*` imports to `client/tls/internal/*`

- (not necessary if boring enabled and imported, if you comment out boring everywhere, make sure to remove that line.) remove line `//go:build !amd64 && !arm64 && !ppc64le && !s390x` inside `client/tls/internal/nistec/p256.go`

#### go modules on the client
- run `go get github.com/consensys/gnark@develop` for compiling and running the latest gnark backends

#### submodules
- to add a repo, jump into the location where it should be cloned, then run `git submodule add cloning-link`

#### how to run mpc library
- command of evaluator `MPCLDIR=~/Code/github.com/markkurossi/mpc ./garbled -v -e -i 109 examples/sha256.mpcl`
- command of garbler `MPCLDIR=~/Code/github.com/markkurossi/mpc ./garbled -v -i 116 examples/sha256.mpcl`

if calling from within the folder, then use `MPCLDIR=./../../ ./garbled -v -e -i 0 examples/sha256.mpcl`
and `MPCLDIR=./../../ ./garbled -v -i 2 examples/sha256.mpcl`

#### ectf computation
- file modified: client_handshake.go in tls_fork and new_ectf.go
- storage of ectf communication values in new_storage_ectf then client.json and proxy.json, (or just a single params.go where both parties read from?)

#### how can we parse data on the client side?
- for handshake traffic of TLS1.3, print out state of `clientHandshakeStateTLS13` which is part of the file `handshake_client_tls13.go`.
- inside `handshake_client_tls13.go`, the functions `establishHandshakeKeys` and `readServerFinished` are of interest to derive handshake traffic keys and application traffic keys (+ server finished data) respectively. 
- important: change the `defaultCipherSuitesTLS13NoAES` order in file `cipher_suites.go` to prioritize `TLS_AES_128_GCM_SHA256` (otherwise, polychacha is selected...), change to:
```go
var defaultCipherSuitesTLS13NoAES = []uint16{
	TLS_AES_128_GCM_SHA256,
	TLS_AES_256_GCM_SHA384,
	TLS_CHACHA20_POLY1305_SHA256,
}
```
- continue with the file `conn.go` to capture ciphertext traffic transcripts and mappings to keys. the file `conn.go` encrypts/decrypts traffic.

#### cipher suite details
- if you want to know how aesgcm encrypts/decrypts, look into cipher_suites.go and the function `aeadAESGCMTLS13` which shows that aes gcm xorNonceAEAD for tls1.3. thats important to know to identify which keys, iv, etc are used for encryption and decryption of speicifc blocks.

#### client side record storage
- inside the conn.go, I modified the halfConn struct to maintain a hashmap of record metadata. with the metadata, ther verifier can hash the e.g. SF record ciphertext, the with the hashmap, immediately locate the SHTS, nonce, and additional data which is required for the decryption. after decryption, the verifier has access to SF and can check the circuit public data if it matches the SHTS value received and this way verify the mapping to the intercepted ciphertext. the additional value inside the `halfConn` struct is `recordMap map[string]recordMeta`, and the following additional types/methods are used to store record metadata inside conn.go
```go
type recordMeta struct {
	additionalData []byte
	nonce          []byte
	typ            string
	// payload==SHTS, if typ==SF
	// payload==plaintext, if typ==RS
	payload []byte
}

func (hc *halfConn) setRecordMeta(ad, nonce, payload []byte, recordHash, typ string) {
	if hc.recordMap == nil {
		hc.recordMap = make(map[string]recordMeta)
	}
	hc.recordMap[recordHash] = recordMeta{
		additionalData: ad,
		nonce:          nonce,
		payload:        payload,
		typ:            typ,
	}
}
```
- once the tls session and the record reception completes on the client side, you can print record map values with `conn.ShowRecordMap()`. Other available functions are:
```go
func (c *Conn) ShowRecordMap() string {
	var buffer bytes.Buffer
	for k := range c.in.recordMap {
		buffer.WriteString(k)
		buffer.WriteString(";")
	}
	return buffer.String()
}
func (c *Conn) RawRecordMap() map[string]recordMeta {
	return c.in.recordMap
}
```
These functions are bound to the `Conn` struct

#### new files
- files which have been added the the tls_fork library are indicated with a `new_` and we maintain the modified crypto functions inside the tls_fork because e.g. aes_gcm decrypt makes use of a specific struct which is already defined in the forked tls lib. this way, functions will make use of the same implementation versions.
