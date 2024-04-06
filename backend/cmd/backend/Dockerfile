FROM golang:1.22.1 as build
WORKDIR /src
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/app main.go

FROM scratch
COPY --from=build /bin/app /bin/app
EXPOSE 8080
CMD ["/bin/app"]