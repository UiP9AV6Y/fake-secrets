package io

type InfiniteReader []byte

func (r InfiniteReader) Read(b []byte) (int, error) {
	m := len(r)

	for i := range b {
		j := i % m

		b[i] = r[j]
	}

	return len(b), nil
}
