FROM golang:1.22.1 as build
WORKDIR /src
COPY . /src/
RUN  GOOS=linux GOARCH=amd64 go build -o /bin/app main.go

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=build /bin/app /bin/app
EXPOSE 8080
CMD ["/bin/app"]