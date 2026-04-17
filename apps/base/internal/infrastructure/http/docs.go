package http

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

func (h *Hooks) handleOpenAPISpec(e *core.RequestEvent) error {
	e.Response.Header().Set("Content-Type", "application/yaml")
	_, err := e.Response.Write(h.openAPISpec)
	return err
}

func (h *Hooks) handleSwaggerUI(e *core.RequestEvent) error {
	specURL := fmt.Sprintf("http://%s/api/openapi.yaml", e.Request.Host)
	e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := fmt.Fprint(e.Response, swaggerHTML(specURL))
	return err
}

func swaggerHTML(specURL string) string {
	return `<!DOCTYPE html>
<html>
<head>
  <title>Preza API</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
  SwaggerUIBundle({
    url: "` + specURL + `",
    dom_id: "#swagger-ui",
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
    layout: "BaseLayout",
    deepLinking: true,
  })
</script>
</body>
</html>`
}
