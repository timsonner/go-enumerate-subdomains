# Enumerate subdomains

Three single-file Go programs look up names under a domain you pass on the command line. Use DNS resolution, plain HTTP, or the public [crt.sh](https://crt.sh/) certificate log. Each file is its own `package main`. They share no module and no common library.

`go-enum-subdomains-dns.go` reads a wordlist, builds `<label>.<domain>`, and calls `net.LookupIP`. A name that resolves is printed with its addresses. Lookup errors produce no line.

`go-enum-subdomains-http.go` reads the same kind of wordlist and calls `http.Get` on `http://<label>.<domain>`. It prints the URL when the status is 200, and it prints the URL plus status code for any other completed response. Connection errors produce no line.

`go-enum-subdomains-ssl.go` sends one GET to `https://crt.sh/?q=%25.<domain>&output=json` (the crt.sh wildcard query `%.domain`). It collects unique `common_name` and `name_value` strings that contain the domain and prints them.

## Prerequisites

- The `go` command. These files import only the standard library (`bufio`, `encoding/json`, `fmt`, `net`, `net/http`, `os`, `strings`, `sync`). There is no `go.mod` and nothing to `go get`.
- Network access for the lookups you run.
- For the DNS and HTTP programs, a wordlist with one label per line. This repo includes `subdomains.txt` (4988 lines, no blank lines).

Commands below were run with `go version go1.22.2 linux/amd64`.

## Install, build, run, test

There is no install step.

Build each program from the repository root:

```bash
go build -o enum-dns go-enum-subdomains-dns.go
go build -o enum-http go-enum-subdomains-http.go
go build -o enum-ssl go-enum-subdomains-ssl.go
```

Run from source:

```bash
go run go-enum-subdomains-dns.go example.com subdomains.txt
go run go-enum-subdomains-http.go example.com subdomains.txt
go run go-enum-subdomains-ssl.go example.com
```

Or run the binaries:

```bash
./enum-dns example.com subdomains.txt
./enum-http example.com subdomains.txt
./enum-ssl example.com
```

The DNS and HTTP programs start one goroutine per wordlist line and wait with `sync.WaitGroup`. There is no concurrency limit. `subdomains.txt` therefore opens one lookup per line at the same time.

This repository has no `_test.go` files and no test script.

## Usage

| Program | Arguments | Printed when the domain argument is missing |
| --- | --- | --- |
| `go-enum-subdomains-dns.go` | `<domain> <wordlist>` | `Usage: go run <script> <domain> <wordlist>` |
| `go-enum-subdomains-http.go` | `<domain> <wordlist>` | `Usage: go run <script> <domain>` |
| `go-enum-subdomains-ssl.go` | `<domain>` | `Usage: go run <script> <domain>` |

That usage line is printed only when `len(os.Args) < 2`, and the process then returns. The DNS and HTTP programs still read the wordlist from `os.Args[2]`. Passing a domain and omitting the wordlist panics with `index out of range [2] with length 2`.

Wordlist lines are used as scanned, with no trimming. A line `www` and a domain `example.com` become `www.example.com` (DNS) or `http://www.example.com` (HTTP).

DNS output for a hit:

```text
www.example.com resolves to [203.0.113.10]
```

HTTP output:

```text
Valid domain: http://www.example.com
Invalid domain or other status code: http://mail.example.com Status Code: 404
```

The request URL is the `http://` string built in the source.

Certificate output, after a JSON body is decoded:

```text
Discovered subdomains:
- www.example.com
- mail.example.com
```

Map iteration order is not stable. `name_value` entries are split on newlines and trimmed. A name is kept when it is non-empty and `strings.Contains` finds the domain in it.

If the crt.sh request fails, stderr gets `Failed to query crt.sh: ...`. If the body is not a JSON array of objects with `common_name` and `name_value`, stderr gets `Failed to decode JSON response: ...`. In both cases the process returns without a non-zero status.

Other errors printed by the DNS and HTTP programs:

```text
Error opening file: ...
Error reading file: ...
```

## Layout

```text
README.MD
go-enum-subdomains-dns.go
go-enum-subdomains-http.go
go-enum-subdomains-ssl.go
subdomains.txt
```

## License

This repository does not include a license file.
