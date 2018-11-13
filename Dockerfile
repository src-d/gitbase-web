FROM debian:stretch-slim

RUN apt-get update && \
  apt-get install -y --no-install-recommends --no-install-suggests \
  ca-certificates \
  && apt-get clean

ADD ./build/bin /bin

ENTRYPOINT ["/bin/gitbase-web"]
CMD ["serve"]
