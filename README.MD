# dbds (Database Develepment Seeder)

A simple database seeding tool to quickly build your development database.

(Not to use in production.)

## Overview

This project provides a command-line tool to seed a database with predefined data. It supports two main commands:

- `dbds init`: Initializes dbds config and base sql script files.
- `dbds rebuild`: Rebuilds the database based on the configuration, which can include dropping and recreating tables and re-inserting data.

## Get dbds

### Download prebuilt binaries

-

### Or build from source

To build this tool from source, you'll need Go installed on your system.

You can download Go from [the official Go website](https://golang.org/dl/).

You can build from source with: `go build`.

## Usage

### Initialize dbds config

To initialize your database seed config, use the following command:

```
dbds init
```

This will create a configuration file `dbds.cfg` and `dbds_script` directory with sql scripts in it.

### Usage of dbds.cfg

Type your **databse connection string** and **database type** in dbds.cfg configuration file.

For example:

```
connectionString:"host=localhost user=db_dev_user password=db_dev_user_password dbname=database_name"
dbType:"postgres"
```

### Usage of dbds_scripts

The `.sql` scripts under `dbds_script` are the SQL statements to be run during `dbds rebuild`.

The `.sql` scripts are run in this order:

1. **1_drop_tables.sql**  
   This script should contain statements to drop the existing tables from the database. It ensures that the schema is clean before creating new tables.

2. **2_create_tables.sql**  
   This script should contain statements to create the necessary tables for the database schema. It defines the structure of the database by creating tables and their relationships.

3. **3_create_indexes.sql**  
   This script should contain statements to create indexes on the tables to optimize query performance. Indexes help speed up data retrieval operations.

4. **4_populate_data.sql**  
   This script should containt statements to populate the tables with initial data. It inserts predefined data that is needed for the application to function correctly or for testing purposes.
