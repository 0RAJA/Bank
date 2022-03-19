#!/bin/sh

set -e #确保脚本在指令返回非零直接返回

echo "run db migrate"
/app/migrate -path /app/migration -database "postgresql://root:secret@postgres12:5432/bank?sslmode=disable" -verbose up

echo "start the app"
exec "$@" # 执行传递给脚本的所有参数

