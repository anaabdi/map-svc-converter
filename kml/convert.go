package kml

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Convert(w http.ResponseWriter, r *http.Request) {
	var req Kml
	if err := parseReq(r, &req); err != nil {
		log.Printf("err: %v", err)
		return
	}

	resp := populateResponse(req)
	respond(w, http.StatusCreated, resp)
}

func populateResponse(req Kml) []Response {
	var landmarks []Response

	for k, v := range req.Document.Folder.Placemark {
		if !v.IsPolygon() {
			if k == 0 {
				log.Printf("invalid structure: expecting polygon but providing a point: %s", v.Name)
				return nil
			}
			continue
		}

		log.Printf("preparing subplaces for: %s", v.Name)

		var subplaces []Subplace
		for k2, v2 := range req.Document.Folder.Placemark {
			if k >= k2 {
				continue
			}

			if v2.IsPolygon() {
				log.Printf("----> done preparing subplaces for: %s", v.Name)
				// done with this landmark, continue to the next one
				break
			}

			log.Printf("subplace is: %s", v2.Name)

			subplaces = append(subplaces, Subplace{
				NameID:       v2.Name,
				NameEn:       v2.Name,
				SubplaceType: SubplaceTypeAll,
				Location:     getSubplaceLocation(v2.Point.Coordinates),
			})
		}

		coordinates := getLandmarkPolygon(v.Polygon.OuterBoundaryIs.LinearRing.Coordinates)
		respGeo := GetGeometry(coordinates)

		landmark := Response{
			Landmark: Landmark{
				NameEn:    v.Name,
				NameID:    v.Name,
				Address:   v.Description,
				StartDate: time.Now().Format(time.RFC3339),
				IsActive:  true,
				City:      "Jakarta",
				Geometry:  respGeo,
			},
			Subplaces: subplaces,
		}

		landmarks = append(landmarks, landmark)
	}

	return landmarks
}

func getSubplaceLocation(coordinates string) (locations []float64) {
	coordinatesPrepared := prepareTheCoordinates(coordinates)
	if len(coordinatesPrepared) != 1 {
		log.Println("invalid point format")
		return
	}

	coords := splitCoordinates(coordinatesPrepared[0])
	for _, v := range coords {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Println(err)
			continue
		}

		if f == 0 {
			continue
		}

		locations = append(locations, f)
	}

	return locations
}

func prepareTheCoordinates(coordinateStr string) []string {
	space := regexp.MustCompile(`\s+`)
	s := space.ReplaceAllString(coordinateStr, " ")
	s = strings.TrimSpace(s)
	return strings.Split(s, " ")
}

func splitCoordinates(coordinate string) []string {
	coordinate = strings.TrimSpace(coordinate)
	return strings.Split(coordinate, ",")
}

func getLandmarkPolygon(coordinates string) [][][]float64 {
	splitted := prepareTheCoordinates(coordinates)

	firstLayer := make([][][]float64, 0)

	secondLayer := make([][]float64, len(splitted))

	for i, v := range splitted {
		secondLayer[i] = make([]float64, 2)
		v2 := splitCoordinates(v)

		for j, v := range v2 {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				log.Println(err)
				continue
			}

			if f == 0 {
				continue
			}

			secondLayer[i][j] = f
		}
	}

	firstLayer = append(firstLayer, secondLayer)
	return firstLayer
}

func parseReq(r *http.Request, body interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/xml") {
		return xml.NewDecoder(r.Body).Decode(&body)
	}

	return errors.New("no supported type")
}

func respond(w http.ResponseWriter, status int, data interface{}) {
	res, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}
