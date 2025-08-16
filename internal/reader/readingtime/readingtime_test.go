// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package readingtime

import (
	"testing"
	"time"
)

var samples = map[string]string{
	"shortenglish": `This is a short paragraph in english, less than 250 chars.`,
	"shortchinese": ` 労問委格名町違載式新青脂通由。割止書円画民京般著治登門画拡下。有国同観教田美森素説砂者徴多。上治速相支存色分繰年活元事集遣逆山`,
	"english": `
	In turpis lacus, sollicitudin non accumsan sed, suscipit eget magna. Morbi id
	neque enim. Aenean ac lacus consectetur, accumsan elit ac, suscipit dui. Donec
	congue mi et nisl bibendum, venenatis fringilla orci tristique. Nullam ullamcorper
	cursus justo, ac iaculis ante euismod a. Fusce dapibus lacus arcu, consectetur
	porttitor odio finibus ac. Integer dictum faucibus egestas. Etiam magna diam, placerat
	sed velit vitae, lobortis accumsan nisi. Sed viverra dui in odio commodo dapibus.
	Sed pulvinar metus finibus, hendrerit diam eu, faucibus lectus. Mauris est tellus,
	convallis et velit sit amet, convallis sagittis nunc. Quisque at ex leo. Donec eget leo
	vel nibh porta molestie. Aenean pellentesque purus non laoreet aliquam.

	In feugiat eget arcu nec sodales. Nunc rutrum felis in tellus venenatis, sit
	amet tincidunt augue varius. Nunc nec dignissim quam. In euismod gravida rhoncus.
	Vivamus eget nibh sed diam malesuada facilisis. Donec ac convallis elit. Fusce
	fermentum tincidunt est. Nunc viverra, eros in gravida convallis, ex augue vehicula
	magna, sed tincidunt metus sem et mauris. In pretium purus odio, a auctor tellus
	ornare vel. Donec ac dolor pulvinar, placerat elit eget, ultrices nisi. Donec
	tincidunt magna eget pretium sodales. In urna lorem, consectetur in fringilla eget,
	rutrum et erat. Proin fringilla, lectus eget commodo consequat, est massa lacinia
	lorem, ut ultricies nunc erat id sapien.

	Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce fermentum id
	sem sed commodo. Ut eget mauris eu lectus mollis aliquam. Fusce convallis, quam
	vel volutpat aliquet, nunc sem rhoncus magna, a iaculis enim ex nec neque.
	Suspendisse vel imperdiet leo. Quisque ultrices semper commodo. Pellentesque nec libero et
	mauris gravida porta vitae id nunc. Fusce sed sem sed augue gravida ultricies at nec
	turpis. Sed semper eu urna sit amet malesuada. Suspendisse blandit condimentum elit,
	in scelerisque tellus convallis eu. Nunc eleifend sem et mauris vestibulum
	mattis. Praesent ultricies pellentesque eros non posuere.
	`,
	"chinese": `
労問委格名町違載式新青脂通由。割止書円画民京般著治登門画拡下。有国同観教田美森素説砂者徴多。上治速相支存色分繰年活元事集遣逆山。身消年森発世財間世変悲原記潟旅好手真今。現通浪口特愛始信川節身方一表著購。郁不使権草定内防並要更一条露加。載交源図訴際属年券重供健三洗。事北残却女鮎朝分要廷込宣政愛無投事。

問警技亮参沼洗請米物模人。誰探重午局新戦報投性病庭。典向載問千著書故表視新権最石車音端乏大。白僚三掲局係仕表広無旧見要最裁。額寄済生年余講前本次載隊劇。権成観始応泉早高拓了経地本稼室目犯井出。暮載必広傷内校岡公南散広転行別釈。康運行関本掲隠泉傷退報告。独変年換差取予口男旅挑講禁姿。出芳工類胸管払時済潟髪内豊。

康浴部問玲玉追球化就店岡問画路投。施先太業阪能敏所陸不供探掲方用。手右演社援発示竹育対橋除際愛功旬転好使公。利時改本項輸属嘆員複携者地剤。天政朝戸祝言月接住世黙極者議編連。囲淑覧重弾必治物健賄開頂外称豊開名銀戸院。政稿調励廃演手生告題営味董演何南峰貨。学横公得行提大品回猿齢利込家前役把煎。天代者内身慢作業署間地日。

中個興本広坂態掲神中能等無滞長対。号処月画界意気様党目購栃欠歌暮。一耳供意盛四俊健必財下画例本判著堺要北王。宮大攻人水一備治首闘振円分建前趣校。目少供午見掲岡安画入情薦続土世始。診読格七久改急目斉実配正。性止月模多様更社発掲雪奇芸量全兵経負。予転済反問止下生買再無旅的。模治明以共会必華浅知館版領送。
	`,
	"korean": `
	세계 인권 선언(世界人權宣言, 영어: Universal Declaration of Human Rights, UDHR)은 1948년 12월 10일 파리에서 열린 제3회 유엔 총회에서 채택된 인권에 관한 세계 선언문이다.[1] 2차 세계대전 전후로 전 세계에 만연하였던 인권침해 사태에 대한 인류의 반성을 촉구하고, 모든 인간의 기본적 권리를 존중해야 한다는 유엔 헌장의 취지를 구체화 하였다.[2] 시민적, 정치적 권리가 중심이지만 노동자의 단결권, 교육에 관한 권리, 예술을 향유할 권리 등 경제적, 사회적, 문화적 권리에 대하여서도 규정하고 있다.[1]
	초안은 1946년 존 험프리가 작성하였다.[3] 인권선언문은 전문과 본문의 30개 조에 개인의 기본적인 자유와 함께 노동권적 권리, 생존권적 권리를 오늘날의 진보적인 국가의 헌법에서 규정하는 인권보장과 같이 자세히 규정하고 있다.[4] 프랑스 파리 샤요 궁(Palais de Chaillot)에서 열린 3번째 회의에서 당시 국제연합 가입국 58개 국가 중 48개 국가가 찬성하여 유엔 총회 결의 217 A (III)로 승인되었다.
	초안은 1946년 존 험프리가 작성하였다.[3] 인권선언문은 전문과 본문의 30개 조에 개인의 기본적인 자유와 함께 노동권적 권리, 생존권적 권리를 오늘날의 진보적인 국가의 헌법에서 규정하는 인권보장과 같이 자세히 규정하고 있다.[4] 프랑스 파리 샤요 궁(Palais de Chaillot)에서 열린 3번째 회의에서 당시 국제연합 가입국 58개 국가 중 48개 국가가 찬성하여 유엔 총회 결의 217 A (III)로 승인되었다.
	초안은 1946년 존 험프리가 작성하였다.[3] 인권선언문은 전문과 본문의 30개 조에 개인의 기본적인 자유와 함께 노동권적 권리, 생존권적 권리를 오늘날의 진보적인 국가의 헌법에서 규정하는 인권보장과 같이 자세히 규정하고 있다.[4] 프랑스 파리 샤요 궁(Palais de Chaillot)에서 열린 3번째 회의에서 당시 국제연합 가입국 58개 국가 중 48개 국가가 찬성하여 유엔 총회 결의 217 A (III)로 승인되었다.
	초안은 1946년 존 험프리가 작성하였다.[3] 인권선언문은 전문과 본문의 30개 조에 개인의 기본적인 자유와 함께 노동권적 권리, 생존권적 권리를 오늘날의 진보적인 국가의 헌법에서 규정하는 인권보장과 같이 자세히 규정하고 있다.[4] 프랑스 파리 샤요 궁(Palais de Chaillot)에서 열린 3번째 회의에서 당시 국제연합 가입국 58개 국가 중 48개 국가가 찬성하여 유엔 총회 결의 217 A (III)로 승인되었다.
	초안은 1946년 존 험프리가 작성하였다.[3] 인권선언문은 전문과 본문의 30개 조에 개인의 기본적인 자유와 함께 노동권적 권리, 생존권적 권리를 오늘날의 진보적인 국가의 헌법에서 규정하는 인권보장과 같이 자세히 규정하고 있다.[4] 프랑스 파리 샤요 궁(Palais de Chaillot)에서 열린 3번째 회의에서 당시 국제연합 가입국 58개 국가 중 48개 국가가 찬성하여 유엔 총회 결의 217 A (III)로 승인되었다.
	초안은 1946년 존 험프리가 작성하였다.[3] 인권선언문은 전문과 본문의 30개 조에 개인의 기본적인 자유와 함께 노동권적 권리, 생존권적 권리를 오늘날의 진보적인 국가의 헌법에서 규정하는 인권보장과 같이 자세히 규정하고 있다.[4] 프랑스 파리 샤요 궁(Palais de Chaillot)에서 열린 3번째 회의에서 당시 국제연합 가입국 58개 국가 중 48개 국가가 찬성하여 유엔 총회 결의 217 A (III)로 승인되었다.
	초안은 1946년 존 험프리가 작성하였다.[3] 인권선언문은 전문과 본문의 30개 조에 개인의 기본적인 자유와 함께 노동권적 권리, 생존권적 권리를 오늘날의 진보적인 국가의 헌법에서 규정하는 인권보장과 같이 자세히 규정하고 있다.[4] 프랑스 파리 샤요 궁(Palais de Chaillot)에서 열린 3번째 회의에서 당시 국제연합 가입국 58개 국가 중 48개 국가가 찬성하여 유엔 총회 결의 217 A (III)로 승인되었다.
	초안은 1946년 존 험프리가 작성하였다.[3] 인권선언문은 전문과 본문의 30개 조에 개인의 기본적인 자유와 함께 노동권적 권리, 생존권적 권리를 오늘날의 진보적인 국가의 헌법에서 규정하는 인권보장과 같이 자세히 규정하고 있다.[4] 프랑스 파리 샤요 궁(Palais de Chaillot)에서 열린 3번째 회의에서 당시 국제연합 가입국 58개 국가 중 48개 국가가 찬성하여 유엔 총회 결의 217 A (III)로 승인되었다.
	초안은 1946년 존 험프리가 작성하였다.[3] 인권선언문은 전문과 본문의 30개 조에 개인의 기본적인 자유와 함께 노동권적 권리, 생존권적 권리를 오늘날의 진보적인 국가의 헌법에서 규정하는 인권보장과 같이 자세히 규정하고 있다.[4] 프랑스 파리 샤요 궁(Palais de Chaillot)에서 열린 3번째 회의에서 당시 국제연합 가입국 58개 국가 중 48개 국가가 찬성하여 유엔 총회 결의 217 A (III)로 승인되었다.
	`,
}

