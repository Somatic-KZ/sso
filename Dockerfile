#
# Контейнер сборки
#
FROM golang:latest as builder

ARG DRONE
ARG DRONE_TAG
ARG DRONE_COMMIT
ARG DRONE_BRANCH

ENV CGO_ENABLED=0

COPY . /go/src/gitea.technodom.kz/intechno/techo-sso
WORKDIR /go/src/gitea.technodom.kz/intechno/techo-sso
RUN \
    if [ -z "$DRONE" ] ; then echo "no drone" && version=`git describe --abbrev=6 --always --tag`; \
    else version=${DRONE_TAG}${DRONE_BRANCH}-`echo ${DRONE_COMMIT} | cut -c 1-7` ; fi && \
    echo "version=$version" && \
    cd cmd/apiserver && \
    go build -a -tags techno-sso -installsuffix techno-sso -ldflags "-X main.version=${version} -s -w" -o /go/bin/techno-sso

#
# Контейнер для получения актуальных SSL/TLS сертификатов
#
FROM alpine as alpine
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
RUN addgroup -S sso && adduser -S sso -G sso

# копируем документацию
RUN mkdir -p /usr/share/techno-sso
COPY --from=builder /go/src/gitea.technodom.kz/intechno/techo-sso/api /usr/share/techno-sso
RUN chown -R sso:sso /usr/share/techno-sso

ENTRYPOINT [ "/bin/techno-sso" ]

#
# Контейнер рантайма
#
FROM scratch
COPY --from=builder /go/bin/techno-sso /bin/techno-sso

# копируем сертификаты из alpine
COPY --from=alpine /etc/ssl/certs /etc/ssl/certs

# копируем документацию
COPY --from=alpine /usr/share/techno-sso /usr/share/techno-sso

# копируем пользователя и группу из alpine
COPY --from=alpine /etc/passwd /etc/passwd
COPY --from=alpine /etc/group /etc/group

ENTRYPOINT ["/bin/techno-sso"]



