package main

import (
	"fmt"
	"io"
	"net/http"

	"codeberg.org/anaseto/goal"
)

func vfHttpGet(ctx *goal.Context, args []goal.V) goal.V {
	url, _ := args[0].BV().(goal.S)
	println(url)
	resp, err := http.Get(string(url))
	if err != nil {
		return goal.Panicf("http.get: bad request %s", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return goal.Panicf("http.get: bad request %s", err)
	}
	return goal.NewS(string(body))
}

func vfHttpRegister(ctx *goal.Context, args []goal.V) goal.V {
	url, _ := args[0].BV().(goal.S)
	if !args[1].IsCallable() {
		return goal.Panicf("http.register: function is not callable")
	}
	fn := args[1]
	http.HandleFunc(string(url), func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		v := fn.Apply(ctx, goal.NewS(string(body)))
		res, ok := v.BV().(goal.S)
		if ok {
			fmt.Fprint(w, string(res))
		} else {
			fmt.Fprint(w, v.String())
		}
	})
	return goal.NewI(1)
}

func vfHttpServe(ctx *goal.Context, args []goal.V) goal.V {
	url, _ := args[0].BV().(goal.S)
	go func() {
		http.ListenAndServe(string(url), nil)
	}()
	return goal.NewI(1)
}
