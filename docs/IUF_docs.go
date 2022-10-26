/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplateIUF = `{
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
        "/iuf/v1/activities": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Activities"
                ],
                "summary": "List IUF activities",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/Activity"
                            }
                        }
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Activities"
                ],
                "summary": "Create an IUF activity",
                "parameters": [
                    {
                        "description": "IUF activity",
                        "name": "activity",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Activity.CreateActivityRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Activity"
                        }
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/iuf/v1/activities/{activity_uid}/sessions": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "List sessions of an IUF activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity uid",
                        "name": "activity_uid",
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
                                "$ref": "#/definitions/Session"
                            }
                        }
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Create a new session of an IUF activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity uid",
                        "name": "activity_uid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "IUF session",
                        "name": "session",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Session.CreateSessionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/Session"
                            }
                        }
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/iuf/v1/activities/{activity_uid}/sessions/{session_uid}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Get a session of an IUF activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity uid",
                        "name": "activity_uid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "session uid",
                        "name": "session_uid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Session"
                        }
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/iuf/v1/activities/{activity_uid}/sessions/{session_uid}/resume": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Resume a stopped session of an IUF activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity uid",
                        "name": "activity_uid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "session uid",
                        "name": "session_uid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Session"
                        }
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/iuf/v1/activities/{activity_uid}/sessions/{session_uid}/stop": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Sessions"
                ],
                "summary": "Stop a running session of an IUF activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity uid",
                        "name": "activity_uid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "session uid",
                        "name": "session_uid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Session"
                        }
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/iuf/v1/activities/{uid}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Activities"
                ],
                "summary": "Get an IUF activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity uid",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Activity"
                        }
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            },
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Activities"
                ],
                "summary": "Patch an IUF activity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "activity uid",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "partial IUF activity",
                        "name": "partial_activity",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Activity.PatchActivityRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Activity"
                        }
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/iuf/v1/stages": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Stages"
                ],
                "summary": "List stages of iuf",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Session"
                        }
                    },
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        }
    },
    "definitions": {
        "Activity": {
            "type": "object",
            "required": [
                "activity_states",
                "input_parameters",
                "operation_outputs",
                "products"
            ],
            "properties": {
                "activity_states": {
                    "description": "History of states",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Activity.State"
                    }
                },
                "input_parameters": {
                    "$ref": "#/definitions/Activity.CreateActivityRequest"
                },
                "operation_outputs": {
                    "description": "Operation outputs from argo",
                    "type": "object",
                    "additionalProperties": true
                },
                "products": {
                    "description": "List of products included in an activity",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Product"
                    }
                }
            }
        },
        "Activity.CreateActivityRequest": {
            "type": "object",
            "required": [
                "bootprep_config_managed",
                "bootprep_config_management",
                "media_dir",
                "name",
                "site_parameters"
            ],
            "properties": {
                "bootprep_config_managed": {
                    "description": "Each item is a path of the bootprep files",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "bootprep_config_management": {
                    "description": "Each item is a path of the bootprep files",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "limit_nodes": {
                    "description": "Each item is the xname of a node",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "media_dir": {
                    "description": "location of media",
                    "type": "string"
                },
                "name": {
                    "description": "Name of activity",
                    "type": "string"
                },
                "site_parameters": {
                    "description": "The inline contents of the site_parameters.yaml file.",
                    "type": "string"
                }
            }
        },
        "Activity.PatchActivityRequest": {
            "type": "object",
            "properties": {
                "comment": {
                    "type": "string"
                },
                "state": {
                    "type": "string"
                }
            }
        },
        "Activity.State": {
            "type": "object",
            "required": [
                "session_name",
                "start_time",
                "state"
            ],
            "properties": {
                "comment": {
                    "type": "string"
                },
                "session_name": {
                    "type": "string"
                },
                "start_time": {
                    "type": "string"
                },
                "state": {
                    "type": "string"
                }
            }
        },
        "Product": {
            "type": "object",
            "required": [
                "name",
                "original_location",
                "validated",
                "version"
            ],
            "properties": {
                "name": {
                    "description": "The name of the product",
                    "type": "string"
                },
                "original_location": {
                    "description": "The original location of the extracted tar in on the physical storage.",
                    "type": "string"
                },
                "validated": {
                    "description": "The flag indicates md5 of a product tarball file has been validated",
                    "type": "boolean"
                },
                "version": {
                    "description": "The version of the product.",
                    "type": "string"
                }
            }
        },
        "Session": {
            "type": "object",
            "required": [
                "products"
            ],
            "properties": {
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Product"
                    }
                },
                "stages": {
                    "description": "The stages that need to be executed.\nThis is either explicitly specified by the Admin, or it is computed from the workflow type.\nAn Stage is a group of Operations. Stages represent the overall workflow at a high-level, and executing a stage means executing a bunch of Operations in a predefined manner.  An Admin can specify the stages that must be executed for an install-upgrade workflow. And Product Developers can extend each stage with custom hook scripts that they would like to run before and after the stage's execution.  The high-level stages allow their configuration would revealing too many details to the consumers of IUF.\nif not specified, we apply all stages",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "Session.CreateSessionRequest": {
            "type": "object"
        }
    }
}`

// SwaggerInfoIUF holds exported Swagger Info so clients can modify it
var SwaggerInfoIUF = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/apis",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "IUF",
	SwaggerTemplate:  docTemplateIUF,
}

func init() {
	swag.Register(SwaggerInfoIUF.InstanceName(), SwaggerInfoIUF)
}
