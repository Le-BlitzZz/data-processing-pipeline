FROM python:3.13.3-alpine3.21

# Add dependencies
RUN apk add --no-cache bash make curl
RUN curl -sSL https://install.python-poetry.org | python3 - && \
    ln -s ~/.local/bin/poetry /usr/local/bin/poetry && \
    poetry config virtualenvs.create false

CMD ["tail", "-f", "/dev/null"]
