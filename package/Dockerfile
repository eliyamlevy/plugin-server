FROM registry.suse.com/bci/bci-micro:latest
RUN  mkdir /home/plugin-server
COPY bin/plugin-server /usr/bin/
COPY files             /home/plugin-server/files
ENTRYPOINT ["plugin-server"]
