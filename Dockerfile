FROM        golang
RUN         mkdir -p /app
WORKDIR     /app
COPY        go.mod .
COPY        go.sum .
RUN         go mod download
COPY        . .
RUN         go build -o app
RUN ls
ENTRYPOINT [ "./app" ]