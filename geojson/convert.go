package geojson

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

func Convert(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := parseReq(r, &req); err != nil {
		log.Printf("err: %v", err)
		return
	}

	resp := populateResponse(req)
	respond(w, http.StatusCreated, resp)
}

func populateResponse(req Request) []Response {
	var landmarks []Response

	for k, v := range req.Features {
		if v.Geometry.Type == Point {
			if k == 0 {
				log.Printf("invalid structure: expecting polygon but providing a point: %s", v.Properties.Name)
				return nil
			}
			continue
		}

		log.Printf("preparing subplaces for: %s", v.Properties.Name)

		var subplaces []Subplace
		for k2, v2 := range req.Features {
			if k >= k2 {
				continue
			}

			if v2.Geometry.Type == Polygon {
				log.Printf("----> done preparing subplaces for: %s", v.Properties.Name)
				// done with this landmark, continue to the next one
				break
			}

			log.Printf("subplace is: %s", v2.Properties.Name)

			coordinates, ok := v2.Geometry.Coordinates.([]interface{})
			if !ok {
				log.Printf("%#v", v2.Geometry.Coordinates)
				log.Printf("invalid structure: expected coordinates with 1D array: %s", v2.Properties.Name)
				return nil
			}

			var locations []float64
			for _, coordinate := range coordinates {
				c, ok := coordinate.(float64)
				if !ok {
					log.Printf("invalid structure: expected coordinates with 1D array: %s", v2.Properties.Name)
					return nil
				}

				if c == 0 {
					continue
				}
				locations = append(locations, c)
			}

			subplaces = append(subplaces, Subplace{
				NameID:       v2.Properties.Name,
				NameEn:       v2.Properties.Name,
				SubplaceType: SubplaceTypeAll,
				Location:     locations,
			})
		}

		coordinates, ok := v.Geometry.Coordinates.([]interface{})
		if !ok {
			log.Printf("%T", v.Geometry.Coordinates)
			log.Printf("invalid structure: expected coordinates with 3D array: %s", v.Properties.Name)
			return nil
		}

		respGeo := GetGeometry(coordinates)

		landmark := Response{
			Landmark: Landmark{
				NameEn:    v.Properties.Name,
				NameID:    v.Properties.Name,
				Address:   v.Properties.Description,
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

func parseReq(r *http.Request, body interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		return json.NewDecoder(r.Body).Decode(&body)
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
