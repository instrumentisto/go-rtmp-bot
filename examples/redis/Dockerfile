FROM golang:onbuild
MAINTAINER FlexConstructor <flexconstructor@gmail.com>
RUN mkdir /opt/video                                  \
    && chmod -644 /opt/video                          \
    && wget -O /opt/video/sample.flv \
    http://techslides.com/demos/sample-videos/small.flv
COPY main.go /go/src
CMD ["go-wrapper", "run","-flv_file","/opt/video/sample.flv","-redis","localhost:6379"]
