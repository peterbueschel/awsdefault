awsdefault — cli tool
=====================

# Usage

## List all available AWS profiles 
- command:

```bash
$ awsdefault ls
```
- example output:

```bash
dev
live
personal
```

## Show the current used AWS profile

- command:

```bash
$ awsdefault is
```
- example output:

```bash
live
```

## Change the default AWS profile to 'personal'

- command:

```bash
$ awsdefault to personal
```

## Disable/unset the AWS profile

command:

```bash
$ awsdefault rm
```

## Show the AWS_ACCESS_KEY_ID of currently used profile

command:

```bash
$ awsdefault id
```

- example output:

```bash
AAAAAAABBBBIIIIII
```

- tip: use it in parameters or bash scripts

```bash
AWS_ACCESS_KEY_ID=$(awsdefault id)
```

## Show the AWS_SECRET_ACCESS_KEY of currently used profile

command:

```bash
$ awsdefault key
```

- example output:

```bash
aaaaeenntrnggg/trntruaelvii
```

## Print the export command for AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY of currently used profile

command:

```bash
$ awsdefault export
```

- example output:

```bash
export AWS_ACCESS_KEY_ID=AAAAAAABBBBIIIIII
export AWS_SECRET_ACCESS_KEY=aaaaeenntrnggg/trntruaelvii
```

- tip: add an alias to your .bashrc or .zshrc like

```bash
alias awsexport='eval $(awsdefault export)'
```

# Installation

## Option 1 — Download binaries

precompiled binaries for Linux, Windows and MacOS are available at the [release] page.

```bash
curl 
```

## Option 2 — Compile it


### Install the dependencies

- *[Go](https://golang.org/doc/install)* is required
- clone this repository: 

```bash
$ go get github.com/peterbueschel/awsdefault
```

- [go-ini](https://github.com/go-ini/ini); used for the handling of the [AWS credentials file](https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html),

```bash
$ go get github.com/go-ini/ini
```

- [urfave/cli](https://github.com/urfave/cli); used for the cli tool itself

```bash
$ go get github.com/urfave/cli
```

### Install the Go binary

#### Linux 

```bash
$ cd $GOPATH/src/github.com/peterbueschel/awsdefault/cmd/awsdefault/ && go install
```

*if everything went well, the binary can now be found in the directory* _$GOPATH/bin_ 

## Configure your environment (only first time)

Set the environment variable `AWS_PROFILE` to `default` ([aws userguide](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html)).

### Linux

Add the following line to your .xinitrc, .zshrc or .bashrc file:

```bash
export AWS_PROFILE=default
```

### Windows

TODO


