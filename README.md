# ImageStack

Service to compress images on the fly

Try out here:

https://imagestack-latest.sliplane.app/?quality=99&width=800&url=https://raw.githubusercontent.com/mukarramali/imagestack/main/docs/ImageStack.png

## Development

### Architecture

![Architecture](./docs/ImageStack.png)

### Usage

`docker compose up`

1. RabbitMQ dashboard: http://localhost:15672
2. Image stack service: http://localhost:8080/?url={your-image-url}

## LoadTesting

`go run ./test/load.go [number-of-requests]`

## What's need to be done

1. Introduce CDN in front
