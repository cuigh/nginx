# ---- Build----
FROM golang:1.17-alpine AS builder
WORKDIR /app
COPY . .
RUN apk update && apk add git
RUN go build -o starter

# ---- Release ----
FROM nginx:1.21.6-alpine
LABEL maintainer="cuigh <noname@live.com>"
WORKDIR /usr/share/nginx/html/
COPY --from=builder /app/starter /usr/local/bin/
COPY default.conf /etc/nginx/conf.d/
CMD ["starter"]