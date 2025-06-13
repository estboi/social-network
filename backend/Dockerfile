FROM golang:1.20 as server

WORKDIR /server

COPY . .

RUN apt update && apt install gcc musl-dev

RUN go build -o sn-backend .

# FROM debian:bullseye-slim

# WORKDIR /server


# COPY --from=server /server/sn-backend .

# COPY --from=server /server/database/migrations ./database/migrations

# COPY --from=server /server/database/database.db ./database/database.db

EXPOSE 8080

CMD [ "./sn-backend" ]