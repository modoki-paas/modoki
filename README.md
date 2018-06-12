# modoki

![Docker Build Status](https://img.shields.io/docker/build/tsuzu/modoki.svg?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/cs3238-tsuzu/modoki)](https://goreportcard.com/report/github.com/cs3238-tsuzu/modoki)
- PaaS modoki(like PaaS)

# Installation
- Install [Docker](https://docker.com)
- Install [docker-compose](https://docs.docker.com/compose/)
- $ mkdir modoki && cd modoki

- $ mkdir cred && cd cred
- $ wget https://raw.githubusercontent.com/cs3238-tsuzu/modoki/master/cred/gen.sh
- $ sh ./gen.sh
- $ cd ../

- $ wget https://github.com/cs3238-tsuzu/modoki/blob/master/production/docker-compose.yml
- Create .env file in the format of [.env.template](https://github.com/cs3238-tsuzu/modoki/blob/master/.env.template)
- $ docker-compose up

# License
- Under the MIT License
- Copyright (c) 2018 Tsuzu