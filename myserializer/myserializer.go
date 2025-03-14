package myserializer

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Marshal(v interface{}) ([]byte, error) {
	return marshalValue(reflect.ValueOf(v))
}

func Unmarshal(data []byte, v interface{}) error {
	return unmarshalValue(strings.TrimSpace(string(data)), reflect.ValueOf(v).Elem())
}

func marshalValue(val reflect.Value) ([]byte, error) {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return []byte("nil"), nil
		}
		return marshalValue(val.Elem())
	}

	var buf bytes.Buffer
	switch val.Kind() {
	case reflect.Struct:
		buf.WriteString("{")
		typeOf := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typeOf.Field(i)
			if field.PkgPath != "" {
				continue
			}
			fieldValue := val.Field(i)
			data, err := marshalValue(fieldValue)
			if err != nil {
				return nil, err
			}
			fmt.Fprintf(&buf, "%s %s ", field.Name, data)
		}
		buf.WriteString("}")

	case reflect.String:
		buf.WriteString(fmt.Sprintf("S\"%s\"", val.String()))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(fmt.Sprintf("N%d", val.Int()))

	case reflect.Float32, reflect.Float64:
		buf.WriteString(fmt.Sprintf("F%v", val.Float()))

	case reflect.Bool:
		if val.Bool() {
			buf.WriteString("B1")
		} else {
			buf.WriteString("B0")
		}

	case reflect.Slice:
		buf.WriteString(fmt.Sprintf("L%d[", val.Len()))
		for i := 0; i < val.Len(); i++ {
			data, err := marshalValue(val.Index(i))
			if err != nil {
				return nil, err
			}
			buf.Write(data)
			if i < val.Len()-1 {
				buf.WriteString(" ")
			}
		}
		buf.WriteString("]")

	default:
		return nil, fmt.Errorf("unsupported type: %s", val.Type())
	}
	return buf.Bytes(), nil
}

func unmarshalValue(data string, v reflect.Value) error {
	data = strings.TrimSpace(data)
	if data == "" {
		return errors.New("empty data")
	}

	if data == "nil" {
		if v.Kind() == reflect.Ptr {
			v.Set(reflect.Zero(v.Type()))
			return nil
		}
		return errors.New("cannot set nil to non-pointer type")
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return unmarshalValue(data, v.Elem())
	}

	if v.Kind() == reflect.Slice {
		if !strings.HasPrefix(data, "L") {
			return errors.New("invalid slice prefix")
		}

		parts := strings.SplitN(data[len("L"):], "[", 2)
		if len(parts) != 2 {
			return errors.New("invalid slice format")
		}

		lengthStr := strings.TrimSpace(parts[0])
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return err
		}

		content := strings.TrimSuffix(parts[1], "]")
		fields := tokenize(content)

		if len(fields) != length {
			return fmt.Errorf("slice length mismatch: expected %d, got %d", length, len(fields))
		}

		slice := reflect.MakeSlice(v.Type(), length, length)
		for i := 0; i < length; i++ {
			if err := unmarshalValue(fields[i], slice.Index(i)); err != nil {
				return err
			}
		}
		v.Set(slice)
		return nil
	}

	if strings.HasPrefix(data, "{") {
		return parseStruct(data, v)
	}

	return parsePrimitive(data, v)
}

func parseStruct(data string, v reflect.Value) error {
	content := strings.TrimSpace(data[1 : len(data)-1])
	fields := tokenize(content)

	if len(fields)%2 != 0 {
		return errors.New("invalid struct format")
	}

	fieldMap := make(map[string]string)
	for i := 0; i < len(fields); i += 2 {
		fieldMap[fields[i]] = fields[i+1]
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if value, ok := fieldMap[field.Name]; ok {
			fieldValue := v.Field(i)
			if err := unmarshalValue(value, fieldValue); err != nil {
				return err
			}
		}
	}
	return nil
}

func parsePrimitive(data string, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		if strings.HasPrefix(data, "S\"") {
			v.SetString(data[2 : len(data)-1])
		} else {
			v.SetString(data)
		}

	case reflect.Int, reflect.Int64:
		if strings.HasPrefix(data, "N") {
			n, _ := strconv.Atoi(data[1:])
			v.SetInt(int64(n))
		}

	case reflect.Float32, reflect.Float64:
		if strings.HasPrefix(data, "F") {
			f, _ := strconv.ParseFloat(data[1:], 64)
			v.SetFloat(f)
		}

	case reflect.Bool:
		v.SetBool(data == "B1")

	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func tokenize(data string) []string {
	var tokens []string
	var buf strings.Builder
	depth := 0
	inQuotes := false

	for _, c := range data {
		switch {
		case c == '"':
			inQuotes = !inQuotes
			buf.WriteRune(c)

		case c == ' ' && depth == 0 && !inQuotes:
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}

		case c == '{' || c == '[':
			depth++
			buf.WriteRune(c)

		case c == '}' || c == ']':
			depth--
			buf.WriteRune(c)

		default:
			buf.WriteRune(c)
		}
	}

	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}

	return tokens
}
