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
        "/forgot-password": {
            "post": {
                "description": "Request a password reset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Forgot password",
                "parameters": [
                    {
                        "description": "Email",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.ForgotPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Check the health of the service",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health Check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Authenticate a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.LoginResponse"
                        }
                    }
                }
            }
        },
        "/logout": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Logout a user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/me": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Get the profile of the logged-in user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Get user profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.UserResponse"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Register a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dto.RegisterResponse"
                        }
                    }
                }
            }
        },
        "/remove-users": {
            "delete": {
                "description": "Remove all users from the database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Remove all users",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/reset-password": {
            "post": {
                "description": "Reset the user's password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Reset password",
                "parameters": [
                    {
                        "description": "Reset Password",
                        "name": "resetPasswordRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.ResetPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Get a list of all users",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.User"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.ForgotPasswordRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "dto.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "dto.LoginResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/dto.UserResponse"
                }
            }
        },
        "dto.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 8
                }
            }
        },
        "dto.RegisterResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/dto.UserResponse"
                }
            }
        },
        "dto.ResetPasswordRequest": {
            "type": "object",
            "required": [
                "new_password",
                "token"
            ],
            "properties": {
                "new_password": {
                    "type": "string",
                    "minLength": 8
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "dto.UserResponse": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/03622bf7-d58b-4997-965c-14ee58c63554/",
	Schemes:          []string{},
	Title:            "Authentication API",
	Description:      "A JWT Authentication Service API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
