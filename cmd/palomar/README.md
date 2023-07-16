# Palomar

Palomar is an Elasticsearch/OpenSearch/Meilisearch frontend and ATP (AT Protocol) repository crawler designed to provide search services for the Bluesky network.

## Prerequisites

- GoLang (version 1.20)
- Running instance of Elasticsearch or OpenSearch or Meilisearch for indexing.

## Building

```
go build
```

## Configuration

Palomar uses environment variables for configuration.

- `ATP_BGS_HOST`: URL of the Bluesky BGS (e.g., `https://bgs.staging.bsky.dev`).
- `SEARCH_ENGINE_TYPE`: "elastic" or "meili"
- `ELASTIC_HTTPS_FINGERPRINT`: Required if using a self-signed cert for your Elasticsearch deployment.
- `ELASTIC_USERNAME`: Elasticsearch username (default: `admin`).
- `ELASTIC_PASSWORD`: Password for Elasticsearch authentication.
- `MEILISEARCH_APIKEY`: apikey of meilisearch.
- `SEARCH_ENGINE_HOSTS`: Comma-separated list of search engine endpoints.
- `READONLY`: Set this if the instance should act as a readonly HTTP server (no indexing).
- `DATABASE_URL`: Url of sqlite or postgresql
- `ADMIN_KEY`: Default admin password
- `ATP_PDS_HOST`: method, hostname, and port of PDS instance

For PostgreSQL, the user and database must already be configured. Some example
SQL commands are:

    CREATE DATABASE palomar;

    CREATE USER ${username} WITH PASSWORD '${password}';
    GRANT ALL PRIVILEGES ON DATABASE palomar TO ${username};


## Running the Application

Once the environment variables are set properly, you can start Palomar by running:

```
./palomar run
```

## Indexing 
For now, there isnt an easy way to get updates from the PDS, so to keep the
index up to date you will periodcally need to scrape the data.

## API

### `/index/:did`
Indexes the content in the given user's repository. It keeps track of the last repository update and only fetches incremental changes.

### `/search?q=QUERY`
Performs a simple, case-insensitive search across the entire application.
