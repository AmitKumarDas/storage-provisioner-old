FROM gcr.io/distroless/static:latest
LABEL maintainers="MayaData Authors"
LABEL description="Storage Provisioner"

COPY ./storage-provisioner storage-provisioner
ENTRYPOINT ["/storage-provisioner"]
