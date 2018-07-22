# modoki

![Docker Build Status](https://img.shields.io/docker/build/tsuzu/modoki.svg?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/cs3238-tsuzu/modoki)](https://goreportcard.com/report/github.com/cs3238-tsuzu/modoki)
- PaaS modoki(like PaaS)

# Installation
- Install [Docker](https://docker.com)
- Install [docker-compose](https://docs.docker.com/compose/)
- $ git clone https://github.com/cs3238-tsuzu/modoki.git
- $ cd production

- $ cd auth
- $ sh ./gen.sh
- Create **./authconfig.json** in the format  of [authconfig.template.json](https://github.com/cs3238-tsuzu/modoki/blob/master/production/auth/authconfig.template.json)
- $ cd ../

- Create **./.env** file in the format of [.env.template](https://github.com/cs3238-tsuzu/modoki/blob/master/.env.template)
- $ docker-compose up

# License
- Under the MIT License
- Copyright (c) 2018 Tsuzu