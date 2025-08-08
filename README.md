
# CRUDProdutosMongoAndGin
It's a simple application created with Go and library Gin for CRUD operations on products.

## Tech Stack Used
**Back-end:** Go and Gin

**Database:** MongoDB

## Environment Variables
To run your project, it's necessary the following environment variables on your `.env` file:

```
GO_ENV=[development|production]
MONGODB_CONNECTION_URI=[URI CONNECTION TO MONGODB]
PORT=[ACCESS PORT TO APPLICATION]
```

> Those variables are actually present on `.env.base`, being only necessary to copy the file to a `.env` file.

## Running locally

### Running on Docker environment
> To run the project with a Docker environment, it's mandatory install *Docker Toolkit*.

With your Docker Engine started, run the following command being on project's root dir:

```bash
  docker compose up --build -d
```
> That command will build your application and run the binary file generated.

That way, with `.env` file configurated accordingly, you will be able to access the application on defined `PORT`.

### Standard Run
> It's mandatory to install Go in your local environment.

It's needed first to install your go project dependency:

```bash
go mod download
```

And build your project:

```bash
go build -o main
```

Finally, being able to run your binary with:

```bash
./main
```

## API Documentation
Swagger documentation is present on `/swagger/index.html` application route, containing all the routes standardized with OpenAPI 3.0.

If you are running locally, access with the link: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)