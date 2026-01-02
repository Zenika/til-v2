# Today I Learned (TiL) - v2 (frontend)

TIL is a simple software that aim to enable article and reflexion-sharing across all Zenika agencies in the world. It
acts as an aggregator (as
daily.dev can do), but on-premises and with filters on content to avoid being spammed with news that don't interest us.

## Getting started

### Developer tools

If you want to develop TiL, the best way is to install NodeJS 18> on your computer, then run `npm install ; npm run dev`. The client will start on
port 5173.

### Docker

If you simply want to run the application, Docker is the best way to do. Simply run `docker build -t til-front:latest .`, and
run it! Full configuration reference below.

## Configuration

All the configuration is done through environment variables. There is only one variable to set:

* **(Mandatory)** `PUBLIC_TIL_SERVER_URL`: The base URL (**without the leading /**) to the backend. In development, you can set it in
  `.env.development` file. The `.env` file must stay at `/api`, because we are using a Nginx reverse proxy on production. You have to specify the
  port, otherwise nginx will mess up. A correct example is: `https://til-backend.35.190.222.151.sslip.io:443`.

## Features

* CRUD to TIL API
* SSE (Server-sent events) to the frontend to refresh the post list automatically
* Settings saved in local storage