# binoc
<img src="binoc.png" width="300" height="300">

Binoc is a GitHub Actions Workflow & Docker Container that can update package build instructions in GitHub Repositories.

## Introduction
Binoc is an automatic software maintainer in a box. It builds off of other projects like [Lookout](https://github.com/alecbcs/lookout) and [Cuppa](https://datadrake/cuppa) to check for upstream software releases on GitHub, GitLab, Sourceforge, FTP, etc.... All Binoc requires is access to a directory of package build instructions and it's able to parse the instructions, check for package updates, and patch the build instructions with newer version information if detected. If Binoc is pointed at a Git Repository it will sperate its update process into seperate branches for each package and submit the updates as pull requests on GitHub. Checkout an example of Binoc in action over at the [container-blueprints](https://github.com/autamus/container-blueprints/pulls) repository.

#### Supported Package Formats

- Spack (`spack`) parses [spack.yaml](https://spack.readthedocs.io/en/latest/configuration.html#yaml-format) files.
- Singularity Registry HPC (`shpc`) parsers the [container.yaml](https://singularity-hpc.readthedocs.io/en/latest/getting_started/developer-guide.html#registry-yaml-files) from shpc.
- Dockerfile (`dockerfile`) parses the first `FROM` discovered (does not handle multistage build `FROM`)

If you'd like to request a special parser, please [open an issue](https://github.com/autamus/binoc/issues).

## Usage

Binoc can be run as either a GitHub Action or a Docker Container. You can have binoc
manage your pull requests (the default) or set `skip_pr` to "true" to manage them on
your own.

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
        with:
          fetch-depth: '0'
        
      - name: Run Binoc Scan
        uses: autamus/binoc@v0.3.4
        with:
          git_token: ${{ secrets.BINOC_GIT_TOKEN }}
          git_username: ${{ secrets.BINOC_GIT_USERNAME }}
          git_email: ${{ secrets.BINOC_GIT_EMAIL }}
          git_name: ${{ secrets.BINOC_GIT_NAME }}
          # Location of the git repository. For all actions this should be '/'.
          repo_path: '/'
          parsers_loaded: 'spack'
          general_action: 'true'
          # pr_skip: 'true'   # Disable Binoc's ability to create pull requests.
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


You can then test binoc on a repository by way of exporting environment variables for the
command. For example, here is for shpc:

```bash
BINOC_REPO_PATH=/path/to/test/ BINOC_PARSERS_LOADED=shpc BINOC_GIT_TOKEN=ghp_xxxx go run binoc.go
```

And here is for checking the Dockerfile in another directory

```bash
BINOC_REPO_PATH=/path/to/test BINOC_PARSERS_LOADED=dockerfile BINOC_GIT_TOKEN=ghp_xxxx go run binoc.go
```

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
