# Echo Swagger Middleware

This module provides middleware for serving Swagger UI and API documentation in an Echo web framework application. It supports serving Swagger JSON files, the Swagger UI, and static files.

## Installation

To use this module, install the required dependencies:

```bash
go get github.com/labstack/echo/v4
go get github.com/swaggo/echo-swagger

```

## Usage

Below is an example of how to use this module to serve Swagger documentation for a `blue-auth` service:

### Example Code

```go
package main

import (
    "github.com/labstack/echo/v4"
    echoSwagger "github.com/bushubdegefu/echo-swagger"
)

func main() {
    app := echo.New()

    // Serve the Swagger JSON file
    app.GET("/blue_auth/docs/doc.json", func(contx echo.Context) error {
        return contx.File("blue-auth/docs/swagger.json")
    }).Name = "blue_auth_docs_json"

    // Serve the Swagger UI
    app.GET("/blue_auth/docs/*", echoSwagger.New(echoSwagger.Config{
        InstanceName: "blue_auth",
        URL:          "/blue_auth/docs/doc.json", // Match the served JSON file
    })).Name = "blue_auth_docs"

    // Start the server
    app.Logger.Fatal(app.Start(":8080"))
}
```

### Directory Structure

Ensure your project directory is structured as follows:

```
project/
├── blue-auth/
│   └── docs/
│       └── swagger.json
├── main.go
```

### Swagger JSON File

The `swagger.json` file should contain the Swagger documentation for your API. You can generate this file using tools like [Swag](https://github.com/swaggo/swag).

### Running the Application

Run the application with:

```bash
go run main.go
```

### Accessing the Documentation

- Swagger JSON: [http://localhost:8080/blue_auth/docs/doc.json](http://localhost:8080/blue_auth/docs/doc.json)
- Swagger UI: [http://localhost:8080/blue_auth/docs/](http://localhost:8080/blue_auth/docs/)

## Configuration

The `New` function accepts a `Config` struct to customize the behavior of the middleware. For example:

```go
echoSwagger.New(echoSwagger.Config{
    InstanceName: "blue_auth",
    URL:          "/blue_auth/docs/doc.json",
})
```

### Config Options

- `InstanceName`: The name of the Swagger instance.
- `URL`: The URL of the Swagger JSON file.

## License
This module is licensed under the MIT License.