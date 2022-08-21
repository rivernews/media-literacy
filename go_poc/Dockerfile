FROM chromedp/headless-shell:latest

RUN mkdir -p /tmp/build

ADD poc /tmp/build/poc

ENTRYPOINT ["/tmp/build/poc"]