package util

import (
	"context"
	"fmt"
	"math/rand/v2"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lcnascimento/go-kit/o11y/baggage"
)

// CanonicalID is a UUID v7 (time-ordered) with a type prefix. The prefix
// is API/log presentation only — storage must persist the raw 16-byte UUID.
type CanonicalID struct {
	prefix string
	uuid   uuid.UUID
}

func NewCanonicalID(prefix string) CanonicalID {
	return CanonicalID{
		prefix: prefix,
		uuid:   uuid.Must(uuid.NewV7()),
	}
}

func (id CanonicalID) String() string {
	return fmt.Sprintf("%s_%s", id.prefix, id.uuid.String())
}

func (id CanonicalID) UUID() uuid.UUID {
	return id.uuid
}

func (id *CanonicalID) IsZero() bool {
	return id.uuid == uuid.Nil
}

func (id *CanonicalID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id *CanonicalID) UnmarshalText(b []byte) error {
	parsed, err := ParseCanonicalID(id.prefix, string(b))
	if err != nil {
		return err
	}

	*id = parsed

	return nil
}

func CanonicalIDFromUUID(prefix string, u uuid.UUID) CanonicalID {
	return CanonicalID{
		prefix: prefix,
		uuid:   u,
	}
}

func ParseCanonicalID(prefix, s string) (CanonicalID, error) {
	raw, ok := strings.CutPrefix(s, prefix+"_")
	if !ok {
		return CanonicalID{}, fmt.Errorf("domain: id %q must have prefix %q", s, prefix+"_")
	}

	u, err := uuid.Parse(raw)
	if err != nil {
		return CanonicalID{}, fmt.Errorf("domain: parsing id %q: %w", s, err)
	}

	return CanonicalID{prefix: prefix, uuid: u}, nil
}

func MustParseCanonicalID(prefix, s string) CanonicalID {
	id, err := ParseCanonicalID(prefix, s)
	if err != nil {
		panic(err)
	}

	return id
}

// Pick randomly selects one item of the given list.
func Pick[T any](items ...T) T {
	return items[rand.IntN(len(items))]
}

// RandomEmail generate a random email.
func RandomEmail() string {
	return fmt.Sprintf("user_%08x@example.com", rand.Uint32())
}

// RandomPhone generate a random Brazilian phone.
func RandomPhone() string {
	return fmt.Sprintf("+55119%08d", rand.IntN(100_000_000))
}

// RandomPastTime generate a random past Time.
func RandomPastTime() time.Time {
	offset := time.Duration(rand.IntN(30*24*60)) * time.Minute

	return time.Now().UTC().Add(-offset).Truncate(time.Millisecond)
}

// ToPointer returns a pointer reference to the given object.
func ToPointer[T any](v T) *T {
	return &v
}

// SafeValue returns the value associated to a pointer.
// If nil, returns the zero value of the given type.
func SafeValue[T any](v *T) T {
	if v == nil {
		return *new(T)
	}

	return *v
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)
}

func SortedKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
	})

	return keys
}

func CorrelationID(ctx context.Context) string {
	cID, ok := ctx.Value(baggage.MemberKeyCorrelationID).(string)
	if ok {
		return cID
	}

	return uuid.New().String()
}
