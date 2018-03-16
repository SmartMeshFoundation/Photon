package gopjnath

/*
#include <pjnath.h>
#include <pjlib-util.h>
#include <pjlib.h>
extern int Testmain();
extern pj_bool_t testmain();
*/
import "C"

func Testmain() {
	C.testmain()
}

