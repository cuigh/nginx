# ---- Build ----
FROM node:alpine AS builder
WORKDIR /app
COPY . .
RUN yarn install
RUN yarn run build

# ---- Release ----
FROM cuigh/nginx
LABEL maintainer="cuigh <noname@live.com>"
COPY --from=builder /app/dist .
COPY --from=builder /app/src/config config/