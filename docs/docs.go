// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/booking": {
            "post": {
                "description": "Create a vaccine booking slot",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create Booking Slot",
                "parameters": [
                    {
                        "description": "booking info",
                        "name": "booking",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.Booking"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entities.Booking"
                        }
                    }
                }
            }
        },
        "/stations/{name}": {
            "get": {
                "description": "Get station slots by location",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "List Station",
                "parameters": [
                    {
                        "type": "string",
                        "description": "select the location. The only available option is PWTC",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Station"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entities.Booking": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string",
                    "example": "Kuala Lumpur"
                },
                "date": {
                    "type": "string",
                    "example": "20210516"
                },
                "firstName": {
                    "type": "string",
                    "example": "Fadhil"
                },
                "lastName": {
                    "type": "string",
                    "example": "Yaacob"
                },
                "location": {
                    "type": "string",
                    "example": "PWTC"
                },
                "mysejahteraId": {
                    "type": "string",
                    "example": "900127015527"
                },
                "phoneNumber": {
                    "type": "string",
                    "example": "0123456789"
                }
            }
        },
        "entities.Station": {
            "type": "object",
            "properties": {
                "availability": {
                    "type": "integer",
                    "example": 10
                },
                "date": {
                    "type": "string",
                    "example": "20210516"
                },
                "location": {
                    "type": "string",
                    "example": "PWTC"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "JKJAV API Server",
	Description: "High performant JKJAV API Server",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
