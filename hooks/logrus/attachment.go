package logrus

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	_colorsLevel = []string{
		"#FF0000", // Panic
		"#CC0000", // Fatal
		"#990000", // Error
		"#9F6000", // Warning
		"#FFFFFF", // Info
		"#000000", // Debug
	}
)

//
// Attachment define Mattermost message attachment [1].
//
// [1] https://docs.mattermost.com/developer/message-attachments.html
//
type Attachment struct {
	AuthorIcon string
	AuthorLink string
	AuthorName string
	Color      string
	Fallback   string
	Fields     Fields
	ImageURL   string
	Pretext    string
	Text       string
	Title      string
	TitleLink  string
}

//
// NewAttachment will create and return new Attachment with default value set
// from `attc` and Color and Fields based on logrus Entry Level and Data.
//
func NewAttachment(defAttc *Attachment, entry *logrus.Entry) (
	attc *Attachment,
) {
	if defAttc == nil {
		return
	}

	attc = &Attachment{
		AuthorIcon: defAttc.AuthorIcon,
		AuthorLink: defAttc.AuthorLink,
		AuthorName: defAttc.AuthorName,
		ImageURL:   defAttc.ImageURL,
		Pretext:    defAttc.Pretext,
		Title:      defAttc.Title,
		TitleLink:  defAttc.TitleLink,
	}

	if entry != nil {
		attc.Color = _colorsLevel[entry.Level]
		attc.SetFields(entry.Data)
		attc.Text = _iconsLevel[entry.Level] + " " + entry.Message
	}

	return
}

func (attc Attachment) marshalAuthor() (str string) {
	if len(attc.AuthorIcon) > 0 {
		str += `"author_icon":"` + attc.AuthorIcon + `",`
	}
	if len(attc.AuthorLink) > 0 {
		str += `"author_link":"` + attc.AuthorLink + `",`
	}
	if len(attc.AuthorName) > 0 {
		str += `"author_name":"` + attc.AuthorName + `",`
	}

	return
}

//
// MarshalJSON will convert Attachment `attc` to JSON.
//
func (attc Attachment) MarshalJSON() (out []byte, err error) {
	var bFields []byte

	bFields, err = attc.Fields.MarshalJSON()
	if err != nil {
		return
	}

	str := "{"

	str += attc.marshalAuthor()

	if len(attc.Color) > 0 {
		str += `"color":"` + attc.Color + `",`
	}
	if len(attc.Fallback) > 0 {
		str += `"fallback":"` + attc.Fallback + `",`
	}
	if len(attc.Fields) > 0 {
		str += `"fields":` + string(bFields) + `,`
	}
	if len(attc.ImageURL) > 0 {
		str += `"image_url":"` + attc.ImageURL + `",`
	}
	if len(attc.Pretext) > 0 {
		str += `"pretext":"` + attc.Pretext + `",`
	}
	if len(attc.Text) > 0 {
		str += `"text":"` + attc.Text + `",`
	}
	if len(attc.Title) > 0 {
		str += `"title":"` + attc.Title + `",`
	}
	if len(attc.TitleLink) > 0 {
		str += `"title_link":"` + attc.TitleLink + `"`
	} else {
		str = strings.TrimRight(str, ",")
	}

	str += "}"
	out = []byte(str)

	return
}

//
// SetFields will convert logrus Fields data `in` into our Fields.
//
func (attc *Attachment) SetFields(in logrus.Fields) {
	attc.Fields = make(Fields, 0)

	if len(in) == 0 {
		return
	}

	for k, v := range in {
		attc.Fields = append(attc.Fields, Field{
			Short: true,
			Title: k,
			Value: fmt.Sprintf("%v", v),
		})
	}
}
