FROM php:7.0.8-fpm-alpine

VOLUME /var/run/docker.sock /docker.sock
VOLUME /ak/projects/www

ADD execution/periodical_tasks/go-cron /root/go-cron
RUN chmod +x /root/go-cron

ADD common/container-discovery-wrapper /bin/container-discovery-wrapper
RUN chmod +x /bin/container-discovery-wrapper

CMD /bin/container-discovery-wrapper \
  --command="/root/go-cron" \
  --service-name="periodical-tasks" \
  --service-exposed-port=1 \
  --service-dependency="bootstrapper" \
  --discovery-server="$DISCOVERY_SERVER" \
  --ak-instance-identifier="$AK_INSTANCE_IDENTIFIER"
