## Based on https://github.com/gbaeke/go-template/blob/main/Dockerfile
FROM golang:1.23-alpine AS build

### Step 1: prepare

# 1. Install git for go.mod
RUN apk add --no-cache git

# 2. Install newest certs
RUN apk --no-cache add ca-certificates

# 3. Add a group and a user to make it more secure 
# https://www.docker.com/blog/understanding-the-docker-user-instruction/
# https://stackoverflow.com/a/55757473/16748093
ENV USER=goapp
ENV GROUPNAME=$USER
ENV UID=12345
ENV GID=23456
RUN addgroup --gid $GID $GROUPNAME
RUN adduser --disabled-password --gecos "" --ingroup $GROUPNAME --no-create-home --uid $UID $USER

### Step 2: build

# 1. Copy source code to /src folder
COPY . /src
WORKDIR /src

# 2. Download dependencies
RUN go mod download

# 3. Build application with CGO_ENABLED=0 to disable cgo and static linking
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /app/telemetry .


### Step 3: create execuitable image

# https://hub.docker.com/_/scratch — empty layer
FROM scratch AS final 

# 1. Copy compiled app and configs
COPY --from=build /app/telemetry /app

# 2. Copy ca certs
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 3. Copy user goapp
COPY --from=0 /etc/passwd /etc/passwd
USER goapp

# 4. Set entrypoint for running the image
ENTRYPOINT ["/app"]
