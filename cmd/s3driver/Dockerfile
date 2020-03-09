FROM debian:stretch
LABEL description="csi-s3 slim image"

RUN apt-get update && \
  apt-get install -y curl fuse && \
  rm -rf /var/lib/apt/lists/*

COPY ./_output/s3driver /s3driver
ENTRYPOINT ["/s3driver"]
