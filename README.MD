# dbds (Database Develeopment Seeder)

A simple database seeding tool to quickly build your development database.

## Overview

This project provides a command-line tool to seed a database with predefined data. It supports two main commands:

- `dbds init`: Initializes dbds config and base sql script files.
- `dbds rebuild`: Rebuilds the database based on the configuration, which can include dropping and recreating tables and re-inserting data.

## Get dbds

### Download prebuilt binaries

-

### Or build from source

To build this tool from source, you'll need Go installed on your system. You can download Go from [the official Go website](https://golang.org/dl/).
You can build from source with: `go build`.

## Usage

### Initialize the Database

To initialize your database, use the following command:

```sh
./dbds init
