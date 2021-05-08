FROM node:16-alpine

RUN apk add --no-cache \
    udev \
    ttf-freefont \
    chromium

ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true

RUN mkdir -p /home/node/app/node_modules && chown -R node:node /home/node/app

WORKDIR /home/node/app

COPY package*.json ./

RUN npm install

COPY --chown=node:node . .

USER node

CMD [ "node", "app.js" ]
