# Migration

!!! Important

    In version 2.0.0 the database backend was changed from [Storm](https://github.com/asdine/storm) to [Gorm](https://gorm.io/). This means you have to migrate 
    your database. The following steps will guide you through the process.

## Migrate from Storm to Gorm

!!! Information

    Ensure you have the latest version of the ticker (`>= 2.0.0`) service installed.

1. Stop the ticker service
2. Configure the new database backend in `config.yml`

        database:
            type: "postgres" # postgres, mysql, sqlite
            dsn: "host=localhost port=5432 user=ticker dbname=ticker password=ticker sslmode=disable"

3. Start the migration

        ticker db migrate --storm.path /path/to/storm.db

4. Start the ticker service
