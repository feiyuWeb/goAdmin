// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "运行地址：http://localhost/swagger/index.html",
        "contact": {},
        "license": {
            "name": "MIT //localhost:80"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/auth/info": {
            "get": {
                "description": "用户信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "单个用户信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "telephone",
                        "name": "telephone",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.UserDto"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"请求失败\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/login": {
            "post": {
                "description": "用户登陆",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "用户登陆",
                "parameters": [
                    {
                        "type": "string",
                        "description": "telephone",
                        "name": "telephone",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "password",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"请求失败\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/register": {
            "post": {
                "description": "用户注册",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "telephone",
                        "name": "telephone",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "password",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"请求失败\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/books": {
            "get": {
                "description": "书籍列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "书籍"
                ],
                "summary": "获取书籍列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "title",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "pageNum",
                        "name": "pageNum",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "pageSize",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Book"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"请求失败\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "创建书籍",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "书籍"
                ],
                "summary": "创建书籍列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "title",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "year",
                        "name": "year",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "author",
                        "name": "actor",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "desc",
                        "name": "desc",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Book"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"请求失败\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/books/{id}": {
            "put": {
                "description": "书籍列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "书籍"
                ],
                "summary": "更新书籍列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Book"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"id必传\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "书籍列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "书籍"
                ],
                "summary": "删除书籍列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Book"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"id必传\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/films": {
            "get": {
                "description": "电影列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "电影"
                ],
                "summary": "获取电影列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "pageNum",
                        "name": "pageNum",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "pageSize",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Film"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"请求失败\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "创建电影",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "电影"
                ],
                "summary": "创建电影列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "year",
                        "name": "year",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "address",
                        "name": "address",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "actor",
                        "name": "actor",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "desc",
                        "name": "desc",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Film"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"请求失败\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/films/{id}": {
            "put": {
                "description": "电影列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "电影"
                ],
                "summary": "更新电影列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Film"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"id必传\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "电影列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "电影"
                ],
                "summary": "删除电影列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Film"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"id必传\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/getCode": {
            "get": {
                "description": "获取验证码",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "验证码"
                ],
                "summary": "获取验证码",
                "responses": {
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"请求失败\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/upload": {
            "post": {
                "description": "文件上传",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文件上传"
                ],
                "summary": "文件上传",
                "parameters": [
                    {
                        "type": "string",
                        "description": "file",
                        "name": "file",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{ \"code\": 200, \"message\": \"上传成功\" }",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"请求失败\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/users": {
            "get": {
                "description": "用户列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "获取用户列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "telephone",
                        "name": "telephone",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "pageNum",
                        "name": "pageNum",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "pageSize",
                        "name": "pageSize",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"请求失败\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/users/{id}": {
            "put": {
                "description": "用户信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "更新用户信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"id必传\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "删除用户",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    },
                    "400": {
                        "description": "{ \"code\": 400, \"message\": \"请求失败\" }",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Book": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string",
                    "example": "作者"
                },
                "createAt": {
                    "type": "string",
                    "example": "创建时间"
                },
                "desc": {
                    "type": "string",
                    "example": "描述"
                },
                "id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string",
                    "example": "书籍名称"
                },
                "updateAt": {
                    "type": "string",
                    "example": "更新时间"
                },
                "year": {
                    "type": "string",
                    "example": "年份"
                }
            }
        },
        "model.Film": {
            "type": "object",
            "properties": {
                "actor": {
                    "type": "string",
                    "example": "演员"
                },
                "address": {
                    "type": "string",
                    "example": "出品地区"
                },
                "createAt": {
                    "type": "string",
                    "example": "创建时间"
                },
                "desc": {
                    "type": "string",
                    "example": "描述"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string",
                    "example": "电影名称"
                },
                "updateAt": {
                    "type": "string",
                    "example": "更新时间"
                },
                "year": {
                    "type": "string",
                    "example": "年份"
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "createAt": {
                    "type": "string",
                    "example": "创建时间"
                },
                "desc": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "ip": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "example": "用户名称"
                },
                "password": {
                    "description": "字段不暴露给用户，则使用 ` + "`" + `json:\"-\"` + "`" + ` 修饰",
                    "type": "string",
                    "example": "密码"
                },
                "telephone": {
                    "type": "string",
                    "example": "手机号码"
                },
                "updateAt": {
                    "type": "string",
                    "example": "更新时间"
                }
            }
        },
        "model.UserDto": {
            "type": "object",
            "properties": {
                "createAt": {
                    "type": "string",
                    "example": "创建时间"
                },
                "desc": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "telephone": {
                    "type": "string"
                },
                "updateAt": {
                    "type": "string",
                    "example": "更新时间"
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
	Version:     "2.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Golang Gin API",
	Description: "An example of gin",
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
