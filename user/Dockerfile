FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.22

WORKDIR /app

# Install certificates
RUN apk add --no-cache ca-certificates

# Copy binary
COPY --from=builder /app/server .

# Expose application port
EXPOSE 8080

# Start application
CMD ["./server"]


############ Meaning of each line in the Dockerfile:

*** 1. FROM golang:1.25 AS builder ***
`FROM golang:1.25 AS builder`: This line specifies the base image.
-> start my image from this existing image
-> withou this 
FROM alpine
RUN go version
-> you will get 
go: command not found
-> so it is doing
```
Download image
↓
golang:1.25
↓
Start a container from it
↓
Execute the remaining instructions
```
-> With 
AS builder
-> it is making it to import AS Builder

*** 2. WORKDIR /app ***
-> this is not created locally but it is created in the container
-> from this point onwards execute all commands in the /app directory
-> it is equivalent to
mkdir -p /app
cd /app
-> it is helpful so that again and again we donot have to write /app in every command

*** 3. COPY go.mod go.sum ./ ***
-> copying these first bcz it work on caching and layer by layer caching
-> because these will not change frequently so it will be cached and next time if you build the image it will not download again

*** 4. RUN go mod download ***
-> it will download all the dependencies mentioned in go.mod and go.sum

*** 5. COPY . . ***
-> copy all the files from the current directory to /app in the container

*** 6. RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server ***
-> this is very important line
-> CGO_ENABLED=0
-> it is disabling CGO which is used for calling C code from Go
-> in go sometime you import C code and it will not work in alpine image so we are disabling it
-> as we are also not using any C code in our project so it is safe to disable it
-> also it will save some space in the final image
-> GOOS=linux
-> it is setting the target OS to linux
GOOS examples
-> for linux
GOOS=linux
-> for mac
GOOS=darwin
-> for windows
GOOS=windows
-> they will not work interchangeably
-> go build -o server ./cmd/server
-> go build 
-> building the binary 
-> - o server
-> instead of giving it default name it will give the name server to the binary
-> .cmd/server
-> build all the files that are imported here


*** 7. FROM alpine:3.22 ***
-> why we imported this as we have already imorted golang image
-> because golang image is very big in size and we do not need it in the final 
-> as now we donot need compiler so this image does not have go compiler and is very small in size
Running it doesnot require:
Go
go command
compiler
source code
It is like a C program.

*** 8. WORKDIR /app ***

*** 9. RUN apk add --no-cache ca-certificates ***

While building the Alpine image, install the trusted SSL/TLS certificate bundle 
without keeping the package manager cache, so the final image stays small
```
apk

apk is Alpine Linux package manager.

Think of it like:

Operating System	Package Manager
Ubuntu	apt
Debian	apt
CentOS	yum
Fedora	dnf
Alpine	apk
```
-> add means installing the package
-> add git means installing git
-> Normally, apk downloads a package index and keeps it in the image.
-> it doesnot keep the downloaded cache after installation.
->Benefit: Smaller image size.

-> ca-certificates
-> installs a collection of trusted Certificate Authority (CA) certificates.
-> certificates are used to verify the identity of HTTPS websites.
-> when your code call
-> http.Get("https://api.github.com")
-> TLS handshake checks whether GitHubs certificate was issued by a trusted CA.
✅ If ca-certificates is installed, the verification succeeds.
❌ If its missing, you may see an error like:
x509: certificate signed by unknown authority

*** 10. COPY --from=builder /app/server . ***
-> --builder because /app/server is in the builder 
-> this line is copying the binary from the builder stage to the final image
COPY --from=<stage-name> <source> <destination>
-> . refer to working directory which is /app in this case

*** 11. EXPOSE 8080 ***
-> it does not means exposing the port to external world
-> here it means this container is expecting to listen on this port
-> it just provides the information that this container is listening on this port
-> nothing else it is doing

*** 12. CMD ["./server"] ***
-> just executing ./server binary
-> can RUN ./server wroks
-> you want
Build image
↓
Copy server
↓
Finish image
-> if you use RUN then your server runs 
-> and then server is waiting for request and it will not finish the image
-> so we use CMD which will run the server when the container is started
-> cmd means it is just copied when someone run the container it will run this command
-> so as this is building phase that is why we are using CMD instead of RUN


```
             podman build
                  │
                  ▼

        Temporary Container 1
     (from golang:1.25 image)

     ┌─────────────────────┐
     │ Go compiler         │
     │ Source code         │
     │ go build            │
     │ server binary       │
     └─────────────────────┘
              │
              │ COPY --from=builder
              ▼
        Temporary Container 2
      (from alpine:3.22 image)

     ┌─────────────────────┐
     │ ca-certificates     │
     │ /app/server         │
     └─────────────────────┘
              │
              ▼
        Final Image Created
```

What happens to these containers?
Container 1 (builder) → Used only to compile your Go code. After the build, it is discarded.
Container 2 (runtime) → Used to assemble the final image. The final image is based on this stage.
Container 3 (application) → Your actual application container when you execute podman run
Docker/Podman creates a third container, which is the actual running application.