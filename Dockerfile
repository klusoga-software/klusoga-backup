FROM alpine:3
ENTRYPOINT ["/klusoga-backup-agent"]
COPY klusoga-backup-agent /