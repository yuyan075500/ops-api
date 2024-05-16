// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/api/v1/group": {
            "put": {
                "description": "组相关接口",
                "tags": [
                    "组管理"
                ],
                "summary": "更新组信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "组信息",
                        "name": "group",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.GroupUpdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"更新成功\", \"data\": nil}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "组相关接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "组管理"
                ],
                "summary": "创建组",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "组信息",
                        "name": "group",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.GroupCreate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"创建成功\", \"data\": nil}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/group/users": {
            "put": {
                "description": "组相关接口",
                "tags": [
                    "组管理"
                ],
                "summary": "更新组用户",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "用户信息",
                        "name": "users",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.GroupUpdateUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"更新成功\", \"data\": nil}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/group/{id}": {
            "delete": {
                "description": "组相关接口",
                "tags": [
                    "组管理"
                ],
                "summary": "删除组",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "组ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"删除成功\", \"data\": nil}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/groups": {
            "get": {
                "description": "组相关接口",
                "tags": [
                    "组管理"
                ],
                "summary": "获取组列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "分页",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "分页大小",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "组名称",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"获取列表成功\", \"data\": []}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/menus": {
            "get": {
                "description": "菜单关接口",
                "tags": [
                    "菜单管理"
                ],
                "summary": "获取菜单列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "分页",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "分页大小",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"data\": []}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user": {
            "put": {
                "description": "用户相关接口",
                "tags": [
                    "用户管理"
                ],
                "summary": "更新用户信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "用户信息",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dao.UserUpdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"更新成功\", \"data\": nil}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "用户相关接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户管理"
                ],
                "summary": "创建用户",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "用户信息",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.UserCreate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"创建成功\", \"data\": nil}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/avatarUpload": {
            "post": {
                "description": "用户相关接口",
                "tags": [
                    "用户管理"
                ],
                "summary": "头像上传",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "头像",
                        "name": "avatar",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"头像更新成功\", \"data\": nil}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/info": {
            "get": {
                "description": "认证相关接口",
                "tags": [
                    "用户认证"
                ],
                "summary": "获取用户信息",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"获取用户信息成功\", \"data\": {}}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/list": {
            "get": {
                "description": "用户相关接口",
                "tags": [
                    "用户管理"
                ],
                "summary": "获取所有的用户列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"获取列表成功\", \"data\": []}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/menu": {
            "get": {
                "description": "菜单关接口",
                "tags": [
                    "菜单管理"
                ],
                "summary": "获取用户菜单",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"data\": []}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/reset_mfa/{id}": {
            "put": {
                "description": "用户相关接口",
                "tags": [
                    "用户管理"
                ],
                "summary": "MFA重置",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "用户ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"重置成功\", \"data\": nil}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/reset_password": {
            "put": {
                "description": "用户相关接口",
                "tags": [
                    "用户管理"
                ],
                "summary": "密码更新",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "用户信息",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dao.UserPasswordUpdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"更新成功\", \"data\": nil}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/user/{id}": {
            "delete": {
                "description": "用户相关接口",
                "tags": [
                    "用户管理"
                ],
                "summary": "删除用户",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "用户ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"删除成功\", \"data\": nil}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/users": {
            "get": {
                "description": "用户相关接口",
                "tags": [
                    "用户管理"
                ],
                "summary": "获取查询的用户列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "分页",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "分页大小",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "用户姓名",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"获取列表成功\", \"data\": []}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "认证相关接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户认证"
                ],
                "summary": "登录",
                "parameters": [
                    {
                        "description": "用户名密码",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.UserLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"认证成功\", \"token\": \"用户令牌\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/logout": {
            "post": {
                "description": "认证相关接口",
                "tags": [
                    "用户认证"
                ],
                "summary": "注销",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 用户令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 0, \"msg\": \"注销成功\", \"data\": nil}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dao.UserPasswordUpdate": {
            "type": "object",
            "required": [
                "id",
                "password",
                "re_password"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                },
                "password": {
                    "type": "string"
                },
                "re_password": {
                    "type": "string"
                }
            }
        },
        "dao.UserUpdate": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "is_active": {
                    "type": "boolean"
                },
                "phone_number": {
                    "type": "string"
                }
            }
        },
        "service.GroupCreate": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "is_role_group": {
                    "type": "boolean",
                    "default": false
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "service.GroupUpdate": {
            "type": "object",
            "required": [
                "id",
                "name"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "service.GroupUpdateUser": {
            "type": "object",
            "required": [
                "id",
                "users"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "service.UserCreate": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password",
                "phone_number",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone_number": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "service.UserLogin": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
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
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
