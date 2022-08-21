set -e

go build
docker build -f Dockerfile -t local/chromedp .
rm poc
docker run --rm -it local/chromedp