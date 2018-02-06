FROM jfloff/alpine-python:2.7

COPY . /app/
WORKDIR /app

# Copy in the entrypoint script -- this installs prerequisites on container start.
COPY .docker-entrypoint.sh /entrypoint.sh

# install requirements
# this way when you build you won't need to install again
# and since COPY is cached we don't need to wait
COPY requirements.txt /tmp/requirements.txt
COPY .docker-build-requirements.txt /tmp/build-requirements.txt

# Run the dependencies installer and then allow it to be run again if needed.
RUN /entrypoint.sh -B /tmp/build-requirements.txt -r /tmp/requirements.txt
RUN rm -f /requirements.installed

ENTRYPOINT ["/usr/bin/python", "-mion"]