func TestEstimateReadingTime(t *testing.T) {
	expected := map[string]int{
		"shortenglish": 1,
		"shortchinese": 1,
		"english":      2,
		"chinese":      2,
		"korean":       5,
	}

	for language, sample := range samples {
		got := textReadingTime(sample, 200, 500)
		want := expected[language]
		if got != want {
			t.Errorf(`Wrong reading time, got %d instead of %d for %s`, got, want, language)
		}
	}
}

func BenchmarkEstimateReadingTime(b *testing.B) {
	for b.Loop() {
		for _, sample := range samples {
			textReadingTime(sample, 200, 500)
		}
	}
}

func TestISO8601DurationParsing(t *testing.T) {
	var scenarios = []struct {
		duration string
		expected time.Duration
	}{
		// Live streams and radio.
		{"PT0M0S", 0},
		// https://www.youtube.com/watch?v=HLrqNhgdiC0
		{"PT6M20S", (6 * time.Minute) + (20 * time.Second)},
		// https://www.youtube.com/watch?v=LZa5KKfqHtA
		{"PT5M41S", (5 * time.Minute) + (41 * time.Second)},
		// https://www.youtube.com/watch?v=yIxEEgEuhT4
		{"PT51M52S", (51 * time.Minute) + (52 * time.Second)},
		// https://www.youtube.com/watch?v=bpHf1XcoiFs
		{"PT80M42S", (1 * time.Hour) + (20 * time.Minute) + (42 * time.Second)},
		// Hours only
		{"PT2H", 2 * time.Hour},
		// Seconds only
		{"PT30S", 30 * time.Second},
		// Hours and minutes
		{"PT1H30M", (1 * time.Hour) + (30 * time.Minute)},
		// Hours and seconds
		{"PT2H45S", (2 * time.Hour) + (45 * time.Second)},
		// Empty duration
		{"PT", 0},
	}

	for _, tc := range scenarios {
		result, err := parseISO8601Duration(tc.duration)
		if err != nil {
			t.Errorf("Got an error when parsing %q: %v", tc.duration, err)
		}

		if tc.expected != result {
			t.Errorf(`Unexpected result, got %v for duration %q`, result, tc.duration)
		}
	}
}

