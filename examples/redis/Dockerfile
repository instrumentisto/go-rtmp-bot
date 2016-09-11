FROM golang:onbuild
MAINTAINER FlexConstructor <flexconstructor@gmail.com>
RUN mkdir /opt/video                                  \
    && chmod -644 /opt/video                          \
    && wget -O /opt/video/big_buck_bunny_720p_2mb.flv \
    http://www.sample-videos.com/video/flv/720/big_buck_bunny_720p_2mb.flv
COPY main.go /go/src
CMD ["go-wrapper", "run","-flv_file","/opt/video/big_buck_bunny_720p_2mb.flv","-redis","localhost:6379"]
