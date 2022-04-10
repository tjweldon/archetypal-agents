# Archetypal Agents

## Setup

### Prerequisites

 - Install [Golang v1.18](https://go.dev/dl/).
 - A web browser.
 - That's pretty much it for now.

### Build

 - Refer to Golang's documentation for instructions on how to install Golang on your OS.
 - Use the following command from the project root to download [p5.js](https://p5js.org/download/) into the correct directory:
```shell
mkdir resources && cd resources
wget https://github.com/processing/p5.js/releases/download/v1.4.1/p5.js
```
 - Verify that this is in the correct place using git status, the directory is ignored so the working tree should still be clean.
 - build the project using the following command from the project root:
```shell
go build
```
 - This will download and install all the necessary dependencies and compile a native binary for your OS. You don't need to do this every time you make changes but it will highlight any missing pre-requisites.
 - You should see a new executable file in the directory called `archetypal-agents`. Execute the binary from the project root with the following command:
```shell
./archetypal-agents
```

## Development

This is just what I can think of right now, but we can add to it as we go.

### Conventions/Style
Go has a notion of public and private attributes and functions. Published parts of a package or interface are denoted with `ClassCase`. This is not 
just a convention it's part of the syntax! At the minimum, we should try to document (in a comment above the definition) all published objects.

Testing is pretty easy in go so if you find it useful, then feel free to add them to the project (within reason), however we're not going to enforce any
arbitrary minimums on coverage. The idea is to keep the cost of change down in the early going and as things crystallise we may want to introduce tests to
ensure interfaces are maintained under refactoring and so on.

### Workflow

Go has a convenient way to build and execute a project without needing to create the executable file in a way that is reminiscent of interpreted languages
```shell
go run main.go
```
This will install any missing package dependencies, compile the source code to a binary in a temporary file and then execute it in your current working directory, connecting up in command line arguments, stdin, stdout and stderr exactly as you would expect if you were just execute the binary.
This turns out to be useful for development where it is easy to forget to build the new binary after making changes to source code.

Try not to check compiled binaries into the repo, they're pretty much just dead weight, and likely only function on your specific system.

If you want to make changes create a branch from main and when you're ready submit a PR for review. Async code review is not discouraged but should
be considered insufficient for approval, at least initially (applies to me too). We can review this once some PR's have been merged. The intention here
is to establish a shared understanding of the codebase and objectives as early as possible.

This repo is intended to track digital assets associated to the project, not just code but design docs and image assets etc.

