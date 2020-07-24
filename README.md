# Bullshit generator

The bullshit generator can generate articles with a given topic and a word count requirement.   

There is also a live site at [https://howtobullshit.me](https://howtobullshit.me).

## Getting started
1. Install [docker](https://www.docker.com).
2. Clone the repository and get into it.
3. Build the docker image.
    ```
    docker build -t bullshit .
    ```
4. Start the docker container.
    ```
    docker run -d -p 1234:10000 bullshit 
    ```
    This runs the api server on port `1234`
5. Now you can test it with `curl`.
    ```
    curl http://127.0.0.1:1234/bullshit \                                                                                                                                           16:16:07 
    -X POST \
    -H 'Content-Type: application/json; charset=utf-8' \
    -d '{"Topic": "hi", "MinLen": 100}' 
    ```
## API

The generator API is open for public. For more details on how to use it, please reference the [code](https://github.com/StillFantastic/bullshit/blob/5cf1a7fc9c70442a213cd8941450695dd13fa76c/index.html#L100a)  
Note that there is a request rate limit of 1000 reqeusts within 1000 seconds for each IP.