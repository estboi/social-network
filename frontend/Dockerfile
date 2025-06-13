FROM node:20.8.0 as frontend 

WORKDIR /frontend

COPY . .

RUN npm install && npm run build

FROM nginx:alpine

WORKDIR /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

RUN rm -rf ./*

COPY --from=frontend /frontend/build .

ENTRYPOINT ["nginx", "-g", "daemon off;"]