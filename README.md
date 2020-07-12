# feedfinder

A tiny Go package and executable to find RSS and Atom feed URLs
from the content of a Web page.

## Package usage example
 
```go
package main

import (
    "fmt"
    "github.com/nlpodyssey/feedfinder"
    "log"
)

func main() {
    feeds, err := feedfinder.FindFeeds("https://blog.golang.org")
    
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
feedfinder https://blog.golang.org
```

## License

Feedfinder is licensed under the
[BSD 2-Clause "Simplified" License](https://github.com/nlpodyssey/feedfinder/blob/master/LICENSE).
