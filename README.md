# Map Svc Converter

## Getting Started

This is just a toy project

to convert geojson or kml format into specific format

### Prerequisites

- install the latest go (https://golang.org/doc/install)
- already have a exported kml format for specific area on map

## How to Run the server

Simple by doing this will open and ready to serve in port 7777

- go run main.go

## How to get the KML format

- Draw you map in https://www.google.com/maps/d/u/0/
- Then export it as KML

## How to convert the KML format

- Copy the content of exported KML file
- Paste it into the request body like below:

```
curl -X POST \
  http://localhost:7777/kml/convert \
  -H 'Content-Type: application/xml' \
  -d 'content_here'
```

## How to convert the GeoJson format
- Copy the content of exported KML file
- Convert it to GeoJson format first in this website (https://mapbox.github.io/togeojson/)
- Copy the result, and paste it into the request body like below:

```
curl -X POST \
  http://localhost:7777/geojson/convert \
  -H 'Content-Type: application/json' \
  -d '{content here}'
```

## API Postman Collections

https://www.getpostman.com/collections/1a90bd8598cf9133af9a

## Authors

* **Abdi Pratama**
