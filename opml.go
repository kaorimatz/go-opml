package opml

import (
	"encoding/xml"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type xmlOPML struct {
	XMLName xml.Name `xml:"opml"`
	Head    xmlHead  `xml:"head"`
	Body    xmlBody  `xml:"body"`
}

func (xo *xmlOPML) ToOPML() *OPML {
	o := &OPML{
		Version:         xo.Head.Version,
		Title:           xo.Head.Title,
		OwnerName:       xo.Head.OwnerName,
		OwnerEmail:      xo.Head.OwnerEmail,
		OwnerID:         (*url.URL)(xo.Head.OwnerID),
		Docs:            (*url.URL)(xo.Head.Docs),
		VertScrollState: xo.Head.VertScrollState,
		WindowTop:       xo.Head.WindowTop,
		WindowLeft:      xo.Head.WindowLeft,
		WindowBottom:    xo.Head.WindowBottom,
		WindowRight:     xo.Head.WindowRight,
		Outlines:        xo.Body.Outlines.ToOutlines(),
	}
	if xo.Head.ExpansionState == nil {
		o.ExpansionState = nil
	} else {
		o.ExpansionState = xo.Head.ExpansionState.expansionState
	}
	if xo.Head.DateCreated == nil {
		o.DateCreated = time.Time{}
	} else {
		o.DateCreated = time.Time(*xo.Head.DateCreated)
	}
	if xo.Head.DateModified == nil {
		o.DateModified = time.Time{}
	} else {
		o.DateModified = time.Time(*xo.Head.DateModified)
	}
	return o
}

func (xo *xmlOPML) FromOPML(o *OPML) {
	xo.Head = xmlHead{
		Version:         o.Version,
		Title:           o.Title,
		OwnerName:       o.OwnerName,
		OwnerEmail:      o.OwnerEmail,
		OwnerID:         (*xmlURL)(o.OwnerID),
		Docs:            (*xmlURL)(o.Docs),
		VertScrollState: o.VertScrollState,
		WindowTop:       o.WindowTop,
		WindowLeft:      o.WindowLeft,
		WindowBottom:    o.WindowBottom,
		WindowRight:     o.WindowRight,
	}
	if o.ExpansionState == nil {
		xo.Head.ExpansionState = nil
	} else {
		xo.Head.ExpansionState = &xmlExpansionState{o.ExpansionState}
	}
	if o.DateCreated.IsZero() {
		xo.Head.DateCreated = nil
	} else {
		xo.Head.DateCreated = (*xmlTime)(&o.DateCreated)
	}
	if o.DateModified.IsZero() {
		xo.Head.DateModified = nil
	} else {
		xo.Head.DateModified = (*xmlTime)(&o.DateModified)
	}
	xo.Body.Outlines.FromOutlines(o.Outlines)
}

type xmlHead struct {
	Version         string             `xml:"version,attr"`
	Title           string             `xml:"title,omitempty"`
	DateCreated     *xmlTime           `xml:"dateCreated,omitempty"`
	DateModified    *xmlTime           `xml:"dateModified,omitempty"`
	OwnerName       string             `xml:"ownerName,omitempty"`
	OwnerEmail      string             `xml:"ownerEmail,omitempty"`
	OwnerID         *xmlURL            `xml:"ownerId,omitempty"`
	Docs            *xmlURL            `xml:"docs,omitempty"`
	ExpansionState  *xmlExpansionState `xml:"expansionState,omitempty"`
	VertScrollState int                `xml:"vertScrollState,omitempty"`
	WindowTop       int                `xml:"windowTop,omitempty"`
	WindowLeft      int                `xml:"windowLeft,omitempty"`
	WindowBottom    int                `xml:"windowBottom,omitempty"`
	WindowRight     int                `xml:"windowRight,omitempty"`
}

type xmlBody struct {
	Outlines xmlOutlines `xml:"outline"`
}

type xmlOutline struct {
	Text         string        `xml:"text,attr"`
	Type         string        `xml:"type,attr,omitempty"`
	IsComment    bool          `xml:"isComment,attr,omitempty"`
	IsBreakpoint bool          `xml:"isBreakpoint,attr,omitempty"`
	Created      *xmlTime      `xml:"created,attr,omitempty"`
	Categories   xmlCategories `xml:"category,attr,omitempty"`
	XMLURL       *xmlURL       `xml:"xmlUrl,attr,omitempty"`
	Description  string        `xml:"description,attr,omitempty"`
	HTMLURL      *xmlURL       `xml:"htmlUrl,attr,omitempty"`
	Language     string        `xml:"language,attr,omitempty"`
	Title        string        `xml:"title,attr,omitempty"`
	Version      string        `xml:"version,attr,omitempty"`
	URL          *xmlURL       `xml:"url,attr,omitempty"`
	Outlines     xmlOutlines   `xml:"outline,omitempty"`
}

func (xo *xmlOutline) ToOutline() *Outline {
	o := &Outline{
		Text:         xo.Text,
		Type:         xo.Type,
		IsComment:    xo.IsComment,
		IsBreakpoint: xo.IsBreakpoint,
		Categories:   []string(xo.Categories),
		XMLURL:       (*url.URL)(xo.XMLURL),
		Description:  xo.Description,
		HTMLURL:      (*url.URL)(xo.HTMLURL),
		Language:     xo.Language,
		Title:        xo.Title,
		Version:      xo.Version,
		URL:          (*url.URL)(xo.URL),
		Outlines:     xo.Outlines.ToOutlines(),
	}
	if xo.Created == nil {
		o.Created = time.Time{}
	} else {
		o.Created = time.Time(*xo.Created)
	}
	return o
}

func (xo *xmlOutline) FromOutline(o *Outline) {
	xo.Text = o.Text
	xo.Type = o.Type
	xo.IsComment = o.IsComment
	xo.IsBreakpoint = o.IsBreakpoint
	if o.Created.IsZero() {
		xo.Created = nil
	} else {
		xo.Created = (*xmlTime)(&o.Created)
	}
	xo.Categories = xmlCategories(o.Categories)
	xo.XMLURL = (*xmlURL)(o.XMLURL)
	xo.Description = o.Description
	xo.HTMLURL = (*xmlURL)(o.HTMLURL)
	xo.Language = o.Language
	xo.Title = o.Title
	xo.Version = o.Version
	xo.URL = (*xmlURL)(o.URL)
	xo.Outlines.FromOutlines(o.Outlines)
}

type xmlOutlines []*xmlOutline

func (xos xmlOutlines) ToOutlines() []*Outline {
	if xos == nil {
		return nil
	}

	outlines := make([]*Outline, len(xos))
	for i, o := range xos {
		outlines[i] = o.ToOutline()
	}
	return outlines
}

func (xos *xmlOutlines) FromOutlines(os []*Outline) {
	for _, o := range os {
		var xo xmlOutline
		xo.FromOutline(o)
		*xos = append(*xos, &xo)
	}
}

type xmlTime time.Time

var (
	xmlTimeLayouts = [...]string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}
)

