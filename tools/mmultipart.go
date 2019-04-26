package tools

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"github.com/utrack/clay/v2/transport/httpruntime"
)

func MultipartFormMarshalerGetter(params httpruntime.ContentTypeOptions) httpruntime.Marshaler {
	return MarshalerMultipartForm{
		params,
	}
}

type MarshalerMultipartForm struct {
	params map[string]string
}

func (MarshalerMultipartForm) ContentType() string {
	return "multipart/form-data"
}

func (m MarshalerMultipartForm) Unmarshal(reader io.Reader, dst interface{}) error {

	var msg proto.Message
	var ok bool

	// We need this dirty hack because `dst` contains double pointer to real struct.
	// And we can't extract an interface from a double pointer without some reflection
	v := reflect.ValueOf(dst)
	for {
		if kind := v.Kind(); kind != reflect.Ptr && kind != reflect.Interface {
			break
		}
		vv := v.Interface()
		msg, ok = vv.(proto.Message)
		if ok {
			break
		}
		v = v.Elem()
	}
	if msg == nil {
		return fmt.Errorf("Form unmarshaller supports only `proto.Message` request type, got %T", dst)
	}

	var mpReader *multipart.Reader

	fmt.Println(m.params)
	if boundary, ok := m.params["boundary"]; ok {
		mpReader = multipart.NewReader(reader, boundary)
	} else {
		fmt.Print(fmt.Errorf("no boundary in params: %v", m.params))
	}

	vs := make(map[string][]string)
	for {
		part, err := mpReader.NextPart()
		if err != nil { // i.e. err.Error() == "EOF"
			break
		}
		value, err := ioutil.ReadAll(part)
		vs[part.FormName()] = append(vs[part.FormName()], string(value))
	}

	filter := &utilities.DoubleArray{}
	err := runtime.PopulateQueryParameters(msg, vs, filter)
	if err != nil {
		return err
	}

	return nil
}

func (m MarshalerMultipartForm) Marshal(w io.Writer, src interface{}) error {
	return errors.New("not implemented")
}
