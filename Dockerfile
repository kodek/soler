# iron/go is the alpine image with only ca-certificates added
#FROM iron/go
FROM golang:1.10.3
WORKDIR /app
# Now just add the binary
COPY server_main .

ARG travis_commit
ENV TRAVIS_COMMIT=$travis_commit

ENV TZ America/Los_Angeles

ENV ENABLE_SOLAREDGE_POLLING=true
ENV ENABLE_SENSE_POLLING=true

RUN echo $TZ > /etc/timezone && \
    apt-get update && apt-get install -y tzdata && \
    rm /etc/localtime && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    apt-get clean

ENTRYPOINT [
  "./server_main",
  "--enable_solaredge_polling=$ENABLE_SOLAREDGE_POLLING",
  "--enable_sense_polling=$ENABLE_SENSE_POLLING"]
