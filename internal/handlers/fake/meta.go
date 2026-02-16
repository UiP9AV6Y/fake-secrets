package fake

import (
	"crypto/x509/pkix"
	"io/fs"
	"net/http"
	"strconv"
	"time"
)

type StaticMeta struct {
	Method     string `json:"method"`
	RemoteAddr string `json:"remote_address"`
}

func NewStaticMeta(r *http.Request) *StaticMeta {
	result := &StaticMeta{
		Method:     r.Method,
		RemoteAddr: r.RemoteAddr,
	}

	return result
}

type FileMeta struct {
	StaticMeta `json:",inline"`

	Size     int64 `json:"size,omitempty"`
	Modified int64 `json:"modified,omitempty"`
}

func NewFileMeta(info fs.FileInfo, r *http.Request) *FileMeta {
	static := NewStaticMeta(r)
	result := &FileMeta{
		StaticMeta: *static,
		Size:       info.Size(),
		Modified:   info.ModTime().Unix(),
	}

	return result
}

type TokenMeta struct {
	StaticMeta `json:",inline"`

	Seed string `json:"seed,omitempty"`
}

func NewTokenMeta(seed string, r *http.Request) *TokenMeta {
	static := NewStaticMeta(r)
	result := &TokenMeta{
		StaticMeta: *static,
		Seed:       seed,
	}

	return result
}

type PasswordMeta struct {
	StaticMeta `json:",inline"`

	Length int `json:"length,omitempty"`

	Upper   bool `json:"upper"`
	Lower   bool `json:"lower"`
	Numeric bool `json:"numeric"`
	Special bool `json:"special"`
}

func ParsePasswordMeta(r *http.Request) (*PasswordMeta, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	length, err := ParseFormInt(r, "length", 12)
	if err != nil {
		return nil, err
	}

	upper, err := ParseFormBool(r, "upper", false)
	if err != nil {
		return nil, err
	}

	lower, err := ParseFormBool(r, "lower", false)
	if err != nil {
		return nil, err
	}

	numeric, err := ParseFormBool(r, "numeric", false)
	if err != nil {
		return nil, err
	}

	special, err := ParseFormBool(r, "special", false)
	if err != nil {
		return nil, err
	}

	if !upper && !lower && !numeric && !special {
		upper = true
		lower = true
		numeric = true
	}

	static := NewStaticMeta(r)
	result := &PasswordMeta{
		StaticMeta: *static,
		Length:     int(length),
		Upper:      upper,
		Lower:      lower,
		Numeric:    numeric,
		Special:    special,
	}

	return result, nil
}

type SSHMeta struct {
	StaticMeta `json:",inline"`

	Hostname string `json:"hostname,omitempty"`

	Length int `json:"length,omitempty"`
}

func ParseSSHMeta(hostname string, r *http.Request) (*SSHMeta, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	length, err := ParseFormInt(r, "length", 4096)
	if err != nil {
		return nil, err
	}

	static := NewStaticMeta(r)
	result := &SSHMeta{
		StaticMeta: *static,
		Hostname:   hostname,
		Length:     int(length),
	}

	return result, nil
}

type TLSMeta struct {
	StaticMeta `json:",inline"`

	Hostname     string `json:"hostname,omitempty"`
	Organization string `json:"organization,omitempty"`

	Length   int   `json:"length,omitempty"`
	ValidFor int64 `json:"valid_for,omitempty"`
	ValidAt  int64 `json:"valid_at,omitempty"`
}

func ParseTLSMeta(hostname string, start time.Time, r *http.Request) (*TLSMeta, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	length, err := ParseFormInt(r, "length", 2048)
	if err != nil {
		return nil, err
	}

	validAt, err := ParseFormInt(r, "valid_at", start.Unix())
	if err != nil {
		return nil, err
	}

	validFor, err := ParseFormInt(r, "valid_for", 60*60*24*90) // 90 days
	if err != nil {
		return nil, err
	}

	organization := ParseFormString(r, "organization", "Acme Co")
	static := NewStaticMeta(r)
	result := &TLSMeta{
		StaticMeta:   *static,
		Hostname:     hostname,
		Organization: organization,
		Length:       int(length),
		ValidFor:     validFor,
		ValidAt:      validAt,
	}

	return result, nil
}

func (m *TLSMeta) NotBefore() time.Time {
	return time.Unix(m.ValidAt, 0)
}

func (m *TLSMeta) NotAfter() time.Time {
	return time.Unix(m.ValidAt+m.ValidFor, 0)
}

func (m *TLSMeta) Subject() pkix.Name {
	result := pkix.Name{
		Organization: []string{m.Organization},
	}

	return result
}

func (m *TLSMeta) SubjectAltNames() []string {
	result := []string{
		m.Hostname,
	}

	return result
}

func ParseFormString(r *http.Request, field, fallback string) string {
	value := r.FormValue(field)
	if value == "" {
		return fallback
	}

	return value
}

func ParseFormBool(r *http.Request, field string, fallback bool) (bool, error) {
	value := r.FormValue(field)
	if value == "" {
		return fallback, nil
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		return fallback, err
	}

	return result, nil
}

func ParseFormInt(r *http.Request, field string, fallback int64) (int64, error) {
	value := r.FormValue(field)
	if value == "" {
		return fallback, nil
	}

	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil || result <= 0 {
		return fallback, err
	}

	return result, nil
}
