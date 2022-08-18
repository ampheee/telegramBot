package telegram

type Update struct {
	UpdateId int    `json:"update_id"`
	Message  string `json:"update_message"`
}

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

func twoSum(nums []int, target int) []int {
	var temp []int
	for i := 0; i < len(nums); i++ {
		for a := len(nums) - 1; a > i; a-- {
			if nums[i]+nums[a] == target && a != i {
				temp = []int{i, a}
			}
		}
	}
	return temp
}
