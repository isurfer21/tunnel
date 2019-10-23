#!/usr/bin/bash env
echo "Tunnel's automated test-case execution"

echo "\nExecute command using cross-domain POST request"
curl --request POST \
  --url 'http://127.0.0.1:9999/terminal' \
  --header 'authorization: Basic YWRtaW46MTIzNDU2' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data cmd=pwd

echo "\n\nExecute command using cross-domain GET request with jsonp calback"
curl --request GET \
  --url 'http://127.0.0.1:9999/terminal?callback=jsonp&cmd=pwd' \
  --header 'authorization: Basic YWRtaW46MTIzNDU2' \
  --data cmd=pwd

echo "\n\nAuthenticate user using cross-domain POST request"
curl --request POST \
  --url 'http://127.0.0.1:9999/authenticate' \
  --header 'authorization: Basic YWRtaW46MTIzNDU2' \
  --header 'content-type: application/x-www-form-urlencoded'

echo "\n\nAuthenticate user using cross-domain GET request with jsonp calback"
curl --request GET \
  --url 'http://127.0.0.1:9999/authenticate?callback=jsonp' \
  --header 'authorization: Basic YWRtaW46MTIzNDU2'

echo "\n\nServer handshake using cross-domain GET request"
curl --request GET \
  --url 'http://127.0.0.1:9999/'

echo "\n\nServer handshake using cross-domain GET request with jsonp calback"
curl --request GET \
  --url 'http://127.0.0.1:9999/?callback=jsonp'

echo "\n\nDone!"