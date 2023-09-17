Work in progress

# Database Table Migration

We use golang-migrate for DB schema changes. To generate a set of migrate up and down sql script:

- Run "migrate create -seq -ext=.sql -dir=<directory_path> -seq <filename>"

After filing up the scripts with your DB commands, run the following to migrate

- "migrate -path=directory_path> -database='postgres://dev:honeybbee8988@localhost:5432/stable_wallet?sslmode=disable' up"
