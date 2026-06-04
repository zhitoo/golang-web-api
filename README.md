## make artisan 
```
go build -o artisan artisan.go
```
## create migration file
```
./artisan migrate:create migration_file_name
```
## migrate up
```
./artisan migrate:up
```

## migrate down
```
./artisan migrate:down 1[steps]
```