# syntax=docker/dockerfile:1.2
FROM ubuntu:20.04 AS base

RUN apt update && apt install -y git && `rm -rf /var/lib/apt/lists/*`

FROM base AS dev

RUN apt update && apt install -y build-essential && `rm -rf /var/lib/apt/lists/*`

RUN --mount=type=secret,id=mysecret cat /run/secrets/mysecret

RUN make

FROM base AS prod

COPY --from=dev /binary /binary

CMD ["/binary"]
