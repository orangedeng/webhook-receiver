ARG ARCH
FROM ${ARCH}alpine:3.14.2
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ENV TARGETPATH=${TARGETPLATFORM:-"."}
COPY bin/$TARGETPATH/webhook-receiver /webhook-receiver
CMD ["/webhook-receiver"]
