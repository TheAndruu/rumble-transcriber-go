# Build stage
FROM golang:1.21 AS builder

WORKDIR /app

COPY main.go .
RUN go mod init rumble-transcriber && \
    go build -o transcriber main.go

# Final stage
FROM ubuntu:22.04

WORKDIR /app

# Install dependencies: ffmpeg, yt-dlp, git, build-essential, curl, cmake, Python, and pip
RUN apt-get update && apt-get install -y \
    ffmpeg \
    python3 \
    python3-pip \
    git \
    build-essential \
    curl \
    cmake && \
    pip3 install yt-dlp

# Install whisper.cpp
RUN git clone https://github.com/ggerganov/whisper.cpp.git && \
    cd whisper.cpp && \
    cmake -B build && \
    cmake --build build --config Release && \
    mv build/bin/whisper-cli /app/whisper && \
    cd /app

# Install Pyannote.audio, PyTorch, and huggingface_hub
RUN pip3 install pyannote.audio torch huggingface_hub

# Copy the diarization script
COPY diarize.py /app/diarize.py

# Download a Whisper model (base English model)
RUN curl -L -o ggml-base.en.bin https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.en.bin

# Copy the Go binary
COPY --from=builder /app/transcriber /app/transcriber

# Set executable permissions
RUN chmod +x /app/transcriber /app/whisper

# Command to run the transcriber
ENTRYPOINT ["/app/transcriber"]