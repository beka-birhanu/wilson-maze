package main

func main() {
	maz, _ := New(10, 10)

	if maz.PopulateReward(struct {
		RewardOne      int32
		RewardTwo      int32
		RewardTypeProb float32
	}{RewardOne: 1, RewardTwo: 5, RewardTypeProb: 0.9}) != nil {
		return
	}
}
