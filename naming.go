package cpe

import (
	"fmt"
	"strings"
)

// Item reprecents a CPE item.
type Item struct {
	part       PartAttr
	vendor     StringAttr
	product    StringAttr
	version    StringAttr
	update     StringAttr
	edition    StringAttr
	language   StringAttr
	sw_edition StringAttr
	target_sw  StringAttr
	target_hw  StringAttr
	other      StringAttr
}

// NewItem returns empty Item.
func NewItem() *Item {
	return &Item{
		part:       PartNotSet,
		vendor:     Any,
		product:    Any,
		version:    Any,
		update:     Any,
		edition:    Any,
		language:   Any,
		sw_edition: Any,
		target_sw:  Any,
		target_hw:  Any,
		other:      Any,
	}
}

func NewItemFromWfn(wfn string) (*Item, error) {
	if strings.HasPrefix(wfn, "wfn:[") {
		wfn = strings.TrimPrefix(wfn, "wfn:[")
	} else {
		return nil, cpeerr{reason: err_invalid_wfn}
	}

	if strings.HasSuffix(wfn, "]") {
		wfn = strings.TrimSuffix(wfn, "]")
	} else {
		return nil, cpeerr{reason: err_invalid_wfn}
	}

	item := NewItem()
	for _, attr := range strings.Split(wfn, ",") {
		sepattr := strings.Split(attr, "=")
		if len(sepattr) != 2 {
			return nil, cpeerr{reason: err_invalid_wfn}
		}

		n, v := sepattr[0], sepattr[1]
		switch n {
		case "part":
			item.part = NewPartAttrFromWfnEncoded(v)
		case "vendor":
			item.vendor = NewStringAttrFromWfnEncoded(v)
		case "product":
			item.product = NewStringAttrFromWfnEncoded(v)
		case "version":
			item.version = NewStringAttrFromWfnEncoded(v)
		case "update":
			item.update = NewStringAttrFromWfnEncoded(v)
		case "edition":
			item.edition = NewStringAttrFromWfnEncoded(v)
		case "language":
			item.language = NewStringAttrFromWfnEncoded(v)
		case "sw_edition":
			item.sw_edition = NewStringAttrFromWfnEncoded(v)
		case "target_sw":
			item.target_sw = NewStringAttrFromWfnEncoded(v)
		case "target_hw":
			item.target_hw = NewStringAttrFromWfnEncoded(v)
		case "other":
			item.other = NewStringAttrFromWfnEncoded(v)
		}
	}

	return item, nil
}

func NewItemFromUri(uri string) (*Item, error) {
	if strings.HasPrefix(uri, "cpe:/") {
		uri = strings.TrimPrefix(uri, "cpe:/")
	} else {
		return nil, cpeerr{reason: err_invalid_wfn}
	}

	item := NewItem()
	for i, attr := range strings.Split(uri, ":") {
		switch i {
		case 0:
			item.part = NewPartAttrFromUriEncoded(attr)
		case 1:
			item.vendor = NewStringAttrFromUriEncoded(attr)
		case 2:
			item.product = NewStringAttrFromUriEncoded(attr)
		case 3:
			item.version = NewStringAttrFromUriEncoded(attr)
		case 4:
			item.update = NewStringAttrFromUriEncoded(attr)
		case 5:
			editions := strings.Split(attr, "~")
			if len(editions) == 1 {
				item.edition = NewStringAttrFromUriEncoded(editions[0])
			} else if len(editions) == 6 {
				item.edition = NewStringAttrFromUriEncoded(editions[1])
				item.sw_edition = NewStringAttrFromUriEncoded(editions[2])
				item.target_sw = NewStringAttrFromUriEncoded(editions[3])
				item.target_hw = NewStringAttrFromUriEncoded(editions[4])
				item.other = NewStringAttrFromUriEncoded(editions[5])
			} else {
				return nil, cpeerr{reason: err_invalid_wfn}
			}
		}
	}
	return item, nil
}

