#!/bin/sh

set -e

## 这个是在容器中运行的;
until nc -w2 mysql 3306 && nc -w2 redis 6379
do
  echo "mysql or redis is unavailable - sleeping" >&2
  sleep 1
done
## 后面是真正要运行的程序
exec "$@"
