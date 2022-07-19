FROM golang:1.17.9 as golang

#use go module
ENV GO111MODULE=on 

#Specify the directory to run the application
WORKDIR /chack-server

#Go to the above directory.mod and go.copy sum
COPY . .

#Cache can be used if the above file is unchanged
RUN go mod download
RUN go build .

RUN go get github.com/pilu/fresh

EXPOSE 1323

#Start the server with the fresh command
CMD ["fresh"]