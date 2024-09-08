# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder

COPY . .

RUN go mod download

# creating a static binary on the root of filesystem of the image
# here the static binary is named sfss
RUN CGO_ENABLED=0 GOOS=linux go build -o /sfss

EXPOSE 3000

# runnning the static binary
CMD ["/sfss"]

# this is a distroless build stage whihc is focused on security and performance
# only contains essential components required to run specific application
FROM gcr.io/distroless/base-debian11 AS prod_build

WORKDIR /

COPY --from=builder /sfss /sfss

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/sfss"]

