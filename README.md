# Roster

> A demo project for understanding Go kit

## Installation

1. Install PostgreSQL
2. Create a database called "roster"

```sh
$ createdb roster
```

3. Install Go, following instructions at <https://golang.org/doc/install> or, for Mac, you can just use homebrew:

```sh
$ brew install go
$ brew install protobuf
$ # Add the below line to your appropriate "rc" file for your shell: (bashrc, zshrc, .fishrc, etc)
$ export GOPATH=$HOME/go; export PATH=$PATH:$GOPATH/bin

```

*Note:* Older installation guides indicate setting more environment variables than necessary. The page at <https://github.com/golang/go/wiki/SettingGOPATH> contains current information on Go's environment variables.

4. Follow instructions at <https://grpc.io/docs/quickstart/go.html> to install gRPC and Protocol Buffers 3. Note that, for the Protocol Buffers step, if you use homebrew on a Mac you can use:

```sh
$ brew install protobuf
```

5. Get the code

```sh
$ go get -u github.com/hoop33/roster
$ cd $GOPATH/src/github.com/hoop33/roster
$ make deps
$ dep ensure -update
$ make
```

6. Ensure postgress is working correctly.

```sh
$ psql #logs in as super user
```
```psql
\connect roster
psql=# CREATE USER <db user> WITH SUPERUSER PASSWORD <'password'>;
psql=#\q
```

7. Run the app, which will create the `players` table

```sh
$ ROSTER_USER=<db user> ROSTER_PASSWORD=<db password> ./roster
```

### (Optional) Seed the Database

1. Follow instructions to install <https://github.com/hoop33/jags>
2. `$ jags | sed 's/,,/,N\/A,/g' | sed 's/,R,/,0,/g' > players.csv`
3. Run the following SQL:

```sql
COPY players(name,number,position,height,weight,age,experience,college) 
FROM '<path to file>/players.csv' DELIMITER ',' CSV HEADER;
```

## Troubleshooting

If you see an error around `grpc` libraries, e.g., :

```sh
package google.golang.org/grpc/grpclb/grpc_lb_v1/messages: cannot find package "google.golang.org/grpc/grpclb/grpc_lb_v1/messages" in any of:
	/usr/local/Cellar/go/1.10.3/libexec/src/google.golang.org/grpc/grpclb/grpc_lb_v1/messages (from $GOROOT)
	/Users/<user>/go/src/google.golang.org/grpc/grpclb/grpc_lb_v1/messages (from $GOPATH)
```

You may have a mismatch between the `grpc` libraries in your vendored files and the ones that `protoc` is picking up. To get them in sync, try getting the latest for both:

```sh
$ go get -u google.golang.org/grpc
$ cd $GOPATH/src/github.com/hoop33/roster
$ dep ensure -update google.golang.org/grpc
```

See a discussion at <https://github.com/grpc/grpc-go/issues/581>

## Walkthrough

For each step, check out the tag, build the app, and run:

```sh
$ git checkout <tag>
$ make && ./roster
```

### Step 1: A Simple Service (step_1)

Nothing Go kit-related here; just a simple Go command-line service around a database table.

### Step 2: Add a Logger (step_2)

We use Go kit's built in logger here, but we could have used logrus or any other logger.

### Step 3: Add Endpoints (step_3)

This is where things get a little interesting, Go kit-wise. An endpoint is a Go kit function that takes a request and returns a response, callable from a Go kit transport.

### Step 4: Add an HTTP Transport (step_4)

And now we have a ReST service around our database table, powered by Go kit.

### Step 5: Add a Protocol Buffer Definition (step_5)

Nothing new integrated into our application in this step. We define our protocol buffers.

### Step 6: Add a gRPC Transport (step_6)

Now we can access our data over gRPC.

## Building

The make file has various targets. To quickly build, run:

```sh
$ make quick
```

To run all the linters and tests, run:

```sh
$ make
```

To run a test coverage report, run:

```sh
$ make coverage
```

## Presentation

The accompanying presentation can be found at `presentation/microservices_with_gokit.md` and is designed to be viewed with [Deckset](https://www.deckset.com).

## Acknowledgments

* Thanks to Michael Dimmitt <https://github.com/michaeldimmitt> for troubleshooting the installation steps
* Go kit <https://gokit.io/> 
* sqlx <http://jmoiron.github.io/sqlx/>
* Testify <https://github.com/stretchr/testify>
* go-sqlmock <https://github.com/DATA-DOG/go-sqlmock>
* gRPC <https://grpc.io/>

Apologies for any I've missed.

## License

Copyright &copy; 2018 Rob Warner

Licensed under the [MIT License](https://hoop33.mit-license.org/)

