#!/usr/bin/env python3
import sys
import os
from pyannote.audio import Pipeline
import torch
from huggingface_hub import login

def diarize(audio_file):
    # Get token from environment variable
    hf_token = os.getenv("HF_TOKEN")
    if not hf_token:
        raise ValueError("HF_TOKEN environment variable not set. Please set it with your Hugging Face access token.")

    # Log in to Hugging Face
    login(token=hf_token)

    # Load the pipeline
    pipeline = Pipeline.from_pretrained("pyannote/speaker-diarization-3.1",
                                        use_auth_token=hf_token)
    pipeline.to(torch.device("cpu"))
    diarization = pipeline(audio_file)
    for turn, _, speaker in diarization.itertracks(yield_label=True):
        print(f"{turn.start:.2f} {turn.end:.2f} Speaker_{speaker}")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 diarize.py <audio_file>")
        sys.exit(1)
    audio_file = sys.argv[1]
    diarize(audio_file)