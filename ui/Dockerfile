#build stage
FROM node:alpine as ui-builder
WORKDIR /usr/local/app
COPY . .
RUN npm install
RUN npm run build

#final
FROM nginx:alpine
COPY --from=ui-builder /usr/local/app/dist/matspinner /usr/share/nginx/html
LABEL Name=ui Version=0.0.1
EXPOSE 80
