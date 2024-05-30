FROM nginx:latest

ARG COMPONENT

COPY $COMPONENT/web/ /usr/share/nginx/html

RUN apt update && apt install -y curl

RUN rm -v /etc/nginx/nginx.conf
ADD $COMPONENT/nginx.conf /etc/nginx/

CMD ["nginx", "-g", "daemon off;"]