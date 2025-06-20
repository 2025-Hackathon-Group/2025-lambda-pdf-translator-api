FROM golang:1.24.4-bookworm AS builder

LABEL maintainer="Sam Laister <saml@everbit.dev>"
LABEL version="1.0"
LABEL description="PDF Translator API"

WORKDIR /app

COPY . . 

RUN go mod download && go mod verify

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    go install github.com/mitranim/gow@latest

# Install Python and pdf2docx
RUN apt-get update && apt-get install -y python3 python3-pip python3-venv
RUN python3 -m venv /opt/venv && \
    /opt/venv/bin/pip install pdf2docx

# Add the virtual environment to PATH
ENV PATH="/opt/venv/bin:$PATH"

CMD ["make"]