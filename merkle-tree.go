package git

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"encoding/hex"
	"math"
	"sort"
)

type MerkleTree struct {
	layers [][]*Node
	layer0ElementPosition map[string] int
}

type Node struct {
	Hash bytes.Buffer
	Parent *Node
	left *Node
	right *Node
}

type Buffers []bytes.Buffer

func (Bs Buffers) Len() int          { return len(Bs) }
func (Bs Buffers) Less(i,j int) bool { return bytes.Compare(Bs[i].Bytes(), Bs[j].Bytes()) == -1 }
func (Bs Buffers) Swap(i,j int)      { Bs[i], Bs[j] = Bs[j], Bs[i]}

func(m *MerkleTree) BuildMerkleTree(contents *Buffers){
	sort.Sort(contents)
	m.layers = make([][]*Node,int64(math.Ceil(float64(contents.Len()) / 2) + 1))
	m.layer0ElementPosition = map[string]int{}
	m.layers[0]= contents.toNodes()
	for i,data := range m.layers[0]{
		m.layer0ElementPosition[data.Hash.String()] = i
	}
	for	i := 0 ; i < len(m.layers) - 1 ; i++ {
		layer := make([]*Node,int64(math.Ceil(float64(len(m.layers[i])) / 2)))
		for j := 0 ; j < len(m.layers[i]) ; j = j + 2 {
			if j+1<len(m.layers[i]) {
				//fmt.Printf("%v ---- %v \n",hex.EncodeToString(m.layers[i][j].Hash.Bytes()) ,hex.EncodeToString(m.layers[i][j+1].Hash.Bytes()))
				cHash := ConbinedHash(&m.layers[i][j].Hash,&m.layers[i][j+1].Hash)
				node := Node{
					Hash:   *cHash,
					Parent: nil,
					left:   m.layers[i][j],
					right:  m.layers[i][j+1],
				}
				layer[j/2] = &node
				m.layers[i][j].Parent = &node
				m.layers[i][j+1].Parent = &node
			} else {
				layer[j/2] = m.layers[i][j]
			}
		}
		m.layers[i+1] = layer
	}
}

func (m *MerkleTree) GetRootHash()(err error,hash bytes.Buffer){
	if (len(m.layers[len(m.layers) - 1])) != 1 {
		err = errors.New("invalidate tree")
	}
	hash = m.layers[len(m.layers)-1][0].Hash
	return
}

func (m *MerkleTree) GetHexRoot()(err error,hexHash string){
	if (len(m.layers[len(m.layers) - 1])) != 1 {
		err = errors.New("invalidate tree")
	}
	hexHash = hex.EncodeToString(m.layers[len(m.layers)-1][0].Hash.Bytes())
	return
}

func(m *MerkleTree) GetProof(n *Node)(ns []*Node){
	index := m.layer0ElementPosition[n.Hash.String()]
	ns = make([]*Node,0)
	for i:=0;i< len(m.layers)-1;i++{
		n := getPairElement(index,m.layers[i])
		if n != nil {
			ns = append(ns, n)
			index = int( index / 2)
		}
	}
	return
}

func getPairElement(index int,layer []*Node) (n *Node){
	var pairIndex int
	if index % 2 == 0 {
		pairIndex = index + 1
	} else {
		pairIndex = index - 1
	}
	if pairIndex >= len(layer) {
		n = nil
	} else {
		n = layer[pairIndex]
	}
	return
}

func(Bs Buffers) toNodes() (layer []*Node){
	layer = make([]*Node, Bs.Len(), Bs.Len())
	for i,data := range Bs {
		fmt.Println(hex.EncodeToString(data.Bytes()))
		node := Node{
			data,
			nil,
			nil,
			nil,
		}
		layer[i] = &node
	}
	return
}

func ConbinedHash(b0,b1 *bytes.Buffer) *bytes.Buffer {
	i := bytes.Compare(b0.Bytes(),b1.Bytes())
	var bts []byte
	if i == -1 { //a<b
		b0.Write(b1.Bytes())
		bts = b0.Bytes()
	} else {
		b1.Write(b0.Bytes())
		bts  = b1.Bytes()
	}
	h := sha3.NewLegacyKeccak256()
	h.Write(bts)
	return bytes.NewBuffer(h.Sum(nil))
}

