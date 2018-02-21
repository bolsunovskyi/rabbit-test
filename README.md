# Rabbit event tracker test

Run docker image:
- `docker run --rm -p 4369:4369 -p 5671:5671 -p 5672:5672 -p 25672:25672 -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password rabbitmq`