# Supabase

## What

Supabase is an outside service that we will use for authentication and our database.

## Getting Started Locally

You must first install the Supabase CLI [here](https://supabase.com/docs/guides/local-development/cli/getting-started?queryGroups=access-method&access-method=postgres&queryGroups=platform&platform=macos)

To run supabase locally, you must have docker installed and running in the background. [Docker Desktop (GUI)](https://docs.docker.com/desktop/install/) [Docker Engine (CLI)](https://docs.docker.com/engine/install/)

Once you have docker installed and running, from the command line run:

```
cd backend/internal/supabase
supabase start
```

Once that has completed, you can check if supabase is running locally by running `supabase status`. You can also view your local instance from the supabase studio at [http://127.0.0.1:54323](http://127.0.0.1:54323)

## Creating a New Migration

To create a new migration _do not_ simply add a file to the migrations folder. Please instead run the following command:

```
supabase migrations new your-migration-name
```

Please add a migration name that is descriptive yet concise. If you believe that your name does not cover the full scope of the migration add a comment to the top of the file so others understand the high level changes made.

To apply the migration to your local database, run `supabase db reset` which will apply all unapplied database migrations. You should then see the updated changes at [http://127.0.0.1:54323](http://127.0.0.1:54323).

## Applying Migrations to Prod

If and only if your PR has been merged, you can apply migrations using the supabase CLI with the following steps:

1. Run `supabase link`
2. Select Inspirate Consulting
3. If prompted, enter the DB password (slack a TL)
4. If prompted, log in to Supabase
5. Apply changes to prod with `supabase db push`

You should check the prod Supabase Studio to check and make sure your changes are applied.

## Why Supabase?

Supabase has a lot of built in features that you get for free like auth, edge functions, an always-active database, localization and much more. Using it allows for us to focus our time on the product and delivering for clients rather than fighting with infra.