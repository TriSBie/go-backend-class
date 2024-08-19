# Build stage

# ?? Why go build ? 
# Go build is a command that compiles the Go source code files into executable binary file.

FROM golang:1.23.0-alpine3.20 AS builder
# set the workdir directory
WORKDIR /app
COPY . .
# run output as binary file as main 
RUN go build -o main main.go


# Run stage
FROM alpine:3.20
WORKDIR /app
# copy the binary compiled file from the build stage to the run stage
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
# define a command to run the app
ENTRYPOINT ["./main"]