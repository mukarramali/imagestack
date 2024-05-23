# ImageStack

Service to compress images on the fly

## Development

### Architecture

![Architecture](./docs/ImageStack.png)

### Usage

`docker compose up`

1. RabbitMQ dashboard: http://localhost:15672
2. Image stack service: http://localhost:8080/?url={your-image-url}

## Goal

Currently all the queues get handled by the same service which serves endpoint to the end user.
Goal is run different queues in their own service talking to the same rabbit mq.
And then scale services based on the usage.
