FROM gcr.io/distroless/static:latest
LABEL maintainers="MayaData Authors"
LABEL description="Storage Provisioner"

COPY ./d-provisioner d-provisioner
ENTRYPOINT ["/d-provisioner"]
