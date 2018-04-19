package ice

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/nkbai/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/nat/goice/sdp"
)

func loadData(tb testing.TB, name string) []byte {
	name = filepath.Join("testdata", name)
	f, err := os.Open(name)
	if err != nil {
		tb.Fatal(err)
	}
	defer func() {
		if errClose := f.Close(); errClose != nil {
			tb.Fatal(errClose)
		}
	}()
	v, err := ioutil.ReadAll(f)
	if err != nil {
		tb.Fatal(err)
	}
	return v
}

func TestConnectionAddress(t *testing.T) {
	data := loadData(t, "candidates_ex1.sdp")
	s, err := sdp.DecodeSession(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range s {
		p := candidateParser{
			c:   new(Candidate),
			buf: c.Value,
		}
		if err = p.parse(); err != nil {
			t.Fatal(err)
		}
		log.Trace("c= %s", log.StringInterface(p.c, 3))
	}

	// a=candidate:3862931549 1 udp 2113937151 192.168.220.128 56032
	//     foundation ---┘    |  |      |            |          |
	//   component id --------┘  |      |            |          |
	//      transport -----------┘      |            |          |
	//       priority ------------------┘            |          |
	//  conn. address -------------------------------┘          |
	//           port ------------------------------------------┘
}

func TestParse(t *testing.T) {
	//data := loadData(t, "candidates_ex1.sdp")
	//s, err := sdp.DecodeSession(data, nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//expected := []Candidate{}
	tCases := []struct {
		input    []byte
		expected Candidate
	}{
		{
			input: []byte("candidate:3862931549 1 udp 2113937151 192.168.220.128 56032 typ host generation 0 network-cost 50  alpha   beta    ??"),
			expected: Candidate{
				Foundation:  3862931549,
				ComponentID: 1,
				Priority:    2113937151,
				addr:        "192.168.220.128:56032",
				Type:        CandidateHost,
				transport:   TransportUDP,
				NetworkCost: 50,
				Attributes: Attributes{
					Attribute{
						Key:   []byte("alpha"),
						Value: []byte("beta"),
					},
				},
			}}, // 0
		{
			input: []byte("candidate:842163049 1 udp 1677729535 213.141.156.236 55726 typ srflx raddr"),
			expected: Candidate{
				Foundation:  842163049,
				ComponentID: 1,
				Priority:    1677729535,
				addr:        "213.141.156.236:55726",
				Type:        CandidateServerReflexive,
				transport:   TransportUDP,
			},
		}, {
			input: []byte("candidate:842163049 1 udp 1677729535 b2.cydev.ru 56024 typ srflx raddr 10.1.22.220 rport 56024 generation 0 ufrag eM2ytqY8D5Q07RAn"),
			expected: Candidate{
				Foundation:  842163049,
				ComponentID: 1,
				Priority:    1677729535,
				addr:        "b2.cydev.ru:56024",
				Type:        CandidateServerReflexive,
				relatedAddr: "10.1.22.220:56024",
				Generation:  0,
				transport:   TransportUDP,
				Attributes: Attributes{
					Attribute{
						Key:   []byte("ufrag"),
						Value: []byte("eM2ytqY8D5Q07RAn"),
					},
				},
			},
		},
	}

	for i, c := range tCases {
		parser := candidateParser{
			buf: c.input,
			c:   new(Candidate),
		}
		if err := parser.parse(); err != nil {
			t.Errorf("[%d]: unexpected error %s",
				i, err,
			)
		}
		if !c.expected.Equal(parser.c) {
			t.Errorf("[%d]: %#v != %#v (exp)",
				i, parser.c, c.expected,
			)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	data := loadData(b, "candidates_ex1.sdp")
	s, err := sdp.DecodeSession(data, nil)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	value := s[0].Value
	p := candidateParser{
		c: new(Candidate),
	}
	for i := 0; i < b.N; i++ {
		p.buf = value
		if err = p.parse(); err != nil {
			b.Fatal(err)
		}
		p.c.reset()
	}
}
