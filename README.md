# Definition Mapper

master:  [![CircleCI](https://circleci.com/gh/ernestio/definition-mapper/tree/master.svg?style=shield)](https://circleci.com/gh/ernestio/definition-mapper/tree/master)  
develop: [![CircleCI](https://circleci.com/gh/ernestio/definition-mapper/tree/develop.svg?style=shield)](https://circleci.com/gh/ernestio/definition-mapper/tree/develop)

This service works as a router to redirect all definition mapping requests to its specific definition mapper (aws or vcloud)

## Build status

* Master: [![CircleCI](https://circleci.com/gh/ernestio/definition-mapper/tree/master.svg?style=svg)](https://circleci.com/gh/ernestio/definition-mapper/tree/master)
* Develop: [![CircleCI](https://circleci.com/gh/ernestio/definition-mapper/tree/develop.svg?style=svg)](https://circleci.com/gh/ernestio/definition-mapper/tree/develop)

## Install

This microservice uses make shortcuts to manage all dependencies, install it just running:
```
make deps
make install
```


## Running Tests

```
make test
```

## Contributing

Please read through our
[contributing guidelines](CONTRIBUTING.md).
Included are directions for opening issues, coding standards, and notes on
development.

Moreover, if your pull request contains patches or features, you must include
relevant unit tests.

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2015 ernest.io authors.

Code released under
[the Mozilla Public License Version 2.0](LICENSE).
