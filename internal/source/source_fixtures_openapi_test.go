package source

const validOpenAPIWithMissingMetadata = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "post": {
        "operationId": "createCustomer",
        "responses": {
          "201": {"description": "created"}
        }
      }
    }
  }
}`

const openAPIWithBodylessDelete = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers/{id}": {
      "delete": {
        "operationId": "deleteCustomer",
        "summary": "Delete customer",
        "description": "Delete one customer.",
        "parameters": [
          {"name": "id", "in": "path", "description": "Customer ID"}
        ],
        "responses": {
          "204": {
            "description": "Deleted.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    }
  }
}`

const openAPIWithNullMediaSchemas = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "post": {
        "operationId": "createCustomer",
        "summary": "Create customer",
        "description": "Create one customer.",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": null
            }
          }
        },
        "responses": {
          "201": {
            "description": "Created.",
            "content": {
              "application/json": {
                "schema": null
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    }
  }
}`

const yamlOpenAPIWithNullMediaSchemas = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody:
        content:
          application/json:
            schema: null
      responses:
        "201":
          description: Created.
          content:
            application/json:
              schema:
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const openAPIWithEmptyOAuthScopeDescription = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "get": {
        "operationId": "listCustomers",
        "summary": "List customers",
        "description": "List customers.",
        "responses": {
          "200": {
            "description": "Customers.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "oauth": {
        "type": "oauth2",
        "flows": {
          "authorizationCode": {
            "authorizationUrl": "https://example.com/auth",
            "tokenUrl": "https://example.com/token",
            "scopes": {
              "customers:read": ""
            }
          }
        }
      }
    }
  }
}`

const openAPIWithEmptySecurityScheme = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "get": {
        "operationId": "listCustomers",
        "summary": "List customers",
        "description": "List customers.",
        "responses": {
          "200": {
            "description": "Customers.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "ApiKeyAuth": {}
    }
  }
}`

const yamlOpenAPIWithEmptyOAuthScopeDescription = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    oauth:
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: https://example.com/auth
          tokenUrl: https://example.com/token
          scopes:
            customers:read: ""
`

const yamlOpenAPIWithLateOAuthType = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    oauth:
      flows:
        authorizationCode:
          authorizationUrl: https://example.com/auth
          tokenUrl: https://example.com/token
          scopes:
            customers:read: ""
      type: oauth2
`

const yamlOpenAPIWithInlineOAuthScopes = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    oauth:
      type: oauth2
      flows:
        authorizationCode:
          authorizationUrl: https://example.com/auth
          tokenUrl: https://example.com/token
          scopes: {customers:read: ""}
`

const yamlOpenAPIWithFlowStyleOAuthSecuritySchemes = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes: {oauth: {type: oauth2, flows: {authorizationCode: {authorizationUrl: https://example.com/auth, tokenUrl: https://example.com/token, scopes: {customers:read: ""}}}}}
`

const yamlOpenAPIWithBlockFlowStyleOAuthSecurityScheme = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    oauth: {type: oauth2, flows: {authorizationCode: {authorizationUrl: https://example.com/auth, tokenUrl: https://example.com/token, scopes: {customers:read: ""}}}}
`

const openAPIWithUndescribedParameterRef = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers/{id}": {
      "parameters": [
        {"$ref": "#/components/parameters/CustomerId"}
      ],
      "get": {
        "operationId": "getCustomer",
        "summary": "Get customer",
        "description": "Get one customer.",
        "responses": {
          "200": {
            "description": "Customer.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    },
    "parameters": {
      "CustomerId": {
        "name": "id",
        "in": "path"
      }
    }
  }
}`

const openAPIWithChainedUndescribedParameterRef = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers/{id}": {
      "parameters": [
        {"$ref": "#/components/parameters/CustomerIdAlias"}
      ],
      "get": {
        "operationId": "getCustomer",
        "summary": "Get customer",
        "description": "Get one customer.",
        "responses": {
          "200": {
            "description": "Customer.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    },
    "parameters": {
      "CustomerIdAlias": {"$ref": "#/components/parameters/CustomerId"},
      "CustomerId": {
        "name": "id",
        "in": "path"
      }
    }
  }
}`

const openAPIWithOperationParameterOverride = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers/{id}": {
      "parameters": [
        {"name": "id", "in": "path"}
      ],
      "get": {
        "operationId": "getCustomer",
        "summary": "Get customer",
        "description": "Get one customer.",
        "parameters": [
          {"name": "id", "in": "path", "description": "Customer ID"}
        ],
        "responses": {
          "200": {
            "description": "Customer.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    }
  }
}`

const openAPIWithDuplicateSamePathNames = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "get": {
        "operationId": "listCustomers",
        "summary": "Customers",
        "description": "List customers.",
        "responses": {
          "200": {
            "description": "Customers.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      },
      "post": {
        "operationId": "createCustomer",
        "summary": "Customers",
        "description": "Create customer.",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {"type": "object"}
            }
          }
        },
        "responses": {
          "201": {
            "description": "Created.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    }
  }
}`

const openAPIWithDuplicateOperationIDsDistinctSummaries = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "get": {
        "operationId": "syncCustomer",
        "summary": "List customers",
        "description": "List customers.",
        "responses": {
          "200": {
            "description": "Customers.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    },
    "/customers/{id}/sync": {
      "post": {
        "operationId": "syncCustomer",
        "summary": "Sync one customer",
        "description": "Sync one customer.",
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {"type": "object"}
            }
          }
        },
        "responses": {
          "200": {
            "description": "Synced.",
            "content": {
              "application/json": {
                "schema": {"type": "object"}
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    }
  }
}`

const yamlOpenAPIWithUndescribedParameterRef = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    parameters:
      - $ref: "#/components/parameters/CustomerId"
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      responses:
        "200":
          description: Customer.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  parameters:
    CustomerId:
      name: id
      in: path
`

const yamlOpenAPIWithUndescribedFlowParameter = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      parameters:
        - { name: id, in: path }
      responses:
        "200":
          description: Customer.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithUndescribedFlowParameterRef = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      parameters:
        - { $ref: "#/components/parameters/CustomerId" }
      responses:
        "200":
          description: Customer.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  parameters:
    CustomerId:
      name: id
      in: path
`

const completeOpenAPIYAML = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers with pagination.
      parameters:
        - name: page
          in: query
          description: Page number.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody:
        content:
          application/json:
            schema:
              type: object
      responses:
        "201":
          description: Created.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOperationMissingDescriptionWithResponseDescription = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      responses:
        "200":
          description: Response object description should not describe the operation.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlSuccessResponseWithoutSchemaAndErrorSchema = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      responses:
        "200":
          description: Customer read-back without schema.
        "404":
          description: Missing customer error.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlSchemaPropertyNamedLikeHTTPMethod = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /reports:
    get:
      operationId: getReport
      summary: Get report
      description: Get one report.
      responses:
        "200":
          description: Report response.
          content:
            application/json:
              schema:
                type: object
                properties:
                  get:
                    type: string
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlNestedPathsSchemaPropertyBeforeLaterOperation = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    get:
      operationId: listCustomers
      summary: List customers
      description: List customers.
      responses:
        "200":
          description: Customers.
          content:
            application/json:
              schema:
                type: object
                properties:
                  paths:
                    type: array
  /orders:
    get:
      summary: List orders
      description: List orders.
      responses:
        "200":
          description: Orders.
          content:
            application/json:
              schema:
                type: object
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`
