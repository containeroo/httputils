# httputils

The `httputils` package provides small helpers for parsing HTTP headers and status-code expressions used by command-line flags, environment variables, and configuration files.

## Features

- Parse HTTP status codes and ranges such as `200,300-302,404`.
- Parse `Key=Value` header entries into `http.Header`.
- Canonicalize HTTP header names.
- Preserve duplicate header values when explicitly allowed.
- Preserve commas and additional `=` characters inside header values.

## Installation

```sh
go get github.com/containeroo/httputils@latest
```

## Usage

### ParseStatusCodes

```go
package main

import (
    "fmt"
    "log"

    "github.com/containeroo/httputils"
)

func main() {
    statusCodes, err := httputils.ParseStatusCodes("200,300-302,404")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(statusCodes)
}
```

Output:

```text
[200 300 301 302 404]
```

### ParseHeaders

Pass each header as a separate entry. This avoids treating commas inside valid header values as separators.

```go
package main

import (
    "fmt"
    "log"

    "github.com/containeroo/httputils"
)

func main() {
    headers, err := httputils.ParseHeaders([]string{
        "Content-Type=application/json",
        "Accept=text/html, application/json",
        "X-Trace=one",
        "X-Trace=two",
    }, true)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(headers.Get("Content-Type"))
    fmt.Println(headers.Get("Accept"))
    fmt.Println(headers.Values("X-Trace"))
}
```

Output:

```text
application/json
text/html, application/json
[one two]
```

`ParseHeaders` canonicalizes header names using `net/http`. When `allowDuplicates` is `false`, repeated names are rejected case-insensitively. When it is `true`, all values are retained in the returned `http.Header`.

## Error handling

`ParseStatusCodes` and `ParseHeaders` return errors for malformed input, including invalid status expressions, missing header separators, empty header names, and disallowed duplicate headers.
