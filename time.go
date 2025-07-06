package main

import (
	"time"

	"codeberg.org/anaseto/goal"
)

/*
Convert string to date ordinal
2000-01-01 -> 0
*/
func vfTimeStrDate(ctx *goal.Context, args []goal.V) goal.V {
	s, ok := args[0].BV().(goal.S)
	if ok {
		return goal.NewI(strToDtOrd(string(s)))
	}
	as, ok := args[0].BV().(*goal.AS)
	if ok {
		dates := make([]int64, as.Len())
		for i, d := range as.Slice {
			dates[i] = strToDtOrd(d)
		}
		return goal.NewAI(dates)
	}
	return goal.Panicf("Invalid type: %s", args[0].Type())
}

func vfTimeDateStr(ctx *goal.Context, args []goal.V) goal.V {
	if args[0].IsI() {
		return goal.NewS(dtOrdToStr(args[0].I()))
	}
	as, ok := args[0].BV().(*goal.AI)
	if ok {
		dates := make([]string, as.Len())
		for i, d := range as.Slice {
			dates[i] = dtOrdToStr(d)
		}
		return goal.NewAS(dates)
	}
	return goal.Panicf("Invalid type: %s", args[0].Type())
}

func strToDtOrd(s string) int64 {
	t, e := time.Parse(time.DateOnly, s)
	if e != nil {
		return -1
	}
	return (t.Unix() / 86400) - 10957
}

func dtOrdToStr(o int64) string {
	t := time.Unix((o+10957)*86400, 0).UTC()
	return t.Format(time.DateOnly)
}
