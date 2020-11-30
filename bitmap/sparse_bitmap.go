package bitmap
type SparseBitmap struct {

}
type TierTree struct {
	savedType int
	count int
	savedMap map[int16]struct{}
	savedBit
}