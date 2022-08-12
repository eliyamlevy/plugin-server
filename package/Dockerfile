FROM registry.suse.com/bci/bci-micro:latest
RUN  mkdir -p /home/plugin-server/files
COPY bin/plugin-server /usr/bin/
# COPY files             /home/plugin-server/files
WORKDIR /home/plugin-server/
ENTRYPOINT ["plugin-server"]
