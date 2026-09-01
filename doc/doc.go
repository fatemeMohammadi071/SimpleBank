package doc

import "embed"

// SwaggerFiles holds the embedded Swagger UI assets and the generated
// simple_bank.swagger.json, so they ship inside the server binary.
//
//go:embed swagger
var SwaggerFiles embed.FS
