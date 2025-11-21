FROM alpine:3.11.3
ARG SLACK_WEBHOOK_URL
ENV SLACK_WEBHOOK_URL_DOCKER_ENV $SLACK_WEBHOOK_URL
ENV PRE_RECEIVE_HOOK_ENV_VERSION=3.0
RUN \
  RED='\e[31m' && BOLD='\e[1m' && NOCOLOR='\e[0m' && \
  (test -n "$SLACK_WEBHOOK_URL" || \
       (echo -e "${RED}${BOLD}Must set --build-arg SLACK_WEBHOOK_URL=... when running 'docker build' command.${NOCOLOR}" && exit 1)) && \
  apk add --no-cache git openssh bash curl go jq python3 && \
  export GO111MODULE=on && \
  addgroup githook && \
  adduser githook -D -G githook -h /home/githook -s /bin/bash && \
  passwd -d githook && \
  echo "$PRE_RECEIVE_HOOK_ENV_VERSION" > /home/githook/pre_receive_hook_env_version.txt && \
  chown githook:githook /home/githook/pre_receive_hook_env_version.txt
COPY src/detectcreds /home/githook/src/detectcreds
RUN chown -R githook:githook /home/githook
# RUN su githook -c "export GO111MODULE=on && export GOPATH=/home/githook && pushd /home/githook/src/ && go mod download && popd"
RUN su githook -c "export GO111MODULE=on && export GOPATH=/home/githook && pushd /home/githook/src/detectcreds && go build -o /home/githook/detectcreds.bin && popd"
RUN su githook -c "pushd /home/githook && chmod u+w -R pkg && rm -rf pkg && rm -rf src && popd"
# RUN su githook -c "echo \"${SLACK_WEBHOOK_URL_DOCKER_ENV}\" > /home/githook/slack_webhook_url.txt"
RUN su githook -c "echo -e \"#!/bin/bash\necho ${SLACK_WEBHOOK_URL_DOCKER_ENV}\" > /home/githook/slack_webhook_url.txt"
RUN su githook -c "echo -e \"#!/bin/bash\necho TEST\" > /home/githook/test.txt"
RUN su githook -c "chmod +x /home/githook/test.txt"
RUN su githook -c "chmod +x /home/githook/slack_webhook_url.txt"
RUN chown -R githook:githook /home/githook
# mkdir -p /home/githook/ && \
RUN echo 'abc' > /test.txt
RUN echo 'abc' > test.txt
# RUN echo 'abc' > /dev/test.txt
RUN echo 'abc' > /tmp/test.txt
RUN echo 'abc' > /home/githook/test2.txt
RUN chown -R githook:githook /home/githook

VOLUME ["/home/githook"]
WORKDIR /home/githook

CMD ["/usr/sbin/sshd", "-D"]
