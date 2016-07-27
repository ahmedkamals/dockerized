FROM alpine:latest

ADD . /tmp/provision

RUN sh /tmp/provision/codebase/provision.sh

VOLUME ["/ak/projects/www"]

ADD common/container-discovery-wrapper /bin/container-discovery-wrapper
RUN chmod +x /bin/container-discovery-wrapper

CMD /bin/container-discovery-wrapper \
  --command="/root/block-indefinitely" \
  --service-name="codebase" \
  --service-exposed-port=999 \
  --discovery-server="$DISCOVERY_SERVER" \
  --ak-instance-identifier="$AK_INSTANCE_IDENTIFIER"
