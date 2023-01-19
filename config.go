package arboc

import (
	"reflect"

	"github.com/fatih/structtag"
	"github.com/spf13/cobra"
)

type fieldData struct {
	Field reflect.StructField
	Value reflect.Value
}

type flaggerFunc func(cmd *cobra.Command, field *fieldData, tags *structtag.Tags)

type Flagger struct {
	FuncByKind map[reflect.Kind]flaggerFunc
}

var defaultFlagger *Flagger = nil

func GenerateFlags(cmd *cobra.Command, variable interface{}) {
	getDefaultFlagger().Process(cmd, variable)
}

func getDefaultFlagger() *Flagger {
	if defaultFlagger == nil {
		defaultFlagger = newDefaultFlagger()
	}

	return defaultFlagger
}

func (flagger *Flagger) Process(cmd *cobra.Command, ptr interface{}) {
	fields := flagger.getFields(ptr)
	flagger.set(cmd, fields)
}

func (flagger *Flagger) getFunctionByKind(k reflect.Kind) flaggerFunc {
	if f, ok := flagger.FuncByKind[k]; ok {
		return f
	}

	return nil
}

func (flagger *Flagger) getFields(ptr interface{}) []*fieldData {
	v := reflect.ValueOf(ptr)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		panic("argument must be a pointer to struct")
	}

	valueObject := v.Elem()

	typeObject := valueObject.Type()

	count := valueObject.NumField()
	results := make([]*fieldData, 0)
	for i := 0; i < count; i++ {
		value := valueObject.Field(i)
		field := typeObject.Field(i)

		if value.CanSet() {
			results = append(results, &fieldData{
				Value: value,
				Field: field,
			})
		}
	}

	return results
}

func (flagger *Flagger) set(cmd *cobra.Command, fields []*fieldData) {
	for _, fd := range fields {
		flagger.setField(cmd, fd)
	}
}

func (flagger *Flagger) setField(cmd *cobra.Command, field *fieldData) {
	if field.Value.CanAddr() {
		tags, err := structtag.Parse(string(field.Field.Tag))
		if err != nil {
			panic(err)
		}

		function := flagger.getFunctionByKind(field.Field.Type.Kind())

		if function == nil {
			return
		}

		function(cmd, field, tags)
	}
}
