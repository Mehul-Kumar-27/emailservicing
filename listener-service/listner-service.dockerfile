FROM alpine:latest
RUN mkdir /app
COPY listnerServiceApp /app
CMD [ "/app/listnerServiceApp" ]