func NewItemFromFormattedString(str string) (*Item, error) {
	if strings.HasPrefix(str, "cpe:2.3:") {
		str = strings.TrimPrefix(str, "cpe:2.3:")
	} else {
		return nil, cpeerr{reason: err_invalid_wfn}
	}

	attrs := strings.Split(str, ":")
	if len(attrs) != 11 {
		return nil, cpeerr{reason: err_invalid_wfn}
	}

	item := NewItem()
	for i, attr := range attrs {
		switch i {
		case 0:
			item.part = NewPartAttrFromFmtEncoded(attr)
		case 1:
			item.vendor = NewStringAttrFromFmtEncoded(attr)
		case 2:
			item.product = NewStringAttrFromFmtEncoded(attr)
		case 3:
			item.version = NewStringAttrFromFmtEncoded(attr)
		case 4:
			item.update = NewStringAttrFromFmtEncoded(attr)
		case 5:
			item.edition = NewStringAttrFromFmtEncoded(attr)
		case 6:
			item.language = NewStringAttrFromFmtEncoded(attr)
		case 7:
			item.sw_edition = NewStringAttrFromFmtEncoded(attr)
		case 8:
			item.target_sw = NewStringAttrFromFmtEncoded(attr)
		case 9:
			item.target_hw = NewStringAttrFromFmtEncoded(attr)
		case 10:
			item.other = NewStringAttrFromFmtEncoded(attr)
		}
	}

	return item, nil
}

// Wfn returns a string of Well-Formed string data model.
func (m *Item) Wfn() string {
	wfn := "wfn:["
	first := true

	for _, it := range []struct {
		name string
		attr Attribute
	}{
		{"part", m.part},
		{"vendor", m.vendor},
		{"product", m.product},
		{"version", m.version},
		{"update", m.update},
		{"edition", m.edition},
		{"language", m.language},
		{"sw_edition", m.sw_edition},
		{"target_sw", m.target_sw},
		{"target_hw", m.target_hw},
		{"other", m.other},
	} {
		if !it.attr.IsEmpty() {
			if first {
				first = false
			} else {
				wfn += ","
			}
			wfn += it.name + "=" + it.attr.WFNEncoded()
		}
	}
	wfn += "]"

	return wfn
}

// Wfn returns a string of uri binding.
func (m *Item) Uri() string {
	uri := "cpe:/"

	l := []struct {
		name string
		attr Attribute
	}{
		{"part", m.part},
		{"vendor", m.vendor},
		{"product", m.product},
		{"version", m.version},
		{"update", m.update},
	}

	for c, it := range l {
		if !it.attr.IsEmpty() {
			uri += it.attr.UrlEncoded()
		}
		if c+1 != len(l) {
			uri += ":"
		}
	}

	if m.target_hw.UrlEncoded() != "" ||
		m.target_sw.UrlEncoded() != "" ||
		m.sw_edition.UrlEncoded() != "" ||
		m.other.UrlEncoded() != "" {
		uri += ":~" + m.edition.UrlEncoded()
		uri += "~" + m.sw_edition.UrlEncoded()
		uri += "~" + m.target_sw.UrlEncoded()
		uri += "~" + m.target_hw.UrlEncoded()
		uri += "~" + m.other.UrlEncoded()
	} else {
		uri += ":" + m.edition.UrlEncoded()
	}

	uri += ":" + m.language.UrlEncoded()
	return strings.TrimRight(uri, ":*")
}

// Wfn returns a formatted string binding.
func (m *Item) Formatted() string {
	fmted := "cpe:2.3"

	for _, it := range []Attribute{
		m.part, m.vendor, m.product, m.version, m.update, m.edition, m.language, m.sw_edition, m.target_sw, m.target_hw, m.other,
	} {
		if !it.IsEmpty() {
			fmted += ":" + it.FmtString()
		} else {
			fmted += ":*"

		}
	}
	return fmted
}

// SetPart sets part of item.  returns error if p is invalid.
func (i *Item) SetPart(p PartAttr) error {
	if !p.IsValid() {
		return cpeerr{reason: err_invalid_type, attr: []interface{}{p, "part"}}
	}

	i.part = p
	return nil
}

