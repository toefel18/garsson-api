package api

type GenericResponse struct {
    Code     int
    Message string
    Data     interface{} `json:",omitempty"`
}
