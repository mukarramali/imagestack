# ImageStack

Service to compress images on the fly

## Development

### Architecture

![Architecture](./docs/ImageStack.png)

### Usage

`docker compose up`

1. RabbitMQ dashboard: http://localhost:15672
2. Image stack service: http://localhost:8080/?url={your-image-url}

## What's need to be done

1. Introduce CDN in front
2. Run different queues in their own service
