# Let's use the offial Golang image

# use official Golang image
FROM golang:1.16.3-alpine3.13

# set working directory
WORKDIR /app

# Copy the source code
COPY . . 

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o api .

#EXPOSE the port
EXPOSE 8000

# Run the executable
CMD ["./api"]

## FROM sets the base image to use. In this case we are using the golang:1.16.3-alpine3.13 image, a lightweight version

## WORKDIR sets the working directory inside the image

## COPY . . copies all the files in the current directory to the working directory

## RUN go get -d -v ./... Is a command to isntall the dependencies before building the image

## RUN go build -o api . build the Go app inside the Image filesystem

## EXPOSE 8000 exposes the port 8000

## CMD ["./api"] sets the command to run when the container starts

