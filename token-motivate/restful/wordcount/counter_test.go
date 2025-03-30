package wordcount

import "testing"

func TestWordCounter_Stat(t *testing.T) {
	tests := []struct {
		language string
		input    string
		total    int
		words    int
	}{
		{"English1", "hello,playground2022", 3, 2},
		{"English2", "hello, playground", 3, 2},
		{"English3", "👍What makes you special what is around you is enjoy your life! 👍👍", 16, 15},
		{"Chinese1", "大明帝国", 4, 4},
		{"EnglishChinese1", "Hello大明帝国", 5, 5},
		{"EnglishChinese2", "Hello 大明帝国", 5, 5},
		{"EnglishChinese3", "Hello，大明帝国。", 7, 5},
		{"link1", "Hello，大明帝国大明帝国。https://picd.zhimg.com/v2-0a95a82ea9f5957f97a2b662d7164b58_r.jpg?source=1940ef5c picture", 12, 10},
		{"Korean1", "hello,새로운승리를향하여！", 12, 10},
		{"Korean2", "hello,당신은 사랑 받기 위해 태어난 사람 ^-^당신의 삶속에서 그 사랑 받고 있지요", 34, 32},
		{"Japanese1", "ハチと先生は、一緒にうちへ帰ります", 17, 16},
		{"Japanese2", "忠犬八公 ハチと先生は、一緒にうちへ帰ります。\nその日も、ハチは、朝、先生と一緒に渋谷駅へ行きました。", 49, 43},
		{"German", "Ein viertel Jahr ist bereits vergangen.", 7, 6},
		{"Spanish", "Primera luz del día antes de salir el sol.", 12, 11},
		{"French", "Le franais est la plus belle langue du monde.", 10, 9},
	}
	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			counter := &WordCounter{}
			if counter.Stat(tt.input); counter.Total != tt.total || counter.Words != tt.words {
				t.Errorf("Total = %v, want %v; Words=%v, want %v",
					counter.Total, tt.total, counter.Words, tt.words)
			}
		})
	}
}
