package fake

import (
	"crypto/x509/pkix"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	nethttp "net/http"
	"time"

	"github.com/UiP9AV6Y/fake-secrets/internal/crypto"
	"github.com/UiP9AV6Y/fake-secrets/internal/http"
)

type StaticMeta struct {
	Method     string `json:"method"`
	RemoteAddr string `json:"remote_address"`
}

func NewStaticMeta(r *nethttp.Request) *StaticMeta {
	result := &StaticMeta{
		Method:     r.Method,
		RemoteAddr: r.RemoteAddr,
	}

	return result
}

func (m *StaticMeta) LogValue() slog.Value {
	return slog.GroupValue(m.LogAttrs()...)
}

func (m *StaticMeta) LogAttrs() []slog.Attr {
	attrs := []slog.Attr{
		slog.String("method", m.Method),
		slog.String("remote_address", m.RemoteAddr),
	}

	return attrs
}

func (m *StaticMeta) String() string {
	return DescribeStruct(m, "StaticMeta")
}

func (m *StaticMeta) StructWriteTo(w io.Writer) (int, error) {
	_, _ = fmt.Fprintf(w, ", method=%s", m.Method)
	_, _ = fmt.Fprintf(w, ", remote_address=%s", m.RemoteAddr)

	return 0, nil
}

type FileMeta struct {
	StaticMeta `json:",inline"`

	Size     int64 `json:"size,omitempty"`
	Modified int64 `json:"modified,omitempty"`
}

func NewFileMeta(info fs.FileInfo, r *nethttp.Request) *FileMeta {
	static := NewStaticMeta(r)
	result := &FileMeta{
		StaticMeta: *static,
		Size:       info.Size(),
		Modified:   info.ModTime().Unix(),
	}

	return result
}

func (m *FileMeta) LogValue() slog.Value {
	return slog.GroupValue(m.LogAttrs()...)
}

func (m *FileMeta) LogAttrs() []slog.Attr {
	attrs := []slog.Attr{
		slog.Int64("size", m.Size),
		slog.Int64("modified", m.Modified),
	}

	return append(attrs, m.StaticMeta.LogAttrs()...)
}

func (m *FileMeta) String() string {
	return DescribeStruct(m, "FileMeta")
}

func (m *FileMeta) StructWriteTo(w io.Writer) (int, error) {
	_, _ = m.StaticMeta.StructWriteTo(w)
	_, _ = fmt.Fprintf(w, ", size=%d", m.Size)
	_, _ = fmt.Fprintf(w, ", modified=%d", m.Modified)

	return 0, nil
}

type TokenMeta struct {
	StaticMeta `json:",inline"`

	Seed string `json:"seed,omitempty"`
}

func NewTokenMeta(seed string, r *nethttp.Request) *TokenMeta {
	static := NewStaticMeta(r)
	result := &TokenMeta{
		StaticMeta: *static,
		Seed:       seed,
	}

	return result
}

func (m *TokenMeta) LogValue() slog.Value {
	return slog.GroupValue(m.LogAttrs()...)
}

func (m *TokenMeta) LogAttrs() []slog.Attr {
	attrs := []slog.Attr{
		slog.String("seed", m.Seed),
	}

	return append(attrs, m.StaticMeta.LogAttrs()...)
}

func (m *TokenMeta) String() string {
	return DescribeStruct(m, "TokenMeta")
}

func (m *TokenMeta) StructWriteTo(w io.Writer) (int, error) {
	_, _ = m.StaticMeta.StructWriteTo(w)
	_, _ = fmt.Fprintf(w, ", seed=%s", m.Seed)

	return 0, nil
}

type PasswordMeta struct {
	StaticMeta `json:",inline"`

	Length int `json:"length,omitempty"`

	Upper   bool `json:"upper"`
	Lower   bool `json:"lower"`
	Numeric bool `json:"numeric"`
	Special bool `json:"special"`
}

