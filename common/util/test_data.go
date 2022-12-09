package util

type intJudge int

func (s intJudge) IsTrue() bool {
	return s%2 == 0
}

func (s intJudge) Value() int {
	return int(s)
}

var testJudge = func() (re []intJudge) {
	for i := 0; i < 100000; i++ {
		re = append(re, intJudge(i))
	}
	return
}()

var testData = func() (re []int) {
	for i := 0; i < 100000; i++ {
		re = append(re, i)
	}
	return
}()
