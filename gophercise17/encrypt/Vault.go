package encrypt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

// Mocked Functions
var DecReaderFunc = DecReader
var EncWriterFunc = EncWriter

type Vault struct {
	enckey   string
	filepath string
	mux      sync.Mutex
	mapping  map[string]string
}

func NewVault(ekey, fp string) *Vault {
	return &Vault{
		enckey: ekey,

		filepath: fp,
	}
}
func (v *Vault) LoadMapping() error {
	fptr, err := os.Open(v.filepath)
	if err != nil {
		v.mapping = make(map[string]string)
		return nil
	}

	reader, err := DecReaderFunc(v.enckey, fptr)
	if err != nil {
		return err
	}
	return v.readmapping(reader)

	/*	var sb strings.Builder
		defer fptr.Close()
		_, err2 := io.Copy(&sb, fptr)
		if err2 != nil {

			return err
		}
		decodedjson, err := Dec(v.enckey, sb.String())
		if err != nil {

			return err
		}
		r := strings.NewReader(decodedjson)
		dec := json.NewDecoder(r)
		err = dec.Decode(&v.mapping)
		if err != nil {

			return err
		}

			return nil
	*/
}
func (v *Vault) readmapping(r io.Reader) error {
	return json.NewDecoder(r).Decode(&v.mapping)
}

func (v *Vault) writemapping(w io.Writer) error {
	return json.NewEncoder(w).Encode(v.mapping)
}
func (v *Vault) savemapping() error {
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	writer, err := EncWriterFunc(v.enckey, f) //rgenerates cipher.stream based on key,returns struct streamwriter{writer,stream}  and passes filewriter to
	if err != nil {                           // to writer and cipher.stream to stream.
		return err // internally streamwriter uses XORKeyStream() on stream and passes it to writer on write
	} //				.|.
	return v.writemapping(writer) // cipher.StreamWriter writes to json.NewEncoder(<writer>) (interface chaining)
	/*var s strings.Builder
	encd := json.NewEncoder(&s)
	err := encd.Encode(v.mapping)
	if err != nil {
		return err
	}
	encodedjson, er := Enc(v.enckey, s.String())
	if er != nil {
		return er
	}

	_, er2 := fmt.Fprintf(f, encodedjson)
	if er2 != nil {
		return er2
	}
	return nil*/
}

func (v *Vault) Set(key, value string) error {
	// v.mux.Lock()
	// defer v.mux.Unlock()
	err := v.LoadMapping()
	if err != nil {
		return err
	}
	v.mapping[key] = value
	return v.savemapping()
}

func (v *Vault) Get(key string) (string, error) {
	// v.mux.Lock()
	// defer v.mux.Unlock()
	err := v.LoadMapping()
	if err != nil {
		fmt.Println(err.Error() + "1")
		return "", err
	}
	val, bol := v.mapping[key]
	if !bol {
		return "", errors.New("Value does not exists")
	}

	return val, nil
}
