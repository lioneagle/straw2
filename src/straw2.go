package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"math"
	//"math/rand"
)

type Item struct {
	id     int32
	weight int32
}

type Bucket struct {
	weight int32
	items  []Item
}

func NewBucket() *Bucket {
	return &Bucket{items: make([]Item, 0)}
}

func (bucket *Bucket) AddItem(id, weight int32) {
	bucket.weight += weight
	bucket.items = append(bucket.items, Item{id: id, weight: weight})
}

func Hash(x int32) uint32 {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(x))
	return crc32.ChecksumIEEE(data)
}

func (bucket *Bucket) Select(x int32) int32 {
	max_item_id := int32(0)
	max_draw := -math.MaxFloat64
	for id, item := range bucket.items {
		draw := -math.MaxFloat64
		if item.weight != 0 {
			h := Hash(x * int32(id+100))
			//fmt.Println("h =", h)
			h &= 0xffff
			//fmt.Println("h2 =", h)
			draw = math.Log(float64(h)/65536.0) / float64(item.weight)
		}

		/*
			fmt.Println("draw =", draw)
			fmt.Println("max_draw =", max_draw)
			fmt.Println("draw > max_draw =", draw > max_draw) //*/

		if draw > max_draw {
			max_item_id = item.id
			max_draw = draw
		}
	}
	return max_item_id
}

type Node struct {
	id     int32
	weight int32
	data   []int32
}

type Nodes struct {
	weight int32
	nodes  []Node
}

func NewNodes() *Nodes {
	nodes := &Nodes{}
	nodes.nodes = make([]Node, 0)
	return nodes
}

func (nodes *Nodes) AddNode(id, weight int32) {
	nodes.weight += weight
	nodes.nodes = append(nodes.nodes, Node{id: id, weight: weight, data: make([]int32, 0)})
}

func (nodes *Nodes) AddNodeData(id, data int32) {
	nodes.nodes[id-1].data = append(nodes.nodes[id-1].data, data)
}

func (nodes *Nodes) ChangeWeight(id, weight int32) {
	delta := weight - nodes.nodes[id-1].weight
	nodes.weight += delta
	nodes.nodes[id-1].weight += delta
}

func (nodes *Nodes) String() string {
	str := fmt.Sprintf("weight = %d\n", nodes.weight)

	for _, node := range nodes.nodes {
		str += fmt.Sprintf("[%d]: weight = %d, counts = %d\n", node.id, node.weight, len(node.data))
	}
	return str
}

func BuildNodes() *Nodes {
	nodes := NewNodes()

	for i := 0; i < 10; i++ {
		nodes.AddNode(int32(i+1), 2)
	}

	//nodes.ChangeWeight(3, 6)
	//nodes.ChangeWeight(7, 3)

	return nodes
}

func BuildBucket(nodes *Nodes) *Bucket {
	bucket := NewBucket()

	for _, node := range nodes.nodes {
		bucket.AddItem(node.id, node.weight)
	}
	return bucket
}

func Run(nodes *Nodes, bucket *Bucket) {
	for i := 0; i < 50000; i++ {
		id := bucket.Select(int32(i))
		nodes.AddNodeData(id, int32(i))
	}
}

func main() {
	nodes := BuildNodes()
	bucket := BuildBucket(nodes)

	Run(nodes, bucket)

	fmt.Printf("nodes = %s\n", nodes)
}
