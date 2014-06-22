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

echo "localhost:27017" > hosts.hosts
time go run mbf.go

# WE DID IT!
# Password is admin:1337h4x0r

# real	0m4.187s
# user	0m0.977s
# sys	0m0.574s

