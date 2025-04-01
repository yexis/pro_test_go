package slice

import (
	"fmt"
	"strconv"
	"testing"
)

// easyjson:json
type ReadData struct {
	CurrentDevice    bool    `json:"-" redis:"-"`
	LogID            string  `json:"-" redis:"-"`
	DumiID           string  `json:"uid" redis:"-"`
	ClientID         string  `json:"-" redis:"-"`
	DeviceID         string  `json:"cuid" redis:"deviceid"`
	SN               string  `json:"sn" redis:"sn"`
	WakeupSN         string  `json:"wakeup_sn" redis:"wakeup_sn"`
	Date             int64   `json:"date" redis:"date"`
	End2End          int     `json:"e2e" redis:"e2e_dci"`
	Weight           string  `json:"-" redis:"dci_weight"`
	AudioPlaying     string  `json:"play,omitempty" redis:"audio_player,omitempty"`
	Snr              string  `json:"snr,omitempty" redis:"snr,omitempty"`
	AudioScore       string  `json:"score,omitempty" redis:"audio_score,omitempty"`
	snr              float32 `json:"-" redis:"-"`
	audioScore       int     `json:"-" redis:"-"`
	WakeUpWordLength int     `json:"word_len,omitempty" redis:"wakeup_word_length,omitempty"`
	RegistTime       int64   `json:"regist_time,omitempty" redis:"regist_time,omitempty"`
	UpdateTime       int64   `json:"update_time,omitempty" redis:"update_time,omitempty"`
	ErrNo            string  `json:"err_no,omitempty" redis:"err_no,omitempty"`
}

func TestSlice(t *testing.T) {
	rd := []ReadData{{LogID: "1"}}
	for i := 0; i < 10000000; i++ {
		go func() {
			rd1 := []ReadData{{LogID: strconv.Itoa(i)}}
			rd = rd1
		}()

		go func() {
			if rd == nil {
				fmt.Println("empty")
			} else {
				//fmt.Println(rd)
			}
		}()
	}

}
