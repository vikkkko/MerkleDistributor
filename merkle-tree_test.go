package MerkleDistributor

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestConbinedHash(t *testing.T){
	hexString0 := "bcf9f204a16e397489a5fcaf7ee8d4b514e5a18b06b3d7ebc757fce96140d5e1"
	hexString1 := "7840b1d90ad73b24be171c1762b4b132bfba21c0f140c595d7c2292fd86b6102"
	b0,err := hex.DecodeString(hexString0)
	if err != nil {
		fmt.Println(err)
	}
	b1,err := hex.DecodeString(hexString1)
	if err != nil {
		fmt.Println(err)
	}
	bNew := ConbinedHash(bytes.NewBuffer(b0),bytes.NewBuffer(b1))
	fmt.Println(hex.EncodeToString(bNew.Bytes()))
	if hex.EncodeToString(bNew.Bytes()) != "7fe5c004d84699ba877988f4469deaa17d362b48c088559fabe9acf30d192396" {
		t.Error("err")
	}
}

func TestBuildMerkleTree(t *testing.T){
	/*
	0xbcf9f204a16e397489a5fcaf7ee8d4b514e5a18b06b3d7ebc757fce96140d5e1
	0x7840b1d90ad73b24be171c1762b4b132bfba21c0f140c595d7c2292fd86b6102
	0xf7dc479044fc49dbd6593936abd447a1d9aff662cbe516e521193a2027566c77
	0xe2609119fe100d0b0939e295615aaf82d4f75277eeae7a475f64a5101457e34c
	0x4610ac227c85a029b3d03fbaffa5d86afe68367d00ef28bdb0edf08de4fed311
	 */
	bts := make(Buffers,5,5)
	hexstrings := [5]string{
		"bcf9f204a16e397489a5fcaf7ee8d4b514e5a18b06b3d7ebc757fce96140d5e1",
		"7840b1d90ad73b24be171c1762b4b132bfba21c0f140c595d7c2292fd86b6102",
		"f7dc479044fc49dbd6593936abd447a1d9aff662cbe516e521193a2027566c77",
		"e2609119fe100d0b0939e295615aaf82d4f75277eeae7a475f64a5101457e34c",
		"4610ac227c85a029b3d03fbaffa5d86afe68367d00ef28bdb0edf08de4fed311",
	}
	for i,data := range hexstrings{
		b , _ := hex.DecodeString(data)
		bts[i] = *bytes.NewBuffer(b)
	}
	var mt MerkleTree
	mt.BuildMerkleTree(&bts)
	_,data := mt.GetHexRoot()
	if data != "307b70792f29e83c416e9b91610892d12adf8ef1f5e8b0689a55cc3e7e23f91f" {
		t.Error("err")
	}
}