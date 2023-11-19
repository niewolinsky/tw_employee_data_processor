package utils

import (
	"encoding/json"
	"net/http"
)

type Wrap map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Wrap, headers http.Header) error {
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	json = append(json, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return nil
}

func WriteJSONBytes(w http.ResponseWriter, status int, data Wrap, headers http.Header) ([]byte, error) {
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil, err
	}

	json = append(json, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return json, nil
}

func WriteJSONCache(w http.ResponseWriter, status int, data Wrap, headers http.Header) ([]byte, error) {
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil, err
	}

	json = append(json, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return json, nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, source interface{}) error {
	requestLimit := 5242880
	r.Body = http.MaxBytesReader(w, r.Body, int64(requestLimit))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(source)
	if err != nil {
		return err
	}

	return nil
}
