#!/bin/sh
go tool pprof mem-0.pprof http://127.0.0.1:3000/debug/pprof/heap
go tool pprof mem-1.pprof http://127.0.0.1:3001/debug/pprof/heap
go tool pprof mem-2.pprof http://127.0.0.1:3002/debug/pprof/heap
go tool pprof mem-3.pprof http://127.0.0.1:3003/debug/pprof/heap
go tool pprof mem-4.pprof http://127.0.0.1:3004/debug/pprof/heap
go tool pprof mem-5.pprof http://127.0.0.1:3005/debug/pprof/heap
go tool pprof mem-6.pprof http://127.0.0.1:3006/debug/pprof/heap
go tool pprof mem-7.pprof http://127.0.0.1:3007/debug/pprof/heap
go tool pprof mem-8.pprof http://127.0.0.1:3008/debug/pprof/heap
go tool pprof mem-9.pprof http://127.0.0.1:3009/debug/pprof/heap

