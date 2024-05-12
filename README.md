# ImageStack

Service to compress images on the fly

## Development

### Usage

`docker compose up`

1. RabbitMQ dashboard: http://localhost:15672
2. App endpoints:

- Check status of any image: http://localhost:8080/status?url={your-image-url}
- To process a new image: http://localhost:8080/submit?url={your-image-url}
