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

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "MIT",
            "url": "https://github.com/Cray-HPE/cray-nls/blob/master/License"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/liveness": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Misc"
                ],
                "summary": "K8s Liveness endpoint",
                "responses": {
                    "204": {
                        "description": ""
                    }
                }
            }
        },
        "/v1/readiness": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Misc"
                ],
                "summary": "K8s Readiness endpoint",
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    }
                }
            }
        },
        "/v1/version": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Misc"
                ],
                "summary": "Get version of cray-nls service",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResponseOk"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    }
                }
            }
        },
        "/v1/workflows": {
            "get": {
                "security": [
                    {
                        "OAuth2Application": [
                            "admin",
                            "read"
                        ]
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Workflow"
                ],
                "summary": "Get status of a ncn workflow",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Label Selector",
                        "name": "labelSelector",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.GetWorkflowResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ResponseError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "ResponseOk": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.GetWorkflowResponse": {
            "type": "object",
            "properties": {
                "label": {
                    "type": "object"
                },
                "name": {
                    "type": "string"
                },
                "status": {
                    "type": "object"
                }
            }
        }
    },
    "tags": [
        {
            "description": "\u003e \u003e \u003e \u003e #### End to end rebuild of worker nodes\n",
            "name": "NCNs"
        },
        {
            "description": "\u003e \u003e \u003e #### Workflow management\n",
            "name": "Workflow"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/apis/nls",
	Schemes:          []string{},
	Title:            "NCN Lifecycle Management API",
	Description:      "## Security\n\n### Authentication\n\nSimilar to other exposed services, authentication is done by keycloak. Keycloak issued jwt token is verified and passed to API gateway.\n\n### Routes/AuthZ\n\nEach route of these APIs are protected by configuring OPA policy.\n\n- **Crawl Phase**\n\n  we will only `admin` and `user` roles. Users have `admin` role are allowed to invoke any APIs. Users with `user` role will only be able to call **GET** APIs.\n\n- **Walk Phase**\n\n  we can introduce more granular permissions/roles based on future requirements.\n\n- **Run Phase**\n\n  we can even go to resources level. For example, `User A` could have all permissions of `ceph nodes` but this user won't be able to rebuild/reboot any k8s nodes. `Monitoring User` can rerun/retry any failed rebuild/reboots but can't initiate such operation.\n\n### Microservices\n\nThe jwt token will be passed down to each microservices and individual microservice should enforce authZ in its own domain. Any credentials needed by each microservice should be obtained in a secure manner. SSH as root should be avoided if possible. However, there are certain operations requires root access via ssh. In those cases, we should use Vault to generate one time, short lived temporary SSH keys. Note that these goals will be achieved phase by phase.\n\n- **Crawl Phase**\n\n  In crawl phase, we execute steps almost identical to what we have today. Most steps need direct root access via SSH. SSH credentials are mounted onto each short lived _Job Pods_ as `hostPath`. JWT tokens needed for other microservice calls are obtained from `ncn-m001` over SSH:\n\n  ```\n  export TOKEN=$(curl -k -s -S -d grant_type=client_credentials \\\n   -d client_id=admin-client \\\n   -d client_secret=`kubectl get secrets admin-client-auth -o jsonpath='{.data.client-secret}' | base64 -d` \\\n   https://api-gw-service-nmn.local/keycloak/realms/shasta/protocol/openid-connect/token | jq -r '.access_token')\n  ```\n\n  > NOTE: This is exactly what our 1.0.x and 1.2.x does\n\n- **Walk Phase**\n\n  - SSH credentials need to be controlled by Vault and only one time credentials should be used\n  - JWT token should be passed from API gateway instead of getting it from `ncn-m001` as root user\n  - Any steps can be performed by make REST/gRPC request to a microservice should not use SSH any more\n\n- **Run Phase**\n\n  Each microservice should implement it's own granular/resources level authZ\n\n### Logging/Audit\n\n- **Request info**\n\n  API Gateway should log user information from validated JWT token so we know \"who did what at when\". Each microservice should also log the same information. Additionally, a unique request id should be passed/logged as well such that we can track a request in every microservice. Note that this is slightly different than what istio tracking is because of async operations. It won't carry istio injected `x-b3-traceid` in some cases.\n\n  Required fields:\n\n  - User Info: `name/id/email`, `roles`\n  - HTTP path: `REST API URI`\n  - HTTP method: `GET|POST|PUT|DELETE`\n  - Resources list: `ncn-w001,ncn-w002...`\n  - Operation Result: `failed|succeed|terminated`\n\n- **Operation logs**\n\n  Each steps of automation should be logged in order to troubleshoot/audit what exactly happened on ncn(s). This is done by _Argo Workflow_ engine.\n\n---\n\n[API Docs](https://cray-hpe.github.io/cray-nls/)\n",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
