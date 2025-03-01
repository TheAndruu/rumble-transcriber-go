#!/usr/bin/env python3
import sys
from pyannote.audio import Pipeline
import torch

def diarize(audio_file):
    # Load pipeline from default cache location
    pipeline = Pipeline.from_pretrained("pyannote/speaker-diarization-3.1")
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