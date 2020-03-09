FROM 041549501818.dkr.ecr.cn-northwest-1.amazonaws.com.cn/csi-s3:1.1.1-dev
LABEL description="csi-s3 testing image"

RUN apt-get update && \
  apt-get install -y \
  git wget make gcc && \
  rm -rf /var/lib/apt/lists/*

RUN wget -q https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz && \
  tar -xf go1.12.5.linux-amd64.tar.gz && \
  rm go1.12.5.linux-amd64.tar.gz && \
  mv go /usr/local

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

WORKDIR /app

# prewarm go mod cache
COPY go.mod .
COPY go.sum .
RUN go mod download

ADD test/test.sh /usr/local/bin

ENTRYPOINT ["/usr/local/bin/test.sh"]
