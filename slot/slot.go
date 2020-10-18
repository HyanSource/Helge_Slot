package slot

import (
	"math/rand"

	"github.com/HyanSource/Helge_Slot/pb"
)

//SymbolName
//0:"10"
//1:"J"
//2:"Q"
//3:"K"
//4:"A"
//5:"Parrot"
//6:"Skull"
//7:"Rum"
//8:"Treasure"
//9:"Captain"

//Table
//2  5  8  11 14
//1  4  7  10 13
//0  3  6  9  12

//遊玩的全域
var PlayGame *Slot

func init() {
	PlayGame = NewSlot()
}

//初始化
func NewSlot() *Slot {
	return &Slot{
		Odds: [][]int{
			[]int{0, 0, 5, 10, 20},
			[]int{0, 0, 5, 10, 20},
			[]int{0, 0, 5, 10, 25},
			[]int{0, 0, 5, 10, 25},
			[]int{0, 0, 5, 15, 25},
			[]int{0, 0, 10, 25, 100},
			[]int{0, 0, 10, 25, 125},
			[]int{0, 0, 15, 50, 150},
			[]int{0, 0, 15, 50, 175},
			[]int{0, 0, 25, 100, 500},
		},
		PayLines: [][]int{
			[]int{1, 4, 7, 10, 13},
			[]int{2, 5, 8, 11, 14},
			[]int{0, 3, 6, 9, 12},
			[]int{2, 4, 6, 10, 14},
			[]int{0, 4, 8, 10, 12},
			[]int{1, 3, 6, 9, 13},
			[]int{1, 5, 8, 11, 13},
			[]int{0, 3, 7, 11, 14},
			[]int{2, 5, 7, 9, 12},
		},
		Reels: [][]int{
			[]int{1, 7, 7, 9, 1, 8, 5, 0, 6, 0, 4, 1, 2, 0, 0, 0, 1, 5, 1, 1, 1, 6, 8, 8, 7, 7, 7, 8, 0, 5, 1, 8, 7, 1, 9, 6, 7, 1, 5, 6, 3, 0, 4, 3, 3, 7, 8, 4, 9, 3},
			[]int{7, 1, 9, 9, 0, 5, 8, 8, 3, 5, 1, 0, 5, 6, 6, 8, 1, 0, 0, 0, 3, 6, 2, 1, 1, 1, 7, 3, 6, 0, 3, 3, 7, 3, 1, 9, 3, 3, 1, 2, 8, 6, 6, 7, 0, 3, 2, 3, 5, 8},
			[]int{5, 1, 5, 7, 7, 0, 0, 5, 0, 0, 0, 0, 1, 2, 4, 7, 7, 7, 1, 4, 6, 2, 1, 9, 6, 0, 1, 1, 1, 1, 2, 5, 5, 0, 7, 9, 8, 8, 4, 3, 4, 7, 2, 2, 6, 9, 0, 6, 1, 6},
			[]int{0, 1, 4, 4, 5, 3, 1, 0, 2, 1, 7, 1, 8, 1, 1, 3, 2, 0, 0, 0, 6, 0, 5, 0, 4, 2, 7, 5, 1, 9, 0, 3, 0, 9, 0, 3, 0, 4, 7, 0, 5, 2, 9, 0, 8, 6, 5, 6, 0, 6},
			[]int{9, 2, 1, 7, 6, 0, 9, 2, 7, 2, 8, 1, 3, 0, 0, 6, 0, 7, 7, 2, 5, 1, 3, 5, 1, 1, 0, 1, 6, 2, 0, 0, 2, 9, 1, 0, 0, 9, 4, 4, 0, 1, 1, 7, 6, 0, 9, 7, 3, 1},
		},
	}
}

type Slot struct {
	Odds     [][]int
	PayLines [][]int
	Reels    [][]int
}

//取得盤面
func (t *Slot) GetTable() [][]int32 {

	a := [][]int32{}

	for i := 0; i < len(t.Reels); i++ {
		a = append(a, []int32{})
		r := rand.Intn(len(t.Reels[i]))
		for j := 0; j < 3; j++ {
			index := (r + j) % len(t.Reels[i])
			a[i] = append(a[i], int32(t.Reels[i][index]))
		}
	}

	return a
}

//取得支付線
func (t *Slot) Result(table [][]int32) *pb.Result {

	//盤面 []int
	//支付線號[]int
	//幾連線[]int
	//賠率 int

	//Table
	resulttable := []int32{}
	for i := 0; i < len(table); i++ {
		resulttable = append(resulttable, table[i]...)
	}

	odds := 0
	paylinesnum := []int32{}
	paylinescount := []int32{}

	for i := 0; i < len(t.PayLines); i++ {
		ii, jj := t.GetPayLinePos(t.PayLines[i][0])
		firstsymbol := table[ii][jj]
		line := 0
		for j := 1; j < len(table); j++ {
			iii, jjj := t.GetPayLinePos(t.PayLines[i][j])
			if firstsymbol == table[iii][jjj] {
				line++
			} else {
				break
			}
		}
		//有賠率
		if line >= 2 {
			odds += t.GetOdds(int(firstsymbol), line)
			paylinesnum = append(paylinesnum, int32(i))
			paylinescount = append(paylinescount, int32(line))
		}

	}

	return &pb.Result{
		Table:         resulttable,
		Paylinesnum:   paylinesnum,
		Paylinescount: paylinescount,
		Odds:          int32(odds),
	}
}

//計算線的位置
func (t *Slot) GetPayLinePos(index int) (int, int) {
	return (index / 3), (index % 3)
}

//取得賠率
func (t *Slot) GetOdds(index int, line int) int {
	return t.Odds[index][line]
}