func ParsePasswordMeta(r *nethttp.Request) (*PasswordMeta, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	length, err := http.ParseFormInt(r, "length", 12)
	if err != nil {
		return nil, err
	}

	upper, err := http.ParseFormBool(r, "upper", false)
	if err != nil {
		return nil, err
	}

	lower, err := http.ParseFormBool(r, "lower", false)
	if err != nil {
		return nil, err
	}

	numeric, err := http.ParseFormBool(r, "numeric", false)
	if err != nil {
		return nil, err
	}

	special, err := http.ParseFormBool(r, "special", false)
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

func (m *PasswordMeta) LogValue() slog.Value {
	return slog.GroupValue(m.LogAttrs()...)
}

func (m *PasswordMeta) LogAttrs() []slog.Attr {
	attrs := []slog.Attr{
		slog.Int("length", m.Length),
		slog.Bool("upper", m.Upper),
		slog.Bool("lower", m.Lower),
		slog.Bool("numeric", m.Numeric),
		slog.Bool("special", m.Special),
	}

	return append(attrs, m.StaticMeta.LogAttrs()...)
}

func (m *PasswordMeta) String() string {
	return DescribeStruct(m, "PasswordMeta")
}

func (m *PasswordMeta) StructWriteTo(w io.Writer) (int, error) {
	_, _ = m.StaticMeta.StructWriteTo(w)
	_, _ = fmt.Fprintf(w, ", length=%d", m.Length)
	_, _ = fmt.Fprintf(w, ", upper=%t", m.Upper)
	_, _ = fmt.Fprintf(w, ", lower=%t", m.Lower)
	_, _ = fmt.Fprintf(w, ", numeric=%t", m.Numeric)
	_, _ = fmt.Fprintf(w, ", special=%t", m.Special)

	return 0, nil
}

type CryptoMeta struct {
	Subject string `json:"subject,omitempty"`

	Length     int               `json:"length,omitempty"`
	Algorithm  crypto.Algorithm  `json:"algorithm,omitempty"`
	ECDSACurve crypto.ECDSACurve `json:"ecdsa_curve,omitempty"`
}

func ParseCryptoMeta(subject string, r *nethttp.Request) (*CryptoMeta, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	length, err := http.ParseFormInt(r, "length", 4096)
	if err != nil {
		return nil, err
	}

	algo, err := http.ParseFormCryptoAlgorithm(r, "algorithm", crypto.AlgorithmRSA)
	if err != nil {
		return nil, err
	}

	curve, err := http.ParseFormECDSACurve(r, "curve", crypto.ECDSACurveP256)
	if err != nil {
		return nil, err
	}

	result := &CryptoMeta{
		Subject:    subject,
		Length:     int(length),
		Algorithm:  algo,
		ECDSACurve: curve,
	}

	return result, nil
}

func (m *CryptoMeta) LogValue() slog.Value {
	return slog.GroupValue(m.LogAttrs()...)
}

func (m *CryptoMeta) LogAttrs() []slog.Attr {
	attrs := []slog.Attr{
		slog.String("subject", m.Subject),
		slog.Int("length", m.Length),
		slog.Any("algorithm", m.Algorithm),
		slog.Any("ecdsa_curve", m.ECDSACurve),
	}

	return attrs
}

func (m *CryptoMeta) String() string {
	return DescribeStruct(m, "CryptoMeta")
}

func (m *CryptoMeta) StructWriteTo(w io.Writer) (int, error) {
	_, _ = fmt.Fprintf(w, ", subject=%s", m.Subject)
	_, _ = fmt.Fprintf(w, ", length=%d", m.Length)
	_, _ = fmt.Fprintf(w, ", algorithm=%s", m.Algorithm)
	_, _ = fmt.Fprintf(w, ", ecdsa_curve=%s", m.ECDSACurve)

	return 0, nil
}

type SSHMeta struct {
	StaticMeta `json:",inline"`
	CryptoMeta `json:",inline"`
}

func ParseSSHMeta(hostname string, r *nethttp.Request) (*SSHMeta, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	crypt, err := ParseCryptoMeta(hostname, r)
	if err != nil {
		return nil, err
	}

	static := NewStaticMeta(r)
	result := &SSHMeta{
		StaticMeta: *static,
		CryptoMeta: *crypt,
	}

	return result, nil
}

func (m *SSHMeta) LogValue() slog.Value {
	return slog.GroupValue(m.LogAttrs()...)
}

func (m *SSHMeta) LogAttrs() []slog.Attr {
	attrs := append([]slog.Attr{}, m.StaticMeta.LogAttrs()...)

	return append(attrs, m.CryptoMeta.LogAttrs()...)
}

func (m *SSHMeta) String() string {
	return DescribeStruct(m, "SSHMeta")
}

func (m *SSHMeta) StructWriteTo(w io.Writer) (int, error) {
	_, _ = m.StaticMeta.StructWriteTo(w)
	_, _ = m.CryptoMeta.StructWriteTo(w)

	return 0, nil
}

type TLSMeta struct {
	StaticMeta `json:",inline"`
	CryptoMeta `json:",inline"`

	Organization string `json:"organization,omitempty"`

	ValidFor int64 `json:"valid_for,omitempty"`
	ValidAt  int64 `json:"valid_at,omitempty"`
}

func ParseTLSMeta(hostname string, start time.Time, r *nethttp.Request) (*TLSMeta, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	crypt, err := ParseCryptoMeta(hostname, r)
	if err != nil {
		return nil, err
	}

	validAt, err := http.ParseFormInt(r, "valid_at", start.Unix())
	if err != nil {
		return nil, err
	}

	validFor, err := http.ParseFormInt(r, "valid_for", 60*60*24*90) // 90 days
	if err != nil {
		return nil, err
	}

	organization := http.ParseFormString(r, "organization", "Acme Co")
	static := NewStaticMeta(r)
	result := &TLSMeta{
		StaticMeta:   *static,
		CryptoMeta:   *crypt,
		Organization: organization,
		ValidFor:     validFor,
		ValidAt:      validAt,
	}

	return result, nil
}

func (m *TLSMeta) LogValue() slog.Value {
	return slog.GroupValue(m.LogAttrs()...)
}

func (m *TLSMeta) LogAttrs() []slog.Attr {
	attrs := []slog.Attr{
		slog.String("organization", m.Organization),
		slog.Int64("valid_for", m.ValidFor),
		slog.Int64("valid_at", m.ValidAt),
	}
	attrs = append(attrs, m.StaticMeta.LogAttrs()...)

	return append(attrs, m.CryptoMeta.LogAttrs()...)
}

func (m *TLSMeta) String() string {
	return DescribeStruct(m, "TLSMeta")
}

func (m *TLSMeta) StructWriteTo(w io.Writer) (int, error) {
	_, _ = m.StaticMeta.StructWriteTo(w)
	_, _ = m.CryptoMeta.StructWriteTo(w)
	_, _ = fmt.Fprintf(w, ", organization=%s", m.Organization)
	_, _ = fmt.Fprintf(w, ", valid_for=%d", m.ValidFor)
	_, _ = fmt.Fprintf(w, ", valid_at=%d", m.ValidAt)

	return 0, nil
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
		m.CryptoMeta.Subject,
	}

	return result
}
