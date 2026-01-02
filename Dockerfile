FROM debian:trixie

COPY --from=til-backend:latest /til /server
COPY --from=til-frontend:latest /usr/share/nginx/html /static
RUN chmod +x /server

ENV TIL_USE_EMBEDDED_FRONTEND=true

ENTRYPOINT ["/server"]