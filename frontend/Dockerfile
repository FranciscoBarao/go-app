FROM node:26-alpine AS builder

WORKDIR /src

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# Vite inlines VITE_* values into the bundle, so the URL is fixed at image build time.
ARG VITE_CATALOG_BASE_URL=http://localhost:8081
ENV VITE_CATALOG_BASE_URL=${VITE_CATALOG_BASE_URL}
RUN npm run build

FROM nginx:alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /src/dist /usr/share/nginx/html

EXPOSE 80
