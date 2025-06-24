# Origo Artifact (PETS ARTIFACT-EVALUATION.md)

Paper title: **ORIGO: Proving Provenance of Sensitive Data with Constant Communication**

Artifacts HotCRP Id: **#27**

Requested Badge: **Available**, **Functional**, and **Reproduced**

## Description
The artifacts linked in this repository contribute to the creation of Figures and Tables in our research paper [Origo](https://petsymposium.org/popets/2025/popets-2025-0069.pdf).
The work Origo is a protocol to verify the provenance of TLS data without using secure two-party computation techniques.
Our artifact is an entire e2e protocol implementation that can be executed on a laptop.

### Security/Privacy Issues and Ethical Concerns (All badges)
Origo is a research project and our implementation has *not* undergone any security audit. Do not use our Origo in a production system and do not use Origo to process sensitive/confidential data.
Origo relies on few dependencies which we do not monitor for changes (e.g. vulnerability patches or similar).
We neither audited external dependencies with respect to security problems and we solely rely on the functionalities they provide.
Be aware that using the referenced tools enable online interactions with external websites and services.
Running Origo allows external services to collect networking data of your devices.
Make sure to obtain data processing agreements with any external site that is queried with the Origo protocol.
If you run Origo in the provided `local` mode, then all data exchanges happen locally and nothing is shared over the Internet.

## Basic Requirements (Only for Functional and Reproduced badges)
In order to run the provided artifacts, no special hardware is required.
With regard to software requirements, we explain tooling and additional requirements below.
The expected times to run our examples align with the benchmarks we provide in our paper.

### Hardware Requirements
This work does not require any specific hardware to reproduce the results.
We used a MacBook Pro 16-inch (2021) configured with an Apple M1 Pro chip and 32 GB of RAM.
All tests and executions have been executed on the defined laptop.

### Software Requirements
Our work entirely relies on publicly accessible software. Our code base can be executed on different machines, different operating systems, and with different dependency versions. However, the provided list shares the software stack we used to test the code base:
- MacOS Sonoma Version 14.6.1, [installation](https://support.apple.com/en-us/105111)
- Golang 1.21.1 darwin/arm64, [installation](https://go.dev/doc/install)
- Git version 2.39.3, [installation](https://git-scm.com/download/mac)
- GNU Make version 3.81, [installation](https://formulae.brew.sh/formula/make)

**Note:** If your software modules deviate from the versions we mention, it remains very likely that the code base continues to work. If not, please take a look into the changelog of the software modules and adapt outdated functionalities that are used within this work.

### Estimated Time and Storage Consumption
Running the Origo protocol (proving smaller chunks (< 64 bytes) of private data) will roughly take about 2 minutes.
The disk requirements remain below 5GB and the RAM allocation depends on the availability of resources. Origo can be executed on a laptop with 8 GB of RAM.

## Environment 
In order to run the code, Git and Golang must be installed.

In the following, we describe how to access our artifacts and all related and necessary data and software components.
Afterwards, we describe how to set up everything and how to verify that everything is set up correctly.

### Accessibility (All badges)
The artifacts can be accessed from a laptop with Internet connectivity.
You are required to download a Github repository to execute the code.

### Set up the environment (Only for Functional and Reproduced badges)
We expect that the following commands are executed from a fixed location such as `/home/user/origo/`.
To download Origo, run the following command `git clone https://github.com/jplaui/origo.git`. Next, `cd` into the folder `origo` and run the local installation git submodule commands provided [here](https://github.com/jplaui/origo/blob/main/docs/01_installation.md#locally-running-the-repo).
Additionally, the repositories [origo-project: tls oracle demo](https://github.com/origo-project/tls-oracle-demo) must be cloned and to for benchmarking individual ZKP circuits, the repository [circuits_janus - gnark_zkp](https://github.com/jplaui/circuits_janus/tree/main) must be cloned.

### Testing the Environment (Only for Functional and Reproduced badges)
In order to make sure that the software has been set up correctly, please go through the following list of bullets and make sure that the expected outputs can be recreated.
- (Git) run `git version` and expect version 2.39.3
- (Golang) run `go version` and expect version go1.21.1
- (Origo local mode) run the command `go mod tidy` in each of the folders `origo/client`, `origo/proxy`, `origo/server`.
- (Origo ZKP circuits) run `git clone https://github.com/jplaui/circuits_janus.git`, then `cd` into the folder `circuits_janus/gnark_zkp` and run `go mod tidy`.
- (Origo real server mode) run `https://github.com/origo-project/tls-oracle-demo.git`.

## Artifact Evaluation (Only for Functional and Reproduced badges)
This section includes all the steps required to evaluate our artifact's functionality and validate our paper's key results and claims.
We highlight our paper's main results and claims in the first subsection. And describe the experiments that support our claims in the subsection after that.

### Main Results and Claims
In total, we provide three main and one side claims in the context of our work. 
We mention additional other claims that can be drawn from the benchmarks we provide with respect to related works.

#### Main Result 1: Origo is an end-to-end practical protocol which uses a highly optimized ZKP circuit (supporting TLS 1.3).
Our results of Table 3 indicate the efficiency of Origo ZKP circuits.
The first and second experiments support this claim.

#### Main Result 2: Origo allows an Internet-wide deployment with efficient online times.
Our results of Figure 14, Figure 12, and Table 2 indicate that Origo is a bandwidth-optimized protocol to prove TLS data in seconds.
The third experiment supports this claim.

### Experiments 
Our list provided next explains each experiment which involves the following parts:
 - How to execute the experiment in detailed steps.
 - What is the expected result.
 - How long does the experiments take and how much space do they consume on disk. (approximately)
 - Which claim and results does our experiments support, and how.

#### Experiment 1: Origo ZKP (zkSNARK) Benchmarks
Make sure to jump into the folder `circuits_janus/gnark_zkp/circuits` and run `go mod tidy`. 
We provide the contents of our go.mod file next to show which version of gnark we use:
```bash go
module circuits

go 1.19

require (
	github.com/consensys/gnark v0.7.2-0.20230518132517-274c883477ec
	github.com/consensys/gnark-crypto v0.11.1-0.20230508024855-0cd4994b7f0b
	github.com/montanaflynn/stats v0.7.1
	github.com/rs/zerolog v1.29.1
)

require (
	github.com/bits-and-blooms/bitset v1.5.0 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fxamacker/cbor/v2 v2.4.0 // indirect
	github.com/google/pprof v0.0.0-20230309165930-d61513b1440d // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/stretchr/testify v1.8.2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)
```
You can see possible circuits with the command `go run main.go -debug -algorithm`. Use the command `go run main.go -debug -gcm -byte-size 256 -iterations 1` to evaluate the AES_GCM zkSNARK circuit and continue to change the size of bytes for the proof using the -byte-size flag. Further, you can increase the number of executions with the -iterations flag. We used the following list of commands to create the aes_gcm benchmarks of Figure 13:
- `go run main.go -debug -gcm -byte-size 16 -iterations 3`
- `go run main.go -debug -gcm -byte-size 32 -iterations 3`
- `go run main.go -debug -gcm -byte-size 64 -iterations 3`
- `go run main.go -debug -gcm -byte-size 128 -iterations 3`
- `go run main.go -debug -gcm -byte-size 256 -iterations 3`
- `go run main.go -debug -gcm -byte-size 512 -iterations 3`
- `go run main.go -debug -gcm -byte-size 1024 -iterations 3`

The command `go run main.go -debug -tls13-oracle -iterations 1` executes the entire Origo ZKP circuit (including the key derivation circuit). With the commands `go run main.go -debug -record -iterations 1` and `go run main.go -debug -authtag -iterations 1`, we create the benchmarks of Table 3. 
This experiment takes longer if byte sizes increase. The number of iterations additionally impacts the runtime of the algorithm. As much system resources as possible are allocated for this job.
The results supports our claims made in our main result 1.

#### Experiment 2: Origo benchmarks
In order to get local Origo end-to-end benchmarks, continue the execution of the command sequence provided [here](https://github.com/jplaui/origo/tree/main?tab=readme-ov-file#running-the-protocol-locally). 
This experiment does not take long and has no drastic consumption of system resources.
The results supports our claims made in our first main result.

#### Experiment 3: Origo WAN (realserver) benchmarks
To recreate the benchmarks of Figure 14, [this repository - tls oracle demo](https://github.com/origo-project/tls-oracle-demo) must be executed and deployed to different machines. We used AWS and selected servers in different locations to run the client. 
This experiment does not take long and has no drastic consumption of system resources.
The results supports the claims made in our second main result.

## Limitations (Only for Functional and Reproduced badges)
The current limitation we have in the public [origo](https://github.com/jplaui/origo) repository is that cli commands must be called sequentially. The purpose of running commands serves for simplifying access to different modules and allows splitting up parts to put them on different machines. Further, we do not fully automate the communication of shared parameters and store them in local folders (e.g. `origo/proxy/local_storage/*`). However, our server-focused implementation of [origo](https://github.com/origo-project/tls-oracle-demo) automates communication and removes sequential calls such that the protocol runs with a single invocation of the client.

## Notes on Reusability (Only for Functional and Reproduced badges)
Our Origo PoC implements a proxy deployment of TLS oracles including a SNI-based routing of session establishments.
Origo has been used and applied by [pluto.xyz](https://pluto.xyz/blog/web-proof-techniques-origo-mode).
