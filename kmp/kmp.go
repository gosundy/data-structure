package kmp

type Kmp struct {
	model []byte
	next  []int
}

func New(model string) *Kmp {
	kmp := &Kmp{model: []byte(model), next: make([]int, len(model))}
	kmp.generateNext()
	kmp.next[0] = -1
	return kmp
}

func (kmp *Kmp) generateNext() {
	k := -1 //left index
	j := 1  //right index
	kmp.next[0] = -1
	for j < len(kmp.model)-1 {
		if k == -1 || kmp.model[k] == kmp.model[j] {
			k++
			j++
			kmp.next[j] = k
		} else {
			k = kmp.next[k]
		}
	}
}

func (kmp *Kmp) Compare(compared string) int {
	k := 0
	i := 0
	for i < len(compared) && k < len(kmp.model) {
		if kmp.model[k] == compared[i] {
			k++
			i++
			continue
		}
		k = kmp.next[k]
		if k == -1 {
			k = 0
			i++
		}
	}
	if k == len(kmp.model) {
		return i - k
	}
	return -1
}
