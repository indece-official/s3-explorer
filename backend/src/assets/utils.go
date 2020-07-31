package assets

import "io/ioutil"

// ReadFile reads a asset file as string
func ReadFile(filename string) (string, error) {
	file, err := Assets.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ReadFileBinary reads a asset file as []byte
func ReadFileBinary(filename string) ([]byte, error) {
	file, err := Assets.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
