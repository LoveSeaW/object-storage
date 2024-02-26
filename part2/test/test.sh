#! /bin/bash
curl -v 127.0.0.1:8080/objects/test2 -XPUT -d "this is objects test2"

curl 192.168.110.134:8800/locate/test2
echo
curl 192.168.110.134:8800/objects/test2
echo