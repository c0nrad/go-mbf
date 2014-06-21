MongoDB Brute Forcer
===================

Multithreaded mongodb authentication brute forcer written in golang.

```bash
Usage of ./go-mbf:
  -database="admin": name of database to use
  -hostname="127.0.01": hostname containing MongoDB
  -passfile="pass.pass": location of password file
  -username="abc": username to bruteforce
  -verbose=false: display each attempt
```

```bash
$ ./go-mbf
Loaded password list! Total words: 14344392

-----WE DID IT abc 123-----
Number of tries: 3984
```

To test it out, try single.sh.

[c0nrad](poptarts4liffe@gmail.com)