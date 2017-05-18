package logrus

//
// Field define a single field in message attachment.
//
type Field struct {
	Short bool
	Title string
	Value string
}

//
// MarshalJSON will convert `field` into a valid JSON.
//
// (1) The conversion will skip empty field Title or Value.
//
func (field Field) MarshalJSON() (out []byte, err error) {
	str := "{"

	// (1)
	if len(field.Title) == 0 || len(field.Value) == 0 {
		goto out
	}

	if field.Short {
		str += `"short":true,`
	} else {
		str += `"short":false,`
	}

	str += `"title": "` + field.Title + `",`
	str += `"value": "` + field.Value + `"`

out:
	str += "}"
	out = []byte(str)

	return
}
