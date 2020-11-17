package cloud

type AddNodeResult struct {
	Result     string `json:"result"`
	InstanceID string `json:"instanceid"`
	ErrMessage string `json:"errMessage"`
}

type AddNodeResults []AddNodeResult
