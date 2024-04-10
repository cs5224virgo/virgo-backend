# Virgo-backend

Virgo chat app backend written in go

## Installation

Do these following steps to set up the project locally:

1. Clone the project && cd into the project
1. Install docker && docker compose ([docs](https://docs.docker.com/engine/install/)). Make sure you can run `docker` and `docker compose`.
1. Run the real thing: `docker compose up --build`

If you are aiming to develop / add features, also do these things:

1. Install sqlc to write database queries: `brew install sqlc` (macos) / `sudo snap install sqlc` (ubuntu) -- ([docs](https://docs.sqlc.dev/en/stable/overview/install.html))
1. `go generate` -- This will run sqlc to generate db code

## Handling migrations / DB workflow

Q: How to add table / model / schema?

A: 

```
touch db/migrations/{three_digit_id}_{name}.up.sql

vim $_

go generate
```

When adding migrations, remember to add both up and down migrations.

Put migration scripts in `/db/migrations` and put queries that get/update/delete in `/db/queries`.

After writing sql files, run `go generate` to generate go code from those sql files. Generated code goes into `/db/generated`. When writing code in `/internal/*`, import and use those generated code to interact with db.

# Dev notes

## Folder structures

```
.
├── config/                 # Config profiles for local and prod
├── cmd/                    # Entrypoint / commands to run different parts of the program
│   ├── root.go             # root.go is the mother of all entrypoints
│   ├── migrate.go          
│   ├── yahoo.go            
│   └── ...                 
├── db/                     # All primitive DB stuff goes here
│   ├── generated/          # Everything generated by sqlc goes here
│   │   ├── db.go           # Actually useless. EDIT: my bad its not useless
│   │   ├── models.go       # Import this file for all the models
│   │   ├── *.sql.go         
│   │   └── ...          
│   ├── migrations/         # Write all your migrations / schema here
│   │   └── ...          
│   ├── queries/            # Write the most basic CRUD here. For other complex DB stuff,
│   │   │                   # put them in internal/datalayer instead, building on top of the CRUDs.
│   │   ├── crawl_job.go    # Categorize queries based on their respective models.
│   │   ├── crawl_job_frontend.go    # If there are many queries, subcategorize based on purpose.
│   │   └── ...          
│   ├── db.go               # Code to connect to db
│   └── migrate.go          # Code to migrates the db
│
├── internal/               # There is actual code inside here, can you believe it?
│   ├── datalayer/          # Database-agnostic code. Everything that's more advanced than pure CRUD.
│   │   └── .../            # Split into their own domains of responsibilty.
│   └── services/           # Only business logic goes in here. Codes that makes the magic happen. 
│       └── .../            # Split into their own domains of responsibilty.
└── ...
```

Please follow naming scheme wherever possible.

## Logging

Use `logger` (`import "github.com/cs5224virgo/virgo/logger"`) to do logging.

- `logger.Info()` to write out debug info, data dump, or code checkpoints. `logger.Infof()` if you need printf-like feature. `logger.Infoln()` if you need more newlines.
- `logger.Warn()` `logger.Warnf()` to write out ignorable errors or generally unexpected situations (but still safe to continue).
- `logger.Error()` `logger.Errorf()` when the error is serious and needs fixing but doesnt break the flow of the program. Usually use before returning.
- `logger.Fatal()` `logger.Fatalf()` log and then panic. Use sparingly. Only non-recoverable errors.

