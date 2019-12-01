# enlabs-test
Enlabs test task solution

## build
``` bash
$ docker-compose up --build

# service available on localhost:8001
# in order to change it, please fix ./config/dev.toml

# the DB already contains a test user
# so you can try to send a transaction request as:

# POST localhost:8001:/user/1/account
# {"state": "win", "amount": "10", "transactionId": "3"}

```