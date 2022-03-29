# wallester
Wallester Go task

# Start test task
- start docker-compose file or use exists postgre database
- start wallester with the following ENV parameters
  - HTTP_PORT - http server port
  - DB_HOST - database host, default localhost
  - DB_PORT - database port, default 5432
  - DB_USER - database user, default 'user'
  - DB_PASSWORD - database user's password, default 'password'
  - DB_NAME - database name, default 'test'
- open in browser localhost:${HTTP_PORT}
