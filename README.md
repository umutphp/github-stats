# github-stats
Command-line tool to get the total traffics statistics on your GitHub account

## How To Use

```bash
github-stats --token A_Valid_Personal_Access_Token
```

For details please use `--help` as argument;

```bash
$ go run main.go --help

NAME:
   github-stats - Get the total visit stats of your GitHub repositories
USAGE:
   github-stats [global options]
VERSION:
   0.0.2
AUTHOR:
   Umut Işık <umutphp@gmail.com>
OPTIONS:
   --day value, -d value           The number of days from today to show the stats (default: 0)
   --show-details value, -s value  Show detailed output or not. 0 to close. Default is 1 (default: 1)
   --token value, -t value         Personal access token got from GitHub to use the API
   --help, -h                      show help (default: false)
   --version, -v                   print the version (default: false)

```

## How To Contribute

All kind of contributions are ok for me :).
