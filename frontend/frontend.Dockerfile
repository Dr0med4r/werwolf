FROM nginx:1.27
COPY nginx.conf /etc/nginx/nginx.conf
COPY . /usr/share/nginx/html
