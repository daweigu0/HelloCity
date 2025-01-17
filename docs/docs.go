// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
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
        "/files/upload_token": {
            "post": {
                "description": "文件上传token获取接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "文件相关接口"
                ],
                "summary": "文件上传token获取接口",
                "parameters": [
                    {
                        "description": "上传文件token获取参数",
                        "name": "UploadToken",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web.reqUploadToken"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":xxx,\"data\":{},\"msg\":\"xxx\"}",
                        "schema": {
                            "$ref": "#/definitions/ginx.Result"
                        }
                    }
                }
            }
        },
        "/users/login": {
            "post": {
                "description": "登录成功返回的token放在响应的header的x-jwt-token里面，登录之后的后续访问需要带上token，放在请求的header里面的Authorization。",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户相关接口"
                ],
                "summary": "用户登录接口",
                "parameters": [
                    {
                        "description": "微信登录的临时登录凭证",
                        "name": "login",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web.loginReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":xxx,\"data\":{},\"msg\":\"xxx\"}",
                        "schema": {
                            "$ref": "#/definitions/ginx.Result"
                        }
                    }
                }
            }
        },
        "/users/signup": {
            "post": {
                "description": "用户注册接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户相关接口"
                ],
                "summary": "用户注册接口",
                "parameters": [
                    {
                        "description": "注册参数",
                        "name": "signup",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/web.SignUpReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":xxx,\"data\":{},\"msg\":\"xxx\"}",
                        "schema": {
                            "$ref": "#/definitions/ginx.Result"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ginx.Result": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "msg": {
                    "type": "string"
                }
            }
        },
        "web.SignUpReq": {
            "type": "object",
            "properties": {
                "mobile": {
                    "type": "string"
                },
                "nick_name": {
                    "type": "string"
                },
                "signup_token": {
                    "type": "string"
                }
            }
        },
        "web.loginReq": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                }
            }
        },
        "web.reqUploadToken": {
            "type": "object",
            "properties": {
                "file_ext": {
                    "type": "string"
                },
                "file_type": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{"http", "https"},
	Title:            "你好同城后端API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
