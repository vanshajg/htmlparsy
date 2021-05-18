The code is live on [Heroku](https://limitless-journey-14259.herokuapp.com/)

# Local Development
Either run it via docker by 
```
docker build -t htmlparsy .
docker run -d -p 8080:8080 htmlparsy
```
or via local install shown below.
## Setting up server
Make sure you have Go version > 1.11 as this repository uses Go modules.
Install dependencies for server by executing the following command in the server folder
```
go mod download
```
Start the server by running the following command in src folder:
```
go run .
```
This runs all files in the current folder bar the test files (*_test.go)

The server should get started on port 8080

To run test cases run the following command in src folder:
```
go test
```
## Setting up Client
Make sure you have a package manager installed. The following example is using yarn.

Install packages:
```
yarn install 
```

Start the local dev server:
```
yarn start
```
The frontend should be available on port 3000
