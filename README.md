# Downloads rumble video and transcribes it

## To build

    docker build -t rumble-transcriber .

## To run

    docker run --rm rumble-transcriber "https://rumble.com/v2nehs7-discover-the-winning-mindset-of-michael-jordan-10-quotes-to-inspire-you-to-.html"


# To have speakers identified

    # On host machine (your laptop)
    export HF_TOKEN=<access token from hugging face>

    docker run --rm rumble-transcriber "https://rumble.com/v5cweph-kamala-harris-finally-gives-softball-interview-and-its-still-a-total-disast.html"


## To access the docker container while its running
    docker run -it --rm -e HF_TOKEN="$HF_TOKEN" --entrypoint /bin/bash rumble-transcriber
