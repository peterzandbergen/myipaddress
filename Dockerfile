FROM golang:latest AS builder
RUN mkdir /work 
ADD . /work/ 
WORKDIR /work 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch
COPY --from=builder /work/main /main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENV port=8080
ENV NAME=myipaddress_1
CMD ["/main"]
