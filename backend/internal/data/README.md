# Data Layer

This folder is our data access layer. Any layers above this such as the handlers will call methods on a Store object, and the Store objects are created in this layer. 

## Folder Structure

**The `/postgres` folder:**
This folder has the implementation
The `postgres/storage.go` file handles database connection setup. We build a pgx connection pool using the configuration in `internal/config.db`. Then `storage.go` pings the database, and returns it.
`postgres/schema/*` are where the database queries will live. We create packages based off of the entity they relate to. Every package contains a `store.go` file which defined the struct and constructor. Then you can make additional `<operation>.go` files for each operation, like `schema/greeting/create.go` to create a greeting. 
Talk about SQL:

## Concepts

### pgx

This is the Go library we are using to talk to PostgreSQL. We pass in SQL strings and arguments to pgx, and pgx sends the arguments to our database, and returns rows back. So pgx is what opens a connection to Postgres, and communicates our arguments and responses to and from our database.

### pgxpool

One connection to postgres can only run one query at a time, so we need to use pgxpool to create a connection pool, so we can run multiple queries concurrently. `pgxpool` handles a pool of connections that can be reused, and controls configurations like how many connections to keep open, etc.