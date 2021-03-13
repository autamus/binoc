# binoc
<img src="binoc.png" width="300" height="300">

Binoc is a GitHub Actions Workflow & Docker Container that can update package build instructions in GitHub Repositories.

## Introduction
Binoc is an automatic software maintainer in a box. It builds off of other projects like [Lookout](https://github.com/alecbcs/lookout) and [Cuppa](https://datadrake/cuppa) to check for upstream software releases on GitHub, GitLab, Sourceforge, FTP, etc.... All Binoc requires is access to a directory of package build instructions and it's able to parse the instructions, check for package updates, and patch the build instructions with newer version information if detected. If Binoc is pointed at a Git Repository it will sperate its update process into seperate branches for each package and submit the updates as pull requests on GitHub. Checkout an example of Binoc in action over at the [container-blueprints](https://github.com/autamus/container-blueprints/pulls) repository.

#### Supported Package Formats
- Spack

(Others Comming Soon)

## Usage
Binoc can be run as either a GitHub Action or a Docker Container

#### Automatic Daily Scan of the Repository containing Spack Packages at 7am PST.
###### .github/workflows/binoc.yaml
```yaml
name: "Automatic Binoc Scan"

on:
  schedule:
    - cron: '0 14 * * *'

jobs:
  auto-scan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v2
        
      - name: Run Binoc Scan
        uses: autamus/binoc@v0.0.8
        with:
          git_token: ${{ secrets.BINOC_GIT_TOKEN }}
          git_username: ${{ secrets.BINOC_GIT_USERNAME }}
          git_email: ${{ secrets.BINOC_GIT_EMAIL }}
          git_name: ${{ secrets.BINOC_GIT_NAME }}
          # Location of the git repository. For all actions this should be '/'.
          repo_path: '/'
          parsers_loaded: 'spack'
          general_action: 'true'
```

## Contributing
If you are interested in contributing to Binoc, we love pull requests! If you've got a favorite packaging format, we're looking for help writing more package parsers! Checkout [parsers/](https://github.com/autamus/binoc/tree/main/parsers) for examples on how to get started and implement the parser interface.

## Building Binoc From Source
#### Development Dependencies

- `GCC`

- `Golang`

#### Build

1. Clone this repository

2. `go run binoc.go` <-- run binoc for testing purposes.

3. `go build binoc.go` <-- will build binoc into an executable on your machine.


## License

Copyright 2021 Alec Scott & Autamus <hi@alecbcs.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
