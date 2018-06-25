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
$ make
```

6. Run the app, which will create the `players` table

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

## License

Copyright &copy; 2018 Rob Warner

Licensed under the [MIT License](https://hoop33.mit-license.org/)

