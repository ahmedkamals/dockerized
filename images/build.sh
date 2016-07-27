#!/bin/bash
prefix="localhost/ak"

# -- Frontend images -- #
docker build -t "$prefix-frontend-static:latest" -f frontend_static/Dockerfile .

# -- Execution images -- #
docker build -t "$prefix-php-fpm:latest" -f execution/php_fpm/Dockerfile .
docker build -t "$prefix-periodical-tasks:latest" -f execution/periodical_tasks/Dockerfile .

# -- Codebase image -- #
docker build -t "$prefix-codebase:latest" -f codebase/Dockerfile .

# -- Data stores -- #
docker build -t "$prefix-data-storage-relational:latest" -f data_storage/relational/Dockerfile .
# docker build -t "$prefix-data-storage-cache:latest" -f data_storage/cache/Dockerfile .
# docker build -t "$prefix-data-storage-key-value:latest" -f data_storage/key_value/Dockerfile .

# docker build -t "$prefix-datastore-session:latest" -f datastore_session/Dockerfile .

#docker build -t "$prefix-bootstrapper-sandbox:latest" -f bootstrapper_sandbox/Dockerfile .
# docker build -t "$prefix-bootstrapper-development:latest" -f bootstrapper_development/Dockerfile .


#docker build -t "$prefix-workers:latest" -f workers/Dockerfile .

# -- Supporting nodes -- #
#docker build -t "$prefix-workers-server:latest" -f workers_server/Dockerfile .
