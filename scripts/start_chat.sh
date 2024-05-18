#!/bin/bash

main() {
        curl -X POST \
          -H "Content-Type: application/json" \
          -d '{"msg-text": "hello"}' \
          http://localhost:3000/api/v1/chat/message             
}


while true; do
        main
        sleep 2
done
