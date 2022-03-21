# feedfinder

A tiny Go package and executable to find RSS and Atom feed URLs
from the content of a Web page.

## Package usage example

```go
package main

import (
    "fmt"
    "github.com/SpecializedGeneralist/feedfinder"
    "log"
)

func main() {
    feeds, err := feedfinder.FindFeeds("https://go.dev", 2)
    
    if err != nil {
        log.Fatal(err)
    }
    
    for _, feed := range feeds {
        fmt.Println(feed)
    }
}
```

### Command line example 

```bash
feedfinder -url https://go.dev -depth 2
```

## License

Feedfinder is licensed under the
[BSD 2-Clause "Simplified" License](https://github.com/SpecializedGeneralist/feedfinder/blob/master/LICENSE).
