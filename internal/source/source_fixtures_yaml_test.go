package source

const swagger2JSONFixture = `{
  "swagger": "2.0",
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
            "schema": {"type": "object"}
          }
        }
      }
    }
  }
}`

const yamlOpenAPIWithNestedSwaggerProperty = `openapi: 3.0.3
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
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  schemas:
    Customer:
      type: object
      properties:
        swagger:
          type: string
`

const yamlOpenAPIWithNestedOpenAPIPropertyOnly = `info:
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
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  schemas:
    Metadata:
      type: object
      properties:
        openapi:
          type: string
`

const yamlEmptySecuritySchemes = `openapi: 3.0.3
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
  securitySchemes: {}
`

const yamlOpenAPIWithEmptyNamedSecurityScheme = `openapi: 3.0.3
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
    apiKeyAuth: {}
`

const yamlNestedSecuritySchemesProperty = `openapi: 3.0.3
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
                  securitySchemes:
                    type: string
`

const yamlNestedSecuritySchemesSchemaPropertyOnly = `openapi: 3.0.3
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
  schemas:
    Agent:
      type: object
      properties:
        securitySchemes:
          type: object
          properties:
            apiKeyAuth:
              type: string
`

const yamlResponseHeaderSchemaOnly = `openapi: 3.0.3
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
          headers:
            X-Request-ID:
              schema:
                type: string
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlRequestBodyNestedExampleSchemaOnly = `openapi: 3.0.3
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
            examples:
              one:
                value:
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

const yamlFlowStyleMediaNestedExampleSchemaOnly = `openapi: 3.0.3
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
          application/json: { examples: { e: { value: { schema: {} } } } }
      responses:
        "200":
          description: Customer.
          content:
            application/json: { examples: { e: { value: { schema: {} } } } }
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithInlineRequestAndResponseContentSchemas = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody: {content: {application/json: {schema: {type: object}}}}
      responses:
        "201": {description: Created, content: {application/json: {schema: {type: object}}}}
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithFlowStyleOperationSchemas = `openapi: 3.0.3
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
          application/json: { schema: { type: object } }
      responses:
        "201":
          description: Created.
          content:
            application/json: { schema: { type: object } }
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithInlineOperationObject = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post: {operationId: createCustomer, summary: Create customer, description: Create one customer., requestBody: {content: {application/json: {schema: {type: object}}}}, responses: {"201": {description: Created., content: {application/json: {schema: {type: object}}}}}}
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithRootFlowStylePaths = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths: {/customers: {get: {operationId: listCustomers, summary: List customers, description: List customers., responses: {"200": {description: Customers., content: {application/json: {schema: {type: object}}}}}}}}
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithFlowStylePathItem = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers: {get: {operationId: listCustomers, summary: List customers, description: List customers., responses: {"200": {description: Customers., content: {application/json: {schema: {type: object}}}}}}}
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithDeprecatedReplacementHint = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /v1/customers:
    get:
      operationId: listLegacyCustomers
      summary: Legacy customers
      description: Legacy endpoint.
      deprecated: true
      x-deprecated-replacement: listCustomers
      responses:
        "200":
          description: Customers.
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

const yamlOpenAPIWithFlowStyleComponentGroups = `openapi: 3.0.3
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
        $ref: "#/components/requestBodies/CustomerWrite"
      responses:
        "201":
          $ref: "#/components/responses/CustomerRead"
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  requestBodies: {CustomerWrite: {content: {application/json: {schema: {type: object}}}}}
  responses: {CustomerRead: {description: Customer., content: {application/json: {schema: {type: object}}}}}
`

