FROM nginx:1.13.3

# forward nginx access and error logs to stdout and stderr of the ingress
# controller process
RUN ln -sf /proc/1/fd/1 /var/log/nginx/access.log \
	&& ln -sf /proc/1/fd/2 /var/log/nginx/error.log

COPY nginx-ingress nginx/templates/nginx.ingress.tmpl nginx/templates/nginx.tmpl /

RUN rm /etc/nginx/conf.d/*

RUN mkdir -p /etc/nginx/secrets
ADD default.pem /etc/nginx/secrets/default

ENTRYPOINT ["/nginx-ingress"]
