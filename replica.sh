# Kill Previous Instances of Mongod
pkill mongod

info() {
	tput setaf 2
	echo $1
	tput sgr0
}


rm -rf ./data/
mkdir -p ./data/rs0
mkdir -p ./data/rs1
mkdir -p ./data/rs2

rm replica1.log*
rm replica2.log*
rm replica3.log*

info 'Spawning mongods'
mongod --replSet set0 --logpath "replica1.log" --dbpath "./data/rs0" --port 13370 --fork --oplogSize 200
mongod --replSet set0 --logpath "replica2.log" --dbpath "./data/rs1" --port 13371 --fork --oplogSize 200
mongod --replSet set0 --logpath "replica3.log" --dbpath "./data/rs2" --port 13372 --fork --oplogSize 200

sleep 5

echo """
config = { '_id': 'set0', members: [
	{ _id: 0, host: 'localhost:13370' },
	{ _id: 1, host: 'localhost:13371' },
	{ _id: 2, host: 'localhost:13372' }
]}
rs.initiate(config)
""" > initRs.js
info 'Initiating replica set!'
mongo --quiet --port 13370 initRs.js

sleep 20

echo """
db = db.getSiblingDB('admin')
db.createUser({user: 'admin', pwd: '1337h4x0r', roles:['dbAdmin']})
""" > addUser.js

info "Adding user admin:1337h4x0r"
mongo --quiet  --port 13370 addUser.js

info "LETS START THE BRUTE FORCING"

echo """localhost:13370
localhost:13371
localhost:13372""" > hosts.hosts

time go run mbf.go

# WE DID IT!
# Password is admin:1337h4x0r

# real	0m3.243s
# user	0m1.008s
# sys	0m0.524s