# Build stage

# ?? Why go build ? 
# Go build is a command that compiles the Go source code files into executable binary file.

FROM golang:1.23.0-alpine3.20 AS builder
# set the workdir directory
WORKDIR /app
COPY . .
# run output as binary file as main 
RUN go build -o main main.go
RUN apk add --no-cache curl 
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-386.tar.gz | tar xvz


# Run stage
FROM alpine:3.20
WORKDIR /app
# copy the binary compiled file from the build stage to the run stage
COPY --from=builder /app/main . 
# copy the migrate binary file from the build stage to the run stage
COPY --from=builder /app/migrate ./migrate
COPY app.env /app/app.env
COPY start.sh .
COPY wait-for.sh .

COPY db/migration ./migration

RUN chmod +x /app/start.sh
RUN chmod +x /app/wait-for.sh
RUN chmod 644 /app/app.env

EXPOSE 8080

# When using CMD and ENTRYPOINT together, the CMD command is passed as an argument to the ENTRYPOINT command.
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
