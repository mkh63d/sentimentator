FROM node:24-alpine  AS builder

WORKDIR /app

COPY package.json package-lock.json ./

RUN --mount=type=cache,target=/root/.npm npm ci

COPY . .

RUN npm run build

FROM nginxinc/nginx-unprivileged:alpine3.22 AS runner

# Install gettext for envsubst
USER root
RUN apk add --no-cache gettext

COPY nginx.conf /etc/nginx/nginx.conf

COPY --chown=nginx:nginx --from=builder /app/dist /usr/share/nginx/html

# Make html directory writable for runtime config generation
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html

COPY start.sh /start.sh
RUN chmod +x /start.sh && chown nginx:nginx /start.sh

USER nginx

EXPOSE 8000

# Use start script to generate runtime config and start nginx
ENTRYPOINT ["/start.sh"]
