FROM ubuntu

WORKDIR /usr/src/spt/

RUN apt-get update && \ 
    apt-get install -y golang npm git ca-certificates imagemagick

COPY go.mod go.sum ./

RUN go mod download && go mod verify

ENV PORT=3000
ENV SPOTIFY_KEY=YOURKEYHERE
ENV OPENAI_API_KEY=YOURKEYHERE
EXPOSE 3000

COPY . .

RUN go build -v -o /usr/local/bin/spt

# wont run idk just build beforehand lol
# RUN npm install vue && \ 
#     cd vue-front && \ 
#     npm install && npm run build

CMD ["spt"]