func TestISO8601DurationParsingErrors(t *testing.T) {
	var errorScenarios = []struct {
		duration    string
		expectedErr string
	}{
		// Missing PT prefix
		{"6M20S", "the period doesn't start with PT"},
		// Unsupported Year specifier
		{"PT1Y", "the 'Y' specifier isn't supported"},
		// Unsupported Week specifier
		{"PT2W", "the 'W' specifier isn't supported"},
		// Unsupported Day specifier
		{"PT3D", "the 'D' specifier isn't supported"},
		// Invalid number for hours (letter at start of number)
		{"PTaH", "invalid character in the period"},
		// Invalid number for minutes (letter at start of number)
		{"PTbM", "invalid character in the period"},
		// Invalid number for seconds (letter at start of number)
		{"PTcS", "invalid character in the period"},
		// Invalid character in the middle of a number
		{"PT1a2H", "invalid character in the period"},
		{"PT3b4M", "invalid character in the period"},
		{"PT5c6S", "invalid character in the period"},
		// Test cases for actual ParseFloat errors (empty number before specifier)
		{"PTH", "strconv.ParseFloat: parsing \"\": invalid syntax"},
		{"PTM", "strconv.ParseFloat: parsing \"\": invalid syntax"},
		{"PTS", "strconv.ParseFloat: parsing \"\": invalid syntax"},
		// Invalid character
		{"PT1X", "invalid character in the period"},
		// Invalid character mixed
		{"PT1H@M", "invalid character in the period"},
	}

	for _, tc := range errorScenarios {
		_, err := parseISO8601Duration(tc.duration)
		if err == nil {
			t.Errorf("Expected an error when parsing %q, but got none", tc.duration)
		} else if err.Error() != tc.expectedErr {
			t.Errorf("Expected error %q when parsing %q, but got %q", tc.expectedErr, tc.duration, err.Error())
		}
	}
}
