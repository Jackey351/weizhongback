// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-04-22 13:03:05.927894085 +0800 CST m=+0.115094373

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "YANFEI API",
        "contact": {},
        "license": {},
        "version": "0.0.1"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ping": {
            "get": {
                "description": "测试服务器是否在线",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "miscellaneous"
                ],
                "summary": "PING-PONG",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controller.Message"
                        }
                    }
                }
            }
        },
        "/wx/info/project_types": {
            "get": {
                "description": "获取所有工程类别",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wx"
                ],
                "summary": "获取所有工程类别",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controller.Message"
                        }
                    }
                }
            }
        },
        "/wx/info/worker_types": {
            "get": {
                "description": "获取所有工种",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wx"
                ],
                "summary": "获取所有工种",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controller.Message"
                        }
                    }
                }
            }
        },
        "/wx/user/new_user": {
            "post": {
                "description": "小程序端新添用户",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wx",
                    "user"
                ],
                "summary": "小程序端新添用户",
                "parameters": [
                    {
                        "description": "create a new user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.WxUserWrapper"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controller.Message"
                        }
                    }
                }
            }
        },
        "/wx/work/publish": {
            "post": {
                "description": "发布工作",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wx"
                ],
                "summary": "发布工作",
                "parameters": [
                    {
                        "type": "string",
                        "description": "工种 0(点工),1(包工) 必填",
                        "name": "type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "点工招聘",
                        "name": "点工示例数据",
                        "in": "body",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.DianWorkWrapper"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controller.Message"
                        }
                    }
                }
            }
        },
        "/wx/work/search": {
            "get": {
                "description": "搜索工作，需要某个筛选参数就加上，否则可以不加",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wx"
                ],
                "summary": "搜索工作",
                "parameters": [
                    {
                        "type": "string",
                        "description": "二级位置信息 选填",
                        "name": "location",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "所需工种 选填",
                        "name": "need",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "工程类别 选填",
                        "name": "type",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "页码，从1开始 必填",
                        "name": "page",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "每页记录数 必填",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controller.Message"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.Message": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "model.DianWorkWrapper": {
            "type": "object",
            "properties": {
                "construction_company": {
                    "type": "string",
                    "example": "飞燕工程队"
                },
                "desc": {
                    "type": "string",
                    "example": "包吃包住"
                },
                "final_treatment": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "location": {
                    "type": "string",
                    "example": "湖北省襄阳市"
                },
                "location_info": {
                    "type": "object",
                    "$ref": "#/definitions/model.LocationInfoWrapper"
                },
                "max_wage": {
                    "type": "string",
                    "example": "200"
                },
                "min_wage": {
                    "type": "string",
                    "example": "100"
                },
                "need": {
                    "type": "string",
                    "example": "钢筋工"
                },
                "pricing_mode": {
                    "type": "string",
                    "example": "点工"
                },
                "project_name": {
                    "type": "string",
                    "example": "主楼建造"
                },
                "required_people": {
                    "type": "string",
                    "example": "11"
                },
                "settlement": {
                    "type": "string",
                    "example": "月薪"
                },
                "type": {
                    "type": "string",
                    "example": "消防"
                },
                "user_id": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "model.LocationInfoWrapper": {
            "type": "object",
            "properties": {
                "addr": {
                    "type": "string",
                    "example": "湖北省襄阳市樊城区武商广场店对面人民公园"
                },
                "latitude": {
                    "type": "string",
                    "example": "32.04278"
                },
                "longitude": {
                    "type": "string",
                    "example": "112.15519"
                },
                "title": {
                    "type": "string",
                    "example": "人民广场"
                }
            }
        },
        "model.WxUserWrapper": {
            "type": "object",
            "properties": {
                "hometown": {
                    "type": "string",
                    "example": "江苏"
                },
                "nick_name": {
                    "type": "string",
                    "example": "飞燕一号"
                },
                "phone": {
                    "type": "string",
                    "example": "133333"
                },
                "real_name": {
                    "type": "string",
                    "example": "张三"
                },
                "sex": {
                    "type": "string",
                    "example": "男"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
