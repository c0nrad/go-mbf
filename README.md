MongoDB Brute Forcer
===================

Multithreaded mongodb authentication brute forcer written in golang. Supports single mongod brute forcing, or replica set brute forcing. 

```bash
Usage of ./go-mbf:
  -database="admin": name of database to use
  -hostname="hosts.hosts": file containing hostnames
  -passfile="pass.pass": location of password file
  -threads=4: number of db connections to use per machine
  -username="admin": username to bruteforce
  -verbose=false: display each attempt
```

```bash
./go-mbf
-------- MongoDB BruteForcer -------
Loaded password list! Total words: 82830
WE DID IT! Password is admin:1337h4x0r
```

To test it out, try [single.sh](single.sh) or [replica.sh](replica.sh)

[c0nrad](mailto:poptarts4liffe@gmail.com)
