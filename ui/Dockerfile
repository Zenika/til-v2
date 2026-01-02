FROM node:18-alpine AS build

WORKDIR /app

COPY . .

RUN npm install

RUN npm run build

FROM nginx:1.27.5-alpine AS run

RUN apk add bind-tools moreutils
COPY default.conf /etc/nginx/conf.d/default.conf
COPY custom-entrypoint.sh /custom-entrypoint.sh
RUN chmod 555 /custom-entrypoint.sh
COPY --from=build /app/build/ /usr/share/nginx/html

ENTRYPOINT ["/custom-entrypoint.sh"]