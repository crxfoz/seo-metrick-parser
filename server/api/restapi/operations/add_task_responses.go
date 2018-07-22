// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// AddTaskCreatedCode is the HTTP code returned for type AddTaskCreated
const AddTaskCreatedCode int = 201

/*AddTaskCreated Task created

swagger:response addTaskCreated
*/
type AddTaskCreated struct {

	/*
	  In: Body
	*/
	Payload *AddTaskCreatedBody `json:"body,omitempty"`
}

// NewAddTaskCreated creates AddTaskCreated with default headers values
func NewAddTaskCreated() *AddTaskCreated {

	return &AddTaskCreated{}
}

// WithPayload adds the payload to the add task created response
func (o *AddTaskCreated) WithPayload(payload *AddTaskCreatedBody) *AddTaskCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add task created response
func (o *AddTaskCreated) SetPayload(payload *AddTaskCreatedBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddTaskCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddTaskUnprocessableEntityCode is the HTTP code returned for type AddTaskUnprocessableEntity
const AddTaskUnprocessableEntityCode int = 422

/*AddTaskUnprocessableEntity Invalid request data

swagger:response addTaskUnprocessableEntity
*/
type AddTaskUnprocessableEntity struct {
}

// NewAddTaskUnprocessableEntity creates AddTaskUnprocessableEntity with default headers values
func NewAddTaskUnprocessableEntity() *AddTaskUnprocessableEntity {

	return &AddTaskUnprocessableEntity{}
}

// WriteResponse to the client
func (o *AddTaskUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(422)
}

// AddTaskBadGatewayCode is the HTTP code returned for type AddTaskBadGateway
const AddTaskBadGatewayCode int = 502

/*AddTaskBadGateway Database is down

swagger:response addTaskBadGateway
*/
type AddTaskBadGateway struct {
}

// NewAddTaskBadGateway creates AddTaskBadGateway with default headers values
func NewAddTaskBadGateway() *AddTaskBadGateway {

	return &AddTaskBadGateway{}
}

// WriteResponse to the client
func (o *AddTaskBadGateway) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(502)
}