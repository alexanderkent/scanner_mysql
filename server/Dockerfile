FROM ubuntu:18.04
RUN apt-get update \
    && apt-get install -y --no-install-recommends netcat \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /app

WORKDIR /app 

ADD server.sh /app/server.sh

RUN chmod +x /app/server.sh 

ENTRYPOINT ["/app/server.sh"]