func tryParseTime(v string) (time.Time, error) {
	var err error
	for _, l := range xmlTimeLayouts {
		t, err := time.Parse(l, v)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, err
}

func (xt *xmlTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	t, err := tryParseTime(v)
	if err != nil {
		return err
	}

	*xt = xmlTime(t)
	return nil
}

func (xt *xmlTime) UnmarshalXMLAttr(attr xml.Attr) error {
	t, err := tryParseTime(attr.Value)
	if err != nil {
		return err
	}

	*xt = xmlTime(t)
	return nil
}

func (xt *xmlTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement((*time.Time)(xt).Format(time.RFC1123), start)
}

func (xt *xmlTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: (*time.Time)(xt).Format(time.RFC1123)}, nil
}

type xmlURL url.URL

func (u *xmlURL) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	url, err := url.Parse(v)
	if err != nil {
		return err
	}

	*u = xmlURL(*url)
	return err
}

func (u *xmlURL) UnmarshalXMLAttr(attr xml.Attr) error {
	url, err := url.Parse(attr.Value)
	if err != nil {
		return err
	}

	*u = xmlURL(*url)
	return nil
}

func (u *xmlURL) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement((*url.URL)(u).String(), start)
}

func (u *xmlURL) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: (*url.URL)(u).String()}, nil
}

type xmlExpansionState struct {
	expansionState []int
}

func (s *xmlExpansionState) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}

	state := []int{}
	for _, str := range strings.Split(v, ",") {
		trimmed := strings.TrimSpace(str)
		if len(trimmed) == 0 {
			continue
		}
		n, err := strconv.Atoi(trimmed)
		if err != nil {
			return err
		}
		state = append(state, n)
	}

	*s = xmlExpansionState{state}
	return nil
}

func (s xmlExpansionState) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	states := make([]string, len(s.expansionState))
	for i, n := range s.expansionState {
		states[i] = strconv.Itoa(n)
	}
	return e.EncodeElement(strings.Join(states, ","), start)
}

type xmlCategories []string

func (c *xmlCategories) UnmarshalXMLAttr(attr xml.Attr) error {
	var categories []string

	for _, str := range strings.Split(attr.Value, ",") {
		trimmed := strings.TrimSpace(str)
		if len(trimmed) == 0 {
			continue
		}
		categories = append(categories, trimmed)
	}

	*c = xmlCategories(categories)
	return nil
}

func (c xmlCategories) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: strings.Join([]string(c), ",")}, nil
}

type OPML struct {
	Version         string
	Title           string
	DateCreated     time.Time
	DateModified    time.Time
	OwnerName       string
	OwnerEmail      string
	OwnerID         *url.URL
	Docs            *url.URL
	ExpansionState  []int
	VertScrollState int
	WindowTop       int
	WindowLeft      int
	WindowBottom    int
	WindowRight     int
	Outlines        []*Outline
}

type Outline struct {
	Text         string
	Type         string
	IsComment    bool
	IsBreakpoint bool
	Created      time.Time
	Categories   []string
	XMLURL       *url.URL
	Description  string
	HTMLURL      *url.URL
	Language     string
	Title        string
	Version      string
	URL          *url.URL
	Outlines     []*Outline
}

type Parser struct {
	XMLDecoder *xml.Decoder
}

func NewParser(r io.Reader) *Parser {
	return &Parser{XMLDecoder: xml.NewDecoder(r)}
}

func (p *Parser) Parse() (*OPML, error) {
	var xmlOPML xmlOPML
	if err := p.XMLDecoder.Decode(&xmlOPML); err != nil {
		return nil, err
	}
	return xmlOPML.ToOPML(), nil
}

func Parse(r io.Reader) (*OPML, error) {
	return NewParser(r).Parse()
}

type Renderer struct {
	XMLEncoder *xml.Encoder
}

func NewRenderer(w io.Writer) *Renderer {
	return &Renderer{XMLEncoder: xml.NewEncoder(w)}
}

func (r *Renderer) Render(opml *OPML) error {
	var xmlOPML xmlOPML
	xmlOPML.FromOPML(opml)
	return r.XMLEncoder.Encode(xmlOPML)
}

func Render(w io.Writer, opml *OPML) error {
	return NewRenderer(w).Render(opml)
}
