#!/usr/bin/env bash
GO_VERSION=1.9.2

docker build -f testdata/Dockerfile.travis -t local/travis . && \
docker run --user root -dit --rm --name travis-debug local/travis:latest /sbin/init
docker ps -a
docker exec --user travis travis-debug bash -l -c "/home/travis/.travis/travis-build/bin/travis compile --no-interactive" && \
docker exec --user travis -it travis-debug bash
docker rm -f travis-debug

#    -e TRAVIS_GO_VERSION=${GO_VERSION} \
#    -v $(pwd)/../../:/home/travis/gopath/src/github.com/peterbueschel/awsdefault/ \
#    --rm \
#    --name travis-debug \
#    -dit travisci/ci-garnet:packer-1512502276-986baf0 /sbin/init
