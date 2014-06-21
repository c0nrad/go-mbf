# Kill Previous Instances of Mongod
pkill mongod

rm -rf ./data/
mkdir -p ./data/db

rm single.log*

mongod --auth --fork --logpath "single.log" --dbpath ./data/db

sleep 2

echo """
db = db.getSiblingDB('admin')
db.createUser({user: 'admin', pwd: '1337h4x0r', roles:['dbAdmin']})
""" > addUser.js
mongo addUser.js

sleep 2

go run mbf.go