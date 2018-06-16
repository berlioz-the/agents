FROM node:carbon

# RUN apt-get update && apt-get install -y vim

ENV NODE_PATH=/usr/local/lib/node_modules

COPY ./ /src/

WORKDIR /src
RUN npm install
RUN npm install --only=dev

EXPOSE 55555
CMD ["node", "."]
