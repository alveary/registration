# [SwarmRegistration](https://swarm-registration.herokuapp.com/)

![Travis-CI](https://travis-ci.org/alveary/swarm-registration.svg) by [Travis-CI](https://travis-ci.org/alveary/swarm-registration)


## run test suite

The registration is using Godep as dependency manager,
so just run `godep go test ./...` to run all included tests.

Install `reflex` when you want to run the tests all the time and reload on file changes

```sh
go get github.com/cespare/reflex
```

With reflex enabled just run the `test-all` file.

```sh
./test-all
```
