package main

func (f *FilesOps) SpoolOn() *FilesOps {
	if f.e != nil {
		return nil
	}

	f.spooling = true

	return f
}

func (f *FilesOps) SpoolOff() *FilesOps {
	if f.e != nil {
		return nil
	}

	f.spooling = false

	return f
}
