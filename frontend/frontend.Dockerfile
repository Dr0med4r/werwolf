# frontend Dockerfile
FROM node:latest AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install
COPY . ./
RUN npx vite build

FROM nginx:1.27
COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=build /app/dist /usr/share/nginx/html
