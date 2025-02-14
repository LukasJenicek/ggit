package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuild(t *testing.T) {
	type args struct {
		root    *Tree
		entries []*Entry
	}
	tests := []struct {
		name string
		args args
		want *Tree
	}{
		{
			name: "Build recursive structure",
			args: args{
				root: &Tree{},
				entries: []*Entry{
					{
						Name:       "hello.txt",
						oid:        []byte{0xde, 0xad, 0xbe, 0xef},
						executable: false,
					},
					{
						Name:       "libs/hello.txt",
						oid:        []byte{0xde, 0xad, 0xbe, 0xef},
						executable: false,
					},
					{
						Name:       "libs/internal/internal.txt",
						oid:        []byte{0xde, 0xad, 0xbe, 0xef},
						executable: false,
					},
					{
						Name:       "world.txt",
						oid:        []byte{0xde, 0xad, 0xbe, 0xef},
						executable: false,
					},
				},
			},
			want: &Tree{
				entries: []Object{
					&Entry{
						Name:       "hello.txt",
						oid:        []byte{0xde, 0xad, 0xbe, 0xef},
						executable: false,
					},
					&Tree{
						entries: []Object{
							&Entry{
								Name:       "libs/hello.txt",
								oid:        []byte{0xde, 0xad, 0xbe, 0xef},
								executable: false,
							},
							&Tree{
								entries: []Object{
									&Entry{
										Name:       "libs/internal/internal.txt",
										oid:        []byte{0xde, 0xad, 0xbe, 0xef},
										executable: false,
									},
								},
							},
						},
					},
					&Entry{
						Name:       "world.txt",
						oid:        []byte{0xde, 0xad, 0xbe, 0xef},
						executable: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			build := Build(tt.args.root, tt.args.entries)

			assert.EqualValues(t, tt.want, build)
		})
	}
}
