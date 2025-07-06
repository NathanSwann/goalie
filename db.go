package main

import (
	"encoding/binary"
	"math"
	"os"
	"strings"

	"codeberg.org/anaseto/goal"
)

func vfDbSave(ctx *goal.Context, args []goal.V) goal.V {
	path, ok := args[1].BV().(goal.S)
	if !ok {
		return goal.Panicf("db.save: Invalid path to save type: %s", args[1].Type())
	}
	dict, ok := args[0].BV().(*goal.D)
	if !ok {
		return goal.Panicf("db.save: Invalid data to save type: %s", args[0].Type())
	}
	_ = dict
	keys := dict.KeyArray()
	i_count := keys.Len()
	schema := ""
	for i := 0; i < i_count; i++ {
		key, ok := keys.At(i).BV().(goal.S)
		if !ok {
			return goal.Panicf("db.save: Invalid key type: %s", keys.At(i).Type())
		}
		data, ok := dict.GetS(string(key))
		if !ok {
			return goal.Panicf("db.save: Unable to find data for key: %s", string(key))
		}
		dt := data.Type()
		println(dt)
		switch dt {
		case "S":
			sd, ok := data.BV().(*goal.AS)
			if !ok {
				return goal.Panicf("db.save: Unable to convert type key: %s", string(key))
			}
			writeString(string(path), string(key), sd.Slice)
		case "N":
			sd, ok := data.BV().(*goal.AF)
			if !ok {
				return goal.Panicf("db.save: Unable to convert type key: %s", string(key))
			}
			writeFloat(string(path), string(key), sd.Slice)
		case "I":
			sd, ok := data.BV().(*goal.AI)
			if !ok {
				id, ok := data.BV().(*goal.AF)

				if !ok {
					return goal.Panicf("db.save: Unable to convert type key: %s %s", string(key), data.BV().Type())
				}
				s := id.Slice
				ns := make([]int64, len(s))
				for i, s := range s {
					ns[i] = int64(s)
				}
				writeInt(string(path), string(key), ns)

			} else {
				writeInt(string(path), string(key), sd.Slice)
			}
		default:
			return goal.Panicf("db.save: Invalid type: %s ", dt)
		}
		schema += string(key) + "\t" + string(dt) + "\n"
	}
	os.WriteFile(string(path)+".schema", []byte(schema), 0644)
	return goal.NewI(1)
}

func writeString(path string, key string, data []string) {
	os.MkdirAll(path, os.ModePerm)
	count_path := path + "/" + key + "#"
	data_path := path + "/" + key
	i_count := len(data)
	count_bytes := make([]byte, 8*i_count)
	data_size := 0
	for i := 0; i < i_count; i++ {
		s := data[i]
		binary.LittleEndian.PutUint64(count_bytes[i*8:], uint64(len(s)))
		data_size += len(s)
	}
	buf_l := 0
	data_bytes := make([]byte, data_size)
	for i := 0; i < i_count; i++ {
		s := data[i]
		copy(data_bytes[buf_l:], []byte(s))
		buf_l += len(s)
	}

	err := os.WriteFile(count_path, count_bytes, 0644)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(data_path, data_bytes, 0644)
	if err != nil {
		panic(err)
	}
	_ = data_path
}

func writeFloat(path string, key string, data []float64) {
	os.MkdirAll(path, os.ModePerm)
	data_path := path + "/" + key
	err := os.WriteFile(data_path, encodeFloats(data), 0644)
	if err != nil {
		panic(err)
	}
}

func writeInt(path string, key string, data []int64) {
	os.MkdirAll(path, os.ModePerm)
	data_path := path + "/" + key
	err := os.WriteFile(data_path, encodeInts(data), 0644)
	if err != nil {
		panic(err)
	}
}

func encodeFloats(fs []float64) []byte {
	buf := make([]byte, len(fs)*8)
	for i, f := range fs {
		u := math.Float64bits(f)
		binary.LittleEndian.PutUint64(buf[i*8:], u)
	}
	return buf
}

func encodeInts(fs []int64) []byte {
	buf := make([]byte, len(fs)*8)
	for i, f := range fs {
		binary.LittleEndian.PutUint64(buf[i*8:], uint64(f))
	}
	return buf
}

func decodeInts(fs []byte) []int64 {
	buf := make([]int64, len(fs)/8)
	for i := range buf {
		bits := binary.LittleEndian.Uint64(fs[i*8:])
		buf[i] = int64(bits)
	}
	return buf
}

func decodeFloats(fs []byte) []float64 {
	buf := make([]float64, len(fs)/8)
	for i := range buf {
		bits := binary.LittleEndian.Uint64(fs[i*8:])
		buf[i] = math.Float64frombits(bits)
	}
	return buf
}

func vfDbGet(ctx *goal.Context, args []goal.V) goal.V {
	path, ok := args[0].BV().(goal.S)
	if !ok {
		return goal.Panicf("db.get: Invalid path to save type: %s", args[1].Type())
	}
	data, err := os.ReadFile(string(path) + ".schema")
	if err != nil {
		return goal.Panicf("db.get: Failed to read schema: %s", err)
	}
	column_meta := strings.Split(string(data), "\n")
	cols := make([]goal.V, 0)
	table_data := make([]goal.V, 0)
	for i := 0; i < len(column_meta)-1; i++ {
		mi := strings.Split(column_meta[i], "\t")
		c, m := mi[0], mi[1]
		cols = append(cols, goal.NewS(c))
		switch m {
		case "S":
			table_data = append(table_data, goal.NewAS(readString(string(path), c)))
		case "N":
			table_data = append(table_data, goal.NewAF(readFloats(string(path), c)))
		case "I":
			table_data = append(table_data, goal.NewAI(readInts(string(path), c)))
		}
	}
	return goal.NewDict(goal.NewArray(cols), goal.NewArray(table_data))
}

func readFloats(path string, key string) []float64 {
	data_path := path + "/" + key
	data, err := os.ReadFile(data_path)
	if err != nil {
		panic(err)
	}
	return decodeFloats(data)
}

func readInts(path string, key string) []int64 {
	data_path := path + "/" + key
	data, err := os.ReadFile(data_path)
	if err != nil {
		panic(err)
	}
	return decodeInts(data)
}

func readString(path string, key string) []string {
	count_path := path + "/" + key + "#"
	data_path := path + "/" + key
	count_data, err := os.ReadFile(count_path)
	if err != nil {
		panic(err)
	}
	string_data, err := os.ReadFile(data_path)
	if err != nil {
		panic(err)
	}
	to_return := []string{}
	buf_loc := uint64(0)
	for i := 0; i < len(count_data)/8; i++ {
		s := binary.LittleEndian.Uint64(count_data[i*8 : (i+1)*8])
		to_return = append(to_return, string(string_data[buf_loc:buf_loc+s]))
		buf_loc += s
	}

	return to_return
}
