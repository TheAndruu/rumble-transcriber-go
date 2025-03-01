# Downloads rumble video and transcribes it

## To build

    # Before transcriber was added, it built with:
    docker build -t rumble-transcriber .

    # Now that it has transcription, use docker secrets to 
    # keep the env var out of the built product
    HF_TOKEN="$HF_TOKEN" docker build --secret id=hf_token,env=HF_TOKEN -t rumble-transcriber .

## To run

    docker run --rm rumble-transcriber "https://rumble.com/v2nehs7-discover-the-winning-mindset-of-michael-jordan-10-quotes-to-inspire-you-to-.html"



## To access the docker container while its running
    docker run -it --rm -e HF_TOKEN="$HF_TOKEN" --entrypoint /bin/bash rumble-transcriber
