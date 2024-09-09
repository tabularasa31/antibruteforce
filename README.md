[![CI/CD](https://github.com/tabularasa31/antibruteforce/actions/workflows/main.yml/badge.svg)](https://github.com/tabularasa31/antibruteforce/actions/workflows/main.yml)   [![Linters](https://github.com/tabularasa31/antibruteforce/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/tabularasa31/antibruteforce/actions/workflows/golangci-lint.yml)    [![Go Report Card](https://goreportcard.com/badge/github.com/tabularasa31/antibruteforce)](https://goreportcard.com/report/github.com/tabularasa31/antibruteforce)

# AntibruteForce

AntibruteForce is a microservice designed to prevent password brute-force attacks during user authentication.

## Overview

This service is called before user authentication and can either allow or block the attempt. It's intended for server-to-server use only and is hidden from end-users.

## How It Works

The service limits the frequency of authentication attempts for various parameter combinations specified in the `config.yml` file:

- No more than N = 10 attempts per minute for a given login.
- No more than M = 100 attempts per minute for a given password (reverse brute-force protection).
- No more than K = 1000 attempts per minute for a given IP (high number due to NAT).

The Generic Cell Rate Algorithm (GCRA), also known as the leaky bucket algorithm, is used to count and limit request frequencies. 

### GCRA Algorithm

GCRA was chosen for its efficiency in managing rate limiting with minimal memory usage. In our context, it works by:

1. Assigning a "bucket" to each login/password/IP.
2. Each request "fills" the bucket by a certain amount.
3. The bucket "leaks" at a constant rate.
4. If a request would overflow the bucket, it's considered a brute-force attempt.

This approach allows for occasional bursts of traffic while still enforcing long-term rate limits.


## Configuration

The service configuration is located in the `config.yml` file. The main configuration parameters are loginLimit, passLimit, and ipLimit - the limits at which the service considers an attempt to be a brute-force attack.

## Architecture

The microservice consists of:
- gRPC API
- Redis database for storing buckets
- PostgreSQL database for storing black/white lists
- Command-line interface for interacting with the service

## Deployment

To deploy the microservice:

1. Clone the repository:
   ```
   git clone https://github.com/tabularasa31/antibruteforce.git
   cd antibruteforce
   ```

2. Set up the environment:
   ```
   # Edit config.yml with your settings
   ```

3. Run the service:
   ```
   make up
   ```

## Testing

Run the test suite with:
```
make test
```

## Performance

The service is designed to handle high loads:
- Tested up to 10,000 requests per second on a standard 4-core server.
- Redis and PostgreSQL should be properly tuned for production environments.
- Consider deploying multiple instances behind a load balancer for very high traffic scenarios.

## License

This project is licensed under the MIT License.
