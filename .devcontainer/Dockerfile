FROM auguwu/devcontainers:latest

# Setup Go v1.17.5
RUN mkdir -p /opt/golang
RUN curl -X GET -L -o /tmp/golang.tar.gz https://go.dev/dl/go1.17.5.linux-amd64.tar.gz
RUN tar -xvf /tmp/golang.tar.gz -C /opt/golang
RUN rm /tmp/golang
RUN sudo install
ENV PATH="/opt/golang/go1.17.5.linux-amd64/bin:$PATH"

CMD ["sleep", "infinity"]
