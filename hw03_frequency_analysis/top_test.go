package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
// var taskWithAsteriskIsCompleted = false

var (
	textOne = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`
	expectedSliceOne = []string{
		"он",        // 8
		"а",         // 6
		"и",         // 6
		"ты",        // 5
		"что",       // 5
		"-",         // 4
		"Кристофер", // 4
		"если",      // 4
		"не",        // 4
		"то",        // 4
	}
	// expectedSliceOneAdditional = []string{"а", "он", "и", "ты", "что", "в", "его", "если", "кристофер", "не"}.

	textTwo          = ""
	expectedSliceTwo []string

	textThree          = "aaa aaa aaa    aaa"
	expectedSliceThree = []string{"aaa"}

	textFour          = "cat and dog, one dog,two cats and one man"
	expectedSliceFour = []string{"and", "one", "cat", "cats", "dog,", "dog,two", "man"}

	textFive          = "ddd ddd    ddd ccc  ccc bbb bbb aaa aaa aaa"
	expectedSliceFive = []string{"aaa", "ddd", "bbb", "ccc"}

	// textSix          = "какой-то человек! с        кем-то когда-то 45 !№%какой-то???? почему-то ПОЧему-То!!!"
	// expectedSliceSixAdditional = []string{"какой-то", "почему-то", "45", "кем-то", "когда-то", "с", "человек"}.
)

func TestTop10(t *testing.T) {
	tests := []struct {
		name          string
		inputText     string
		expectedSlice []string
	}{
		{name: "test1", inputText: textOne, expectedSlice: expectedSliceOne},
		{name: "test2: empty text", inputText: textTwo, expectedSlice: expectedSliceTwo},
		{name: "test3: similar words", inputText: textThree, expectedSlice: expectedSliceThree},
		{name: "test4", inputText: textFour, expectedSlice: expectedSliceFour},
		{name: "test5", inputText: textFive, expectedSlice: expectedSliceFive},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actual := Top10(tc.inputText)
			require.Equal(t, tc.expectedSlice, actual)
		})
	}
}

// func TestTop10Additional(t *testing.T) {
//	tests := []struct {
//		name          string
//		inputText     string
//		expectedSlice []string
//	}{
//		{name: "test1 (additional func)", inputText: textOne, expectedSlice: expectedSliceOneAdditional},
//		{name: "test2 (additional func)", inputText: textSix, expectedSlice: expectedSliceSixAdditional},
//	}
//	for _, tc := range tests {
//		tc := tc
//		t.Run(tc.name, func(t *testing.T) {
//			actual := Top10Additional(tc.inputText)
//			require.Equal(t, tc.expectedSlice, actual)
//		})
//	}
//}
