FROM nginx:alpine

ADD frontend_static/application.conf /etc/nginx/conf.d/default.conf
ADD frontend_static/prepare /bin/prepare
RUN chmod +x /bin/prepare

ADD common/container-discovery-wrapper /bin/container-discovery-wrapper
RUN chmod +x /bin/container-discovery-wrapper

CMD /bin/container-discovery-wrapper \
  --command="/bin/start" \
  --service-name="frontend-static" \
  --service-exposed-port=80 \
  --service-dependency="php-fpm" \
  --discovery-server="$DISCOVERY_SERVER" \
  --ak-instance-identifier="$AK_INSTANCE_IDENTIFIER"
