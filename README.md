# tfaction-go

CLI for [tfaction](https://github.com/suzuki-shunsuke/tfaction).

This CLI is used in GitHub Actions Workflows built with tfaction.
This CLI was introduced for Drift Detection.

https://suzuki-shunsuke.github.io/tfaction/docs/feature/drift-detection

## Usage

```console
$ tfaction help
NAME:
   tfaction - GitHub Actions Workflow for Terraform. https://github/com/suzuki-shunsuke/tfaction-go

USAGE:
   tfaction [global options] command [command options] [arguments...]

VERSION:
   0.1.1 (1263eaf5834fba37f89593415e76f21e7e276846)

COMMANDS:
   version                    Show version
   create-drift-issues        Create GitHub Issues for Terraform drift detection
   pick-out-drift-issues      Pick out GitHub Issues for Terraform drift detection
   get-or-create-drift-issue  Get or Create a GitHub Issue for Terraform drift detection
   help, h                    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level value  log level [$TFACTION_LOG_LEVEL]
   --help, -h         show help
   --version, -v      print the version
```

## LICENSE

[MIT](LICENSE)
