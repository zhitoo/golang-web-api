ARG ELASTIC_VERSION

# https://www.docker.elastic.co/
FROM mirror2.chabokan.net/elasticsearch:${ELASTIC_VERSION}

USER root

RUN set -eux; \
	mkdir /state; \
	chmod 0775 /state; \
	chown elasticsearch:root /state

USER elasticsearch:root

ENTRYPOINT ["/entrypoint.sh"]
