awsdefault
==========
[![Build Status](https://travis-ci.org/peterbueschel/awsdefault.svg?branch=master)](https://travis-ci.org/peterbueschel/awsdefault)
[![Go Report Card](https://goreportcard.com/badge/github.com/peterbueschel/awsdefault)](https://goreportcard.com/report/github.com/peterbueschel/awsdefault)
[![Coverage Status](https://coveralls.io/repos/github/peterbueschel/awsdefault/badge.svg?branch=master)](https://coveralls.io/github/peterbueschel/awsdefault?branch=master)

*Change [Amazon AWS](https://aws.amazon.com) profiles/accounts globally.*

This tool sets one of your AWS profiles, configured in your $HOME/.aws/credentials file (see  [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html)),
 to the default profile.
Means:
- adding parameters like `--profile=my-profile` to the aws cli
- or changing the environment variable with `export AWS_PROFILE=my-profile` in every new terminal window

are not required anymore.

Table of Contents
=================

* [Example Usage](#example-usage)
  * [Cli tool](#use-either-the-cli-tool)
  * [UI tool](#or-the-ui-tool)
    * [Linux](#linux)
    * [Windows](#windows)
* [Installation](#installation)
* [How it works](#how-it-works)
* [License](#license)

# Example Usage

## Use either the cli tool

### Change the default AWS profile to 'personal'

- command:

```bash
$ awsdefault to personal
```

### Disable/unset the AWS profile

command:

```bash
$ awsdefault rm
```

*the complete list of parameters can be found [here](cmd/awsdefault/readme.md#usage)*


## Or the UI tool

### Linux

![awsdefault-gkt3-example1](doc/awsdefault-gtk3-example1.gif?raw=true) 

*Note* [i3block](https://github.com/vivien/i3blocks) was used as statusbar in this example. You can find the config in the [doc folder](cmd/awsdefault-gtk3/doc/i3block-example.conf).

### Windows

TODO

# Installation

## Option 1 — Download binaries

precompiled binaries for Linux, (TODO: Windows and MacOS) are available at the [release] page.


```bash
curl 
```

## Option 2 — Compile it

- for the cli tool, see [here](cmd/awsdefault/readme.md#installation)
- for the gtk3-UI, see [here](cmd/awsdefault-gtk3/readme.md#installation)

## Configure your environment (only first time)

Set the environment variable `AWS_PROFILE` to `default` ([aws userguide](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html)).

### Linux

Add the following line to your .xinitrc, .zshrc or .bashrc file:

```bash
export AWS_PROFILE=default
```

### Windows

TODO

# How it works


The awsdefault (cli or UI) tool creates, changes or deletes the `[default]` profile section in your [AWS credentials file](https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html). This default section contains the same `aws_access_key_id` and `aws_secret_access_key` as stored in one of the other profile sections.
Together with the configured environment variable `AWS_PROFILE=default`, this approach enables or disables the credentials of the specific AWS profile. In other words, the default profile points to the required profile or will be deleted if no profile is needed.

- The environment variable

```bash
$ env | grep AWS_PROFILE
AWS_PROFILE=default
```

- The AWS credentials file — _no default was set; default was unset_

```bash
$ cat ~/.aws/credentials
[live]
aws_access_key_id     = A
aws_secret_access_key = B

[dev]
aws_access_key_id     = C
aws_secret_access_key = D

[personal]
aws_access_key_id     = E
aws_secret_access_key = F
```

- The AWS credentials file — _default was set to personal_ (comment was added automatically)

```bash
$ cat ~/.aws/credentials
; active_profile=personal
[default]
aws_access_key_id     = E
aws_secret_access_key = F
                         
[live]                   
aws_access_key_id     = A
aws_secret_access_key = B
                         
[dev]                    
aws_access_key_id     = C
aws_secret_access_key = D
                         
[personal]               
aws_access_key_id     = E
aws_secret_access_key = F
```

# License

See [LICENSE](LICENSE).