const openAPIWithComponentRequestAndResponseRefs = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "post": {
        "operationId": "createCustomer",
        "summary": "Create customer",
        "description": "Create one customer.",
        "requestBody": {"$ref": "#/components/requestBodies/CustomerWrite"},
        "responses": {
          "201": {"$ref": "#/components/responses/CustomerRead"}
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    },
    "requestBodies": {
      "CustomerWrite": {
        "content": {
          "application/json": {
            "schema": {"type": "object"}
          }
        }
      }
    },
    "responses": {
      "CustomerRead": {
        "description": "Customer read-back.",
        "content": {
          "application/json": {
            "schema": {"type": "object"}
          }
        }
      }
    }
  }
}`

const openAPIWithChainedComponentRequestAndResponseRefs = `{
  "openapi": "3.0.3",
  "info": {"title": "Fixture API", "version": "1.0.0"},
  "paths": {
    "/customers": {
      "post": {
        "operationId": "createCustomer",
        "summary": "Create customer",
        "description": "Create one customer.",
        "requestBody": {"$ref": "#/components/requestBodies/CustomerWriteAlias"},
        "responses": {
          "201": {"$ref": "#/components/responses/CustomerReadAlias"}
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "apiKeyAuth": {"type": "apiKey", "in": "header", "name": "X-API-Key"}
    },
    "requestBodies": {
      "CustomerWriteAlias": {"$ref": "#/components/requestBodies/CustomerWrite"},
      "CustomerWrite": {
        "content": {
          "application/json": {
            "schema": {"type": "object"}
          }
        }
      }
    },
    "responses": {
      "CustomerReadAlias": {"$ref": "#/components/responses/CustomerRead"},
      "CustomerRead": {
        "description": "Customer read-back.",
        "content": {
          "application/json": {
            "schema": {"type": "object"}
          }
        }
      }
    }
  }
}`

const yamlOpenAPIWithComponentResponseRef = `openapi: 3.0.3
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
          $ref: "#/components/responses/CustomerList"
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  responses:
    CustomerList:
      description: Customer list.
      content:
        application/json:
          schema:
            type: object
`

const yamlOpenAPIWithInlineComponentRefs = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody: { $ref: "#/components/requestBodies/CustomerWrite" }
      responses:
        "201": { $ref: "#/components/responses/CustomerRead" }
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  requestBodies:
    CustomerWrite:
      content:
        application/json:
          schema:
            type: object
  responses:
    CustomerRead:
      description: Customer read-back.
      content:
        application/json:
          schema:
            type: object
`

const yamlOpenAPIWithQuotedInlineComponentRefs = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody: { "$ref": "#/components/requestBodies/CustomerWrite" }
      responses:
        "201": { "$ref": "#/components/responses/CustomerRead" }
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  requestBodies:
    CustomerWrite:
      content:
        application/json:
          schema:
            type: object
  responses:
    CustomerRead:
      description: Customer read-back.
      content:
        application/json:
          schema:
            type: object
`

const yamlOpenAPIWithJSONStyleFlowComponentsAndNestedComponentsProperty = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers:
    post:
      operationId: createCustomer
      summary: Create customer
      description: Create one customer.
      requestBody: {"$ref":"#/components/requestBodies/CreateCustomer"}
      responses:
        "200": {"$ref":"#/components/responses/Customer"}
components:
  schemas:
    Customer:
      type: object
      properties:
        components:
          type: string
  requestBodies:
    CreateCustomer:
      content:
        application/json: {"schema":{"type":"object"}}
  responses:
    Customer:
      description: Customer.
      content:
        application/json: {"schema":{"type":"object"}}
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithFlowStyleComponentSchemas = `openapi: 3.0.3
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
        $ref: "#/components/requestBodies/CustomerWrite"
      responses:
        "201":
          $ref: "#/components/responses/CustomerRead"
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
  requestBodies:
    CustomerWrite:
      content:
        application/json: { schema: { type: object } }
  responses:
    CustomerRead:
      description: Customer read-back.
      content:
        application/json: { schema: { type: object } }
`

const yamlOpenAPIWithColonPathKey = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /v1/books:batchGet:
    get:
      operationId: batchGetBooks
      responses:
        "200":
          description: Books.
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

const yamlOpenAPIWithUndescribedPathParameter = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    parameters:
      - name: id
        in: path
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
`

const yamlOpenAPIWithLateUndescribedPathParameter = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
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
    parameters:
      - name: id
        in: path
components:
  securitySchemes:
    apiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
`

const yamlOpenAPIWithInlineOperationParameterArray = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    get:
      operationId: getCustomer
      summary: Get customer
      description: Get one customer.
      parameters: [{name: id, in: path}]
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

const yamlOpenAPIWithInlinePathParameterArray = `openapi: 3.0.3
info:
  title: Fixture API
  version: 1.0.0
paths:
  /customers/{id}:
    parameters: [{name: id, in: path}]
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
`

const completeDocs = `# API Guide

Use the API key header for auth.
Retry 429 responses with backoff.
Respect the rate limit.
Use pagination with the page parameter.
Create requests are idempotent when an idempotency key is supplied.
`
