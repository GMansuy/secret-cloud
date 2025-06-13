FROM node:18-alpine AS builder
ENV NEXT_PUBLIC_BACKEND_URL="rocket-app:443"
WORKDIR /app
COPY . .
RUN npm install && npm run build


FROM nginx:alpine
# Remove default nginx static assets
RUN rm -rf /usr/share/nginx/html/*
# Copy static assets from builder stage
COPY --from=builder /app/out /usr/share/nginx/html

# Configure Nginx to serve the static files
RUN echo 'server { \
    listen 80; \
    server_name _; \
    root /usr/share/nginx/html; \
    index index.html; \
    location = /callback { \
            try_files /callback.html =404;\
        } \
    # Add detailed error logging \
    error_log /var/log/nginx/error.log debug; \
    access_log /var/log/nginx/access.log; \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]