// Part returns part of item.
func (i *Item) Part() PartAttr {
	return i.part
}

// SetVendor sets vendor of item.  returns error if s is invalid.
func (i *Item) SetVendor(s StringAttr) error {
	if !s.IsValid() {
		return cpeerr{reason: err_invalid_attribute_str}
	}

	i.vendor = s
	return nil
}

// Vendor returns vendor of item.
func (i *Item) Vendor() StringAttr {
	return i.vendor
}

// SetProduct sets vendor of item.  returns error if s is invalid.
func (i *Item) SetProduct(s StringAttr) error {
	if !s.IsValid() {
		return cpeerr{reason: err_invalid_attribute_str}
	}

	i.product = s
	return nil
}

// Vendor returns product of item.
func (i *Item) Product() StringAttr {
	return i.product
}

// SetVersion sets version of item.  returns error if s is invalid.
func (i *Item) SetVersion(s StringAttr) error {
	if !s.IsValid() {
		return cpeerr{reason: err_invalid_attribute_str}
	}

	i.version = s
	return nil
}

// Version returns version of item.
func (i *Item) Version() StringAttr {
	return i.version
}

// SetUpdate sets update of item.  returns error if s is invalid.
func (i *Item) SetUpdate(s StringAttr) error {
	if !s.IsValid() {
		return cpeerr{reason: err_invalid_attribute_str}
	}

	i.update = s
	return nil
}

// Update returns update of item.
func (i *Item) Update() StringAttr {
	return i.update
}

// SetEdition sets edition of item.  returns error if s is invalid.
func (i *Item) SetEdition(s StringAttr) error {
	if !s.IsValid() {
		return cpeerr{reason: err_invalid_attribute_str}
	}

	i.edition = s
	return nil
}

// Edition returns edition of item.
func (i *Item) Edition() StringAttr {
	return i.edition
}

// SetLanguage sets language of item.  returns error if s is invalid.
func (i *Item) SetLanguage(s StringAttr) error {
	if !s.IsValid() {
		return cpeerr{reason: err_invalid_attribute_str}
	}

	i.language = s
	return nil
}

// Language returns language of item.
func (i *Item) Language() StringAttr {
	return i.language
}

// SetSwEdition sets sw_edition of item.  returns error if s is invalid.
func (i *Item) SetSwEdition(s StringAttr) error {
	if !s.IsValid() {
		return cpeerr{reason: err_invalid_attribute_str}
	}

	i.sw_edition = s
	return nil
}

// SwEdition returns sw_edition of item.
func (i *Item) SwEdition() StringAttr {
	return i.sw_edition
}


// SetTargetSw sets target_sw of item.  returns error if s is invalid.
func (i *Item) SetTargetSw(s StringAttr) error {
	if !s.IsValid() {
		return cpeerr{reason: err_invalid_attribute_str}
	}

	i.target_sw = s
	return nil
}

// TargetSw returns target_sw of item.
func (i *Item) TargetSw() StringAttr {
	return i.target_sw
}

// SetTargetHw sets target_hw of item.  returns error if s is invalid.
func (i *Item) SetTargetHw(s StringAttr) error {
	if !s.IsValid() {
		return cpeerr{reason: err_invalid_attribute_str}
	}

	i.target_hw = s
	return nil
}

// TargetHw returns target_hw of item.
func (i *Item) TargetHw() StringAttr {
	return i.target_hw
}

// SetOther sets other of item.  returns error if s is invalid.
func (i *Item) SetOther(s StringAttr) error {
	if !s.IsValid() {
		return cpeerr{reason: err_invalid_attribute_str}
	}

	i.other = s
	return nil
}

// Other returns other of item.
func (i *Item) Other() StringAttr {
	return i.other
}

type cpeerr struct {
	reason string
	attr   []interface{}
}

var (
	err_invalid_type          = "\"%#v\" is not valid as %v attribute."
	err_invalid_attribute_str = "invalid attribute string."
	err_invalid_wfn           = "invalid wfn string."
)

func (e cpeerr) Error() string {
	return fmt.Sprintf("cpe:"+e.reason, e.attr...)
}
