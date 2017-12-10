package output

type FileOutput struct {
	name string
}

func NewFileOutput(name string) (*FileOutput, error) {
	if path.Ext(name) != ".json" {
		return nil, errors.New("only support json")
	}
	return &FileOutput{name: name}, nil
}

func (f *FileOutput) Put(data []byte) error {
	if !govalidator.IsJSON(string(data)) {
		return errors.New("illegal json")
	}

	return ioutil.WriteFile(f.name, data, 0660)
}
