# Downloads rumble video and transcribes it

## To build

Now that it has transcription, use docker secrets to 
keep the env var out of the built product

In hugging face, accept the licenses at:
- https://huggingface.co/pyannote/segmentation-3.0
- https://huggingface.co/pyannote/speaker-diarization-3.1

Then generate an access token with READ permissions in hugging face's settings

Set the HF access token on your host machine in the env var: HF_TOKEN

    export HF_TOKEN=<hf_access_token>

Can then build the project with docker secrets:

    HF_TOKEN="$HF_TOKEN" docker build --secret id=hf_token,env=HF_TOKEN -t rumble-transcriber .

## To run

    docker run --rm rumble-transcriber "https://rumble.com/v5cweph-kamala-harris-finally-gives-softball-interview-and-its-still-a-total-disast.html"


## To access the docker container while its running
    docker run -it --rm -e HF_TOKEN="$HF_TOKEN" --entrypoint /bin/bash rumble-transcriber
