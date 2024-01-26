## URL Shortener

A simple URL redirector. I use it as a self-hosted shortening service. But it can be used to handle pretty much any redirection scenario. 

Redirects are set up at startup and read from a `redirects` file. This file is a space-separated list where each line consists of a "short" path and a corresponding URL.

Example:

```
/g https://google.com
/33cs https://news.ycombinator.com/
```

With the above configuration, a request to `https://yourdomain/g` will redirect to Google, and `https://yourdomain/33cs` will redirect to Hacker News.

For any unknown paths, the service redirects to a random Wikipedia article.

## Running

### Build with redirects file

This is the simplest method:

1. Clone the repository.
2. Add a `redirects` file populated with your paths and URLs.
3. Build using Docker: `docker build . -t redirect:latest`

This produces an image that can be run anywhere, with predefined URLs.

Run the container:

`docker run -it --rm -p 8080:8080 redirect:latest`

This method is used in a CI/CD pipeline where a new commit triggers a build of a new container version, which is then pushed to a private registry. Another job handles pulling the new image and restarting the container.

#### Mount a Redirects File

Alternatively, mount a `redirects` file stored on your local filesystem:

`docker run -it --rm -p 8080:8080 -v /path/redirects:/root/redirects redirect:latest`


### Handling New Redirects

Currently, you need to restart the container to add a new redirect. Signal handling for dynamic updates is a planned feature but not yet implemented.

### Adding a Redirect

To add a new redirect:

1. Manually edit the `redirects` file.
2. Alternatively, use the provided script for assistance:

`go build src/cli/add_redirect.go`

This builds an executable. Run it with `./add_redirect --help` for usage instructions.