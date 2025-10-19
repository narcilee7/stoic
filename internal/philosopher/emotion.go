package philosopher

import "strings"

type EmotionAnalyzer struct {
	patterns map[string][]string
}

func NewEmotionAnalyzer() *EmotionAnalyzer {
	return &EmotionAnalyzer{
		patterns: map[string][]string{
			"happy":    {"happy", "joy", "excited", "glad", "cheerful", "😊", "😄", "🎉"},
			"sad":      {"sad", "depressed", "down", "unhappy", "miserable", "😢", "😭", "💔"},
			"angry":    {"angry", "mad", "furious", "irritated", "annoyed", "😠", "😡", "🤬"},
			"anxious":  {"anxious", "worried", "nervous", "stressed", "concerned", "😰", "😟", "😨"},
			"calm":     {"calm", "peaceful", "relaxed", "serene", "tranquil", "😌", "🧘", "☮️"},
			"confused": {"confused", "lost", "uncertain", "puzzled", "don't understand", "🤔", "😕", "😵"},
		},
	}
}

func (ea *EmotionAnalyzer) AnalyzeEmotion(text string) string {
	text = strings.ToLower(text)
	scores := make(map[string]int)

	for emotion, patterns := range ea.patterns {
		for _, pattern := range patterns {
			if strings.Contains(text, pattern) {
				scores[emotion]++
			}
		}
	}

	// 找到得分最高的emotion
	maxScore := 0
	dominantEmotion := "neutral"
	for emotion, score := range scores {
		if score > maxScore {
			maxScore = score
			dominantEmotion = emotion
		}
	}
	return dominantEmotion
}
