FROM ghcr.io/linuxsuren/api-testing:master

WORKDIR /workspace
COPY e2e/* .
COPY helm/api-testing api-testing
RUN apk add curl openssh-client bash
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
RUN chmod 700 get_helm.sh
RUN ./get_helm.sh
RUN chmod +x entrypoint.sh
RUN chmod +x k8s.sh

CMD [ "/workspace/entrypoint.sh" ]
