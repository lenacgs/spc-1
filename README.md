# SIBS Performance Challenge 1

## Build and run

Requirements:

- Golang go1.22+

Compile and install:
```bash
make clean && make
# build will be available at ./build/
# builds workersd, serviced and client
```

Start Workers API service:
```bash
cd build
./workersd -d . -o
# -d . writes logs and lock at current directory
# -o logs to stdout
```

Start Service API service:
```bash
cd build
./serviced -d . -o
# -d . writes logs and lock at current directory
# -o logs to stdout
```

## Client

### Launch client:

```bash
cd build
./client
```

### Interact with client:

Connect to workersd:

```shell
connect
```

or

```shell
connect localhost:8080
```

Connect to serviced:
```shell
connect localhost:8081
```

Put Worker with ID 1 and Location 2:

```shell
put 1 2
```

Get Worker with ID 1:

```shell
get 1
```

Delete Worker with ID 1:

```shell
del 1
```

Make request to Service API (after connecting to localhost:8081)
```shell
#gets worker with ID 1
cache 1
```
