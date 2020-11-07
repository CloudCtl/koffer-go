FROM registry.access.redhat.com/ubi8/ubi-minimal
COPY ./bin/koffer /root/
WORKDIR /root
RUN ./koffer -h
CMD ["bash"]
