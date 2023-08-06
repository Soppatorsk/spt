FROM ubuntu

WORKDIR /usr/src/spt/

RUN apt-get update && \ 
    apt-get install -y golang npm nodejs git ca-certificates imagemagick

COPY go.mod go.sum ./

RUN go mod download && go mod verify

ENV PORT=3000
ENV SPOTIFY_KEY=[REDACTED]
ENV OPENAI_API_KEY=[REDACTED]
EXPOSE 3000

COPY . .

RUN go build -v -o /usr/local/bin/spt

RUN cd vue-front && npm install vue@3.2.26 && npm install && npm run build

CMD ["spt"]