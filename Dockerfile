FROM golang
RUN go install github.com/acheong08/ChatGPT-Proxy-V4@latest
CMD [ "ChatGPT-Proxy-V4" ]
