MIGRATIONS_DIR=./migrations



migrate_create:
	migrate create -seq -ext=.sql -dir=$(MIGRATIONS_DIR) $(name)

migrate_up:
	migrate -path=$(MIGRATIONS_DIR) -database=$(DSN_DATABASE) up
