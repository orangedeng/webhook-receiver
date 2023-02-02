ARG ARCH
FROM ${ARCH}registry.suse.com/bci/bci-micro:15.4
ARG user=webhookreceiver
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ENV TARGETPATH=${TARGETPLATFORM:-"."}

RUN echo "$user:x:1000:1000::/home/$user:/bin/bash" >> /etc/passwd && \
    echo "$user:x:1000:" >> /etc/group && \
    mkdir /home/$user && \
    chown -R $user:$user /home/$user

COPY bin/$TARGETPATH/webhook-receiver /webhook-receiver
USER $user
CMD ["/webhook-receiver"]
