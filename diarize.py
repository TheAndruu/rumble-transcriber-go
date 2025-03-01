#!/usr/bin/env python3
import sys
from pyannote.audio import Pipeline
import torch

def diarize(audio_file):
    # Load pipeline from the default cache (pre-downloaded)
    pipeline = Pipeline.from_pretrained("pyannote/speaker-diarization-3.1")
    pipeline.to(torch.device("cpu"))
    diarization = pipeline(audio_file)
    for turn, _, speaker in diarization.itertracks(yield_label=True):
        # Clean up speaker label from "SPEAKER_00" to "0"
        speaker_id = speaker.replace("SPEAKER_", "")
        print(f"{turn.start:.2f} {turn.end:.2f} Speaker_{speaker_id}", flush=True)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 diarize.py <audio_file>", flush=True)
        sys.exit(1)
    audio_file = sys.argv[1]
    diarize(audio_file)