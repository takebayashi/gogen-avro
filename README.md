gogen-avro
===

[![Build Status](https://travis-ci.org/alanctgardner/gogen-avro.svg?branch=master)](https://travis-ci.org/alanctgardner/gogen-avro)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/alanctgardner/gogen-avro/master/LICENSE)
[![Version 2.0.0](https://img.shields.io/badge/version-2.0.0-lightgrey.svg)](https://gopkg.in/alanctgardner/gogen-avro.v2)

Generate Go structures and serializer / deserializer methods from Avro schemas. Generated serializers/deserializers are 2-8x faster than goavro, and you get compile-time safety for getting and setting fields.

### Installation

gogen-avro is a tool which you install on your system (usually on your GOPATH), and run as part of your build process. To install gogen-avro to `$GOPATH/bin/`, run:

```
go get gopkg.in/alanctgardner/gogen-avro.v2/...
```

### Usage

To generate Go source files from one or more Avro schema files, run:

```
gogen-avro <output directory> <avro schema files>
```

Or use a `go:generate` directive in a source file ([example](https://github.com/alanctgardner/gogen-avro/blob/master/test/primitive/schema_test.go)):

```
//go:generate $GOPATH/bin/gogen-avro . primitives.avsc
```

The generated source files contain structs for each schema, plus a function `Serialize(io.Writer)` to encode the contents into the given `io.Writer`, and `Deserialize<RecordType>(io.Reader)` to read a struct from the given `io.Reader`. 

### Example

The `example` directory contains a simple example project with an Avro schema. Once you've installed gogen-avro on your GOPATH, you can install the example project:

```
# Build the Go source files from the Avro schema using the generate directive 
go generate github.com/alanctgardner/gogen-avro/example

# Install the example project on the gopath
go install github.com/alanctgardner/gogen-avro/example
```

### Type Conversion

Gogen-avro produces a Go struct which reflects the structure of your Avro schema. Most Go types map neatly onto Avro types:

| Avro Type     | Go Type           | Notes                                                                                                                |
|---------------|-------------------|----------------------------------------------------------------------------------------------------------------------|
| null          | interface{}       | This is just a placeholder, nothing is encoded/decoded                                                               |
| boolean       | bool              |                                                                                                                      |
| int, long     | int32,int64       |                                                                                                                      |
| float, double | float32, float64  |                                                                                                                      |
| bytes         | []byte            |                                                                                                                      |
| string        | string            |                                                                                                                      |
| enum          | custom type       | Generates a type with a constant for each symbol                                                                     |
| array<type>   | []<type>          |                                                                                                                      |
| map<type>     | map[string]<type> |                                                                                                                      |
| fixed         | [<n>]byte         | Fixed fields are given a custom type, which is an alias for an appropriately sized byte array                        |
| union         | custom type       | Unions are handled as a struct with one field per possible type, and an enum field to dictate which field to read    |

`union` is more complicated than primitive types. We generate a struct and enum whose name is uniquely determined by the types in the union. For a field whose type is `["null", "int"]` we generate the following:

```
type UnionNullInt struct {
	// All the possible types the union could take on
	Null               interface{}
	Int                int32
	// Which field actually has data in it - defaults to the first type in the list, "null"
	UnionType          UnionNullTypeEnum
}

type UnionNullIntTypeEnum int

const (
	UnionNullIntTypeEnumNull            UnionNullIntTypeEnum = 0
	UnionNullIntTypeEnumInt             UnionNullIntTypeEnum = 1
)
``` 

### TODO / Caveats

This package doesn't implement the entire Avro 1.7.7 specification, specifically:

- Schema resolution
- Framing - generate RPCs and container format readers/writers

### Changelog

2.0
---
- Bug fixes for arrays and maps with record members
- Refactored internals significantly

1.0
---
- Initial release

### Thanks

Thanks to LinkedIn's [goavro library](https://github.com/linkedin/goavro), for providing the encoders for primitives.
