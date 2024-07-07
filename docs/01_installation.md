## Installation

### Locally running the repo
- cd into the client folder and run `git clone --recurse-submodules git@github.com:jplaui/origo.git`
- cd into each folder of `client, proxy, server` and run `go mod tidy` to import all necessary go packages
- now follow the instructions provided in the README.md (running the protocol locally)

### Container
- has been done in the [opex-repositories](https://github.com/opex-research/tls-oracle-demo)

## Git Submodules

- add new submodules with `git submodule add -b <branch_name> <repository_url> <path/to/submodule>` or with `git submodule add <repository_url> <path/to/submodule>`
- then pull all submodules in clone with `git clone --recurse-submodules <main_repository_url>`
