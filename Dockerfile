# Build stage
FROM golang:1.21 AS builder

WORKDIR /app

# Copy Go source files
COPY main.go .

# Build the Go binary
RUN go mod init rumble-transcriber && \
    go build -o transcriber main.go

# Final stage
FROM ubuntu:22.04

WORKDIR /app

# Install dependencies: ffmpeg, yt-dlp, git, build-essential, curl, and cmake
RUN apt-get update && apt-get install -y \
    ffmpeg \
    python3 \
    python3-pip \
    git \
    build-essential \
    curl \
    cmake && \
    pip3 install yt-dlp

# Install whisper.cpp and move the correct binary
RUN git clone https://github.com/ggerganov/whisper.cpp.git && \
    cd whisper.cpp && \
    cmake -B build && \
    cmake --build build --config Release && \
    mv build/bin/whisper-cli /app/whisper && \
    cd /app

    #sh ./models/download-ggml-model.sh base.en && \

# Download a Whisper model (base English model)
RUN curl -L -o ggml-base.en.bin https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.en.bin

# Copy the Go binary from the builder stage
COPY --from=builder /app/transcriber /app/transcriber

# Set executable permissions
RUN chmod +x /app/transcriber /app/whisper

# Command to run the transcriber
ENTRYPOINT ["/app/transcriber"]