FROM mysql:latest

ENV MYSQL_DATABASE "reldb"
ENV MYSQL_ROOT_PASSWORD "buzzword"

ADD common/container-discovery-wrapper /bin/container-discovery-wrapper
RUN chmod +x /bin/container-discovery-wrapper

VOLUME /var/run/docker.sock /docker.sock

CMD ""
ENTRYPOINT /bin/container-discovery-wrapper --command="/entrypoint.sh| |mysqld" \
  --service-name="datastore-relational" \
  --service-exposed-port=3306 \
  --announce-after-output-text="datastore-relational[mysqld]: ready to accept connections" \
  --discovery-server="$DISCOVERY_SERVER" \
  --ak-instance-identifier="$AK_INSTANCE_IDENTIFIER"
