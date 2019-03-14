# Build a Golang app with the Gin framework, and authenticate with JWT

Based on the [Build a Golang app with the Gin framework, and authenticate with Auth0 + JWT ](https://hakaselogs.me/2018-04-20/building-a-web-app-with-go-gin-and-react) article, Implemented the example of Go back-end and React front-end. __It's THE JUST PROTOTYPE OF JWT Token management in go and react.__

## Setup

1. Run `mv .env.sample .env` and update with valid credentials.

``` bash
export USER_ID="admin"
export USER_PASSWORD="14123123"
```

2. Source the environment variables - `source .env`
3. Pull the dependencies.

``` bash
$ go mod init
$ go get 
```
4. Run the project
```
$ go run main.go
```

## Author
* [Archer Yoo](https://jinhoyoo.github.io/)
* [Francis Sunday](https://twitter.com/codehakase) 

## License

This project is licensed under the MIT license. See the [LICENSE](LICENSE) file for more info.

## Thanks

My colleagues, [bkmoon318](https://github.com/bkmoon318) and [DooSikChoi](https://github.com/DooSikChoi). 