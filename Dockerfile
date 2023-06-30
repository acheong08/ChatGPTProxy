FROM golang
RUN go install github.com/acheong08/ChatGPTProxy@latest
CMD [ "ChatGPTProxy" ]
