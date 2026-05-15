package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/health": {
            "get": {
                "description": "Check whether the HTTP server is running.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "system"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Body"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "response.Body": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 0
                },
                "data": {},
                "message": {
                    "type": "string",
                    "example": "ok"
                }
            }
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Feed System API",
	Description:      "Feed System API documentation.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
