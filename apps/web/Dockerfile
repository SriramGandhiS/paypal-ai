# Multi-stage Dockerfile for PayBridge AI React Web Frontend
# Syntax: docker build -t paybridge-web -f apps/web/Dockerfile .

# Stage 1: Build static assets
FROM node:20-alpine AS builder

WORKDIR /app

COPY apps/web/package.json ./
RUN npm install

COPY apps/web/ ./
RUN npm run build

# Stage 2: Serve static bundle via Nginx
FROM nginx:alpine

COPY apps/web/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=15s --timeout=3s --retries=3 \
  CMD wget -qO- http://localhost:80/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
