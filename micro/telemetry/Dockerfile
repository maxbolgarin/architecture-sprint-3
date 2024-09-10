## Based on https://github.com/gbaeke/go-template/blob/main/Dockerfile
FROM golang:1.23-alpine AS build

### Step 1: prepare

# 1. Install git for go.mod
RUN apk add --no-cache git

# 2. Install newest certs
RUN apk --no-cache add ca-certificates

### Step 2: build

# 1. Copy source code to /src folder
COPY . /src
WORKDIR /src

# 2. Download dependencies
RUN go mod download

# 3. Build application with CGO_ENABLED=0 to disable cgo and static linking
RUN CGO_ENABLED=0 go build -o /app/telemetry .


### Step 3: create execuitable image
FROM alpine AS final 

# 1. Copy compiled app and configs
COPY --from=build /app/telemetry /app

# 2. Copy ca certs
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 4. Set entrypoint for running the image
ENTRYPOINT ["/app"]
