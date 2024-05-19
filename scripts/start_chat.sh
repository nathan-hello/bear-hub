#!/bin/bash

random_hello() {
    # Generate a random number between 1 and 10
    num=$(shuf -i 40-80 -n 1)

    # Print "hello" repeated the random number of times
    printf 'hello %.0s' $(seq $num)
    echo   # Add a newline at the end
}
                        
main() {
        local message=$(random_hello)
        curl -X POST \
          -H "Content-Type: application/json" \
          -d '{"msg-text": "'"${message}"'"}' \
          http://localhost:3000/api/v1/chat/message             
}


while true; do
        main
        sleep 0.5
done
