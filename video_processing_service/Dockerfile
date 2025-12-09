FROM golang:1.25-alpine AS build
RUN apk add --no-cache curl libstdc++ libgcc python3 py3-pip ffmpeg nodejs npm && \
    pip3 install --break-system-packages yt-dlp && \
    yt-dlp --version

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate && \
    curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64-musl -o tailwindcss && \
    chmod +x tailwindcss && \
    ./tailwindcss -i cmd/web/styles/input.css -o cmd/web/assets/css/output.css

RUN go build -o main cmd/api/main.go

FROM alpine:latest AS prod

RUN apk add --no-cache python3 py3-pip ffmpeg curl gcompat&& \
    pip3 install --break-system-packages yt-dlp

# 2. Install Deno manually
# We download the official Linux binary, unzip it, and place it in the path.
RUN curl -L https://github.com/denoland/deno/releases/latest/download/deno-x86_64-unknown-linux-gnu.zip -o deno.zip && \
    unzip deno.zip && \
    mv deno /usr/bin/deno && \
    rm deno.zip && \
    chmod +x /usr/bin/deno

RUN command -v ffmpeg >/dev/null 2>&1 || (echo "ffmpeg missing" && exit 1) && \
    command -v ffprobe >/dev/null 2>&1 || (echo "ffprobe missing" && exit 1) && \
    command -v deno >/dev/null 2>&1 || (echo "deno missing" && exit 1) && \
    python3 -c "import yt_dlp" >/dev/null 2>&1 || (echo "yt-dlp python package missing" && exit 1) && \
    yt-dlp --version

WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE 8080
CMD ["./main"]
