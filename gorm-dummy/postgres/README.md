# GORM DUMMY

A simple app to get you started with testing GORM if you've never used it before.

Prerequisites
-------------

You need [docker](https://docs.docker.com/get-docker/) in order to run the example, and potentially [docker-compose](https://docs.docker.com/compose/install/) if it doesn't come by default with docker.

Setup
-----

There are two components here, the database and the Go app.

**Database**

I'm using [Postgres](https://www.postgresql.org/) as the Relational Database engine. You can find a great PostgreSQL tutorial [here](https://www.postgresqltutorial.com/). This tutorial also contains a dummy database with plenty of records you can practice your SQL on: check out [this link.](https://www.postgresqltutorial.com/postgresql-getting-started/postgresql-sample-database/)

First, you should set up the database. At the moment, it contains a single basic entity defined in [scripts](./scripts/). It should have enough fields to test easily most of the functionalities. Later this may include a model with multiple relations.

Start the database with:

```sh
docker compose up
# use -d flag to run in detached mode.
```

The first time you run it it should automatically execute the SQL script mentioned earlier. This has been implemented using the mechanism described [here](https://hub.docker.com/_/postgres) (look out for `Initialization scripts` and `initdb.d`).

Stop the database with:

```sh
docker compose down
```

Note that this command will not remove the data that was inserted the first time running it, or changed by your code. In order to remove the data as well use the command:

```sh
docker compose down -v
```

Starting the container again will recreate the database based on the SQL script again.

> NB: It is not recommended your database in containers in production: https://vsupalov.com/database-in-docker/

**Go app**

Now that you have the database set up, simply test the database by running the Go app. By default, it simply orders the entities, but you should add whichever functionality you want to test out. Start the app with:

```sh
go run main.go
```

The program will connect to the db, execute the code and then exit.

Sample run
----------

Starting the db:

```sh
$ docker compose up
[+] Running 2/2
 - Network gorm-dummy_default  Created
 - Container dummy_gin_db      Created
Attaching to dummy_gin_db
dummy_gin_db  | 
dummy_gin_db  | PostgreSQL Database directory appears to contain a database; Skipping initialization
dummy_gin_db  | 
dummy_gin_db  | 2023-04-13 14:28:21.506 UTC [1] LOG:  starting PostgreSQL 15.2 (Debian 15.2-1.pgdg110+1) on x86_64-pc-linux-gnu, compiled by gcc (Debian 10.2.1-6) 10.2.1 20210110, 64-bit
dummy_gin_db  | 2023-04-13 14:28:21.507 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
dummy_gin_db  | 2023-04-13 14:28:21.507 UTC [1] LOG:  listening on IPv6 address "::", port 5432
dummy_gin_db  | 2023-04-13 14:28:21.515 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
dummy_gin_db  | 2023-04-13 14:28:21.524 UTC [29] LOG:  database system was shut down at 2023-04-13 14:28:16 UTC
dummy_gin_db  | 2023-04-13 14:28:21.530 UTC [1] LOG:  database system is ready to accept connections
```

Running the app:

```
$ go run main.go
...
[2.626ms] [rows:5] SELECT * FROM "entities" ORDER BY name DESC,salary ASC,uuid ASC LIMIT 5
UUID: 5f875d48-c925-40e7-96ae-2d7ee48e685f      NAME: F  OTHER_NAME: F   AGE: 6  SALARY: 3000
UUID: 3be67cd8-476c-4e4d-8f64-141ff3ec1ea1      NAME: F  OTHER_NAME: C   AGE: 7  SALARY: 4000
UUID: 0b402f7f-bc58-4098-9b56-996f7bca92e7      NAME: F  OTHER_NAME: F   AGE: 3  SALARY: 5000
UUID: 5ba21827-7bd1-46c1-b2f0-b8b1022766b7      NAME: E  OTHER_NAME: B   AGE: 9  SALARY: 3000
UUID: 5d448886-8e52-4644-9ea2-cf2f7b99f4d7      NAME: E  OTHER_NAME: E   AGE: 5  SALARY: 3000
```

If ran in detached mode, stop the db:

```
$ docker compose down -v
[+] Running 3/3
 - Container dummy_gin_db        Removed 
 - Volume gorm-dummy_testvolume  Removed
 - Network gorm-dummy_default    Removed 
```