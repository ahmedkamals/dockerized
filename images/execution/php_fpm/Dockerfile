FROM php:7.0.8-fpm-alpine

VOLUME /var/run/docker.sock /docker.sock
VOLUME /ak/projects/www

ADD common/container-discovery-wrapper /bin/container-discovery-wrapper
RUN chmod +x /bin/container-discovery-wrapper

CMD /bin/container-discovery-wrapper \
  --command="php-fpm| |-F" \
  --service-name="php-fpm" \
  --service-exposed-port=9000 \
  --service-dependency="bootstrapper" \
  --discovery-server="$DISCOVERY_SERVER" \
  --ak-instance-identifier="$AK_INSTANCE_IDENTIFIER"
