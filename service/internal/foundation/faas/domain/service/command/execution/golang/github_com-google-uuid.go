// Code generated by 'yaegi extract github.com/google/uuid'. DO NOT EDIT.

package golang

import (
	"github.com/google/uuid"
	"reflect"
)

func init() {
	Symbols["github.com/google/uuid/uuid"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"ClockSequence":       reflect.ValueOf(uuid.ClockSequence),
		"FromBytes":           reflect.ValueOf(uuid.FromBytes),
		"Future":              reflect.ValueOf(uuid.Future),
		"GetTime":             reflect.ValueOf(uuid.GetTime),
		"Group":               reflect.ValueOf(uuid.Group),
		"Invalid":             reflect.ValueOf(uuid.Invalid),
		"Microsoft":           reflect.ValueOf(uuid.Microsoft),
		"Must":                reflect.ValueOf(uuid.Must),
		"MustParse":           reflect.ValueOf(uuid.MustParse),
		"NameSpaceDNS":        reflect.ValueOf(&uuid.NameSpaceDNS).Elem(),
		"NameSpaceOID":        reflect.ValueOf(&uuid.NameSpaceOID).Elem(),
		"NameSpaceURL":        reflect.ValueOf(&uuid.NameSpaceURL).Elem(),
		"NameSpaceX500":       reflect.ValueOf(&uuid.NameSpaceX500).Elem(),
		"New":                 reflect.ValueOf(uuid.New),
		"NewDCEGroup":         reflect.ValueOf(uuid.NewDCEGroup),
		"NewDCEPerson":        reflect.ValueOf(uuid.NewDCEPerson),
		"NewDCESecurity":      reflect.ValueOf(uuid.NewDCESecurity),
		"NewHash":             reflect.ValueOf(uuid.NewHash),
		"NewMD5":              reflect.ValueOf(uuid.NewMD5),
		"NewRandom":           reflect.ValueOf(uuid.NewRandom),
		"NewRandomFromReader": reflect.ValueOf(uuid.NewRandomFromReader),
		"NewSHA1":             reflect.ValueOf(uuid.NewSHA1),
		"NewString":           reflect.ValueOf(uuid.NewString),
		"NewUUID":             reflect.ValueOf(uuid.NewUUID),
		"Nil":                 reflect.ValueOf(&uuid.Nil).Elem(),
		"NodeID":              reflect.ValueOf(uuid.NodeID),
		"NodeInterface":       reflect.ValueOf(uuid.NodeInterface),
		"Org":                 reflect.ValueOf(uuid.Org),
		"Parse":               reflect.ValueOf(uuid.Parse),
		"ParseBytes":          reflect.ValueOf(uuid.ParseBytes),
		"Person":              reflect.ValueOf(uuid.Person),
		"RFC4122":             reflect.ValueOf(uuid.RFC4122),
		"Reserved":            reflect.ValueOf(uuid.Reserved),
		"SetClockSequence":    reflect.ValueOf(uuid.SetClockSequence),
		"SetNodeID":           reflect.ValueOf(uuid.SetNodeID),
		"SetNodeInterface":    reflect.ValueOf(uuid.SetNodeInterface),
		"SetRand":             reflect.ValueOf(uuid.SetRand),

		// type definitions
		"Domain":  reflect.ValueOf((*uuid.Domain)(nil)),
		"Time":    reflect.ValueOf((*uuid.Time)(nil)),
		"UUID":    reflect.ValueOf((*uuid.UUID)(nil)),
		"Variant": reflect.ValueOf((*uuid.Variant)(nil)),
		"Version": reflect.ValueOf((*uuid.Version)(nil)),
	}
}
