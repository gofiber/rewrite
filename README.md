### Install
```
go get -u github.com/gofiber/fiber
go get -u github.com/gofiber/rewrite
```
### Example
```go
package main

import (
  "github.com/gofiber/fiber"
  "github.com/gofiber/rewrite"
)

func main() {
  app := fiber.New()
  
  app.Use(rewrite.New(rewrite.Config{
    Rules: map[string]string{
      "/old":   "/new",
      "/old/*": "/new/$1",
    },
  }))
  
  app.Get("/new", func(c *fiber.Ctx) {
    c.Send("Hello, World!")
  })
  app.Get("/new/*", func(c *fiber.Ctx) {
    c.Send("Wildcard: ", c.Params("*"))
  })
  
  app.Listen(3000)
}

```
### Test
```curl
curl http://localhost:3000/old
curl http://localhost:3000/old/hello
```
