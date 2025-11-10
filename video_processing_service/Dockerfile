FROM golang:1.25-alpine AS build
RUN apk add --no-cache curl libstdc++ libgcc python3 py3-pip ffmpeg && \
    pip3 install --break-system-packages yt-dlp

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
RUN apk add --no-cache python3 py3-pip ffmpeg && \
    pip3 install --break-system-packages yt-dlp && \
    command -v ffmpeg >/dev/null 2>&1 || (echo "ffmpeg missing" && exit 1) && \
    command -v ffprobe >/dev/null 2>&1 || (echo "ffprobe missing" && exit 1) && \
    python3 -c "import yt_dlp" >/dev/null 2>&1 || (echo "yt-dlp python package missing" && exit 1)

WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE 8080
CMD ["./main"]
