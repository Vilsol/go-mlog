package runtime

import (
	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/go-mlog/runtime"
	"testing"
)

func TestOperationDraw(t *testing.T) {
	_, err := runtime.Parse(`draw`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'draw': expecting at least 2 arguments", err.Error())

	_, err = runtime.Parse(`draw clear 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'draw clear 1': expecting at least 4 arguments", err.Error())

	_, err = runtime.Parse(`draw color 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'draw color 1': expecting at least 5 arguments", err.Error())

	_, err = runtime.Parse(`draw line 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'draw line 1': expecting at least 5 arguments", err.Error())

	_, err = runtime.Parse(`draw rect 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'draw rect 1': expecting at least 5 arguments", err.Error())

	_, err = runtime.Parse(`draw lineRect 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'draw lineRect 1': expecting at least 5 arguments", err.Error())

	_, err = runtime.Parse(`draw poly 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'draw poly 1': expecting at least 6 arguments", err.Error())

	_, err = runtime.Parse(`draw linePoly 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'draw linePoly 1': expecting at least 6 arguments", err.Error())

	_, err = runtime.Parse(`draw triangle 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'draw triangle 1': expecting at least 7 arguments", err.Error())

	_, err = runtime.Parse(`draw image 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'draw image 1': expecting at least 6 arguments", err.Error())
}

func TestOperationDrawFlush(t *testing.T) {
	_, err := runtime.Parse(`drawflush`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'drawflush': expecting exactly 1 argument", err.Error())
}

func TestOperationWait(t *testing.T) {
	operations, err := runtime.Parse(`wait 0.1`)
	testza.AssertNoError(t, err)
	context, counter := runtime.ConstructContext(map[string]interface{}{})
	err = runtime.ExecuteContext(operations, context, counter)
	testza.AssertNoError(t, err)

	_, err = runtime.Parse(`wait`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'wait': expecting exactly 1 argument", err.Error())
}

func TestOperationRead(t *testing.T) {
	_, err := runtime.Parse(`read`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'read': expecting exactly 3 arguments", err.Error())
}

func TestOperationWrite(t *testing.T) {
	_, err := runtime.Parse(`write`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'write': expecting exactly 3 arguments", err.Error())
}

func TestOperationEnd(t *testing.T) {
	operations, err := runtime.Parse(`end`)
	testza.AssertNoError(t, err)
	context, counter := runtime.ConstructContext(map[string]interface{}{})
	err = runtime.ExecuteContext(operations, context, counter)
	testza.AssertNoError(t, err)

	_, err = runtime.Parse(`end 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'end 1': requires no arguments", err.Error())
}

func TestOperationPrint(t *testing.T) {
	_, err := runtime.Parse(`print`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'print': expecting exactly 1 argument", err.Error())
}

func TestOperationPrintFlush(t *testing.T) {
	_, err := runtime.Parse(`printflush`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'printflush': expecting exactly 1 argument", err.Error())
}

func TestOperationOp(t *testing.T) {
	operations, err := runtime.Parse(`op add r0 25 12
op sub r1 25 12
op mul r2 25 12
op div r3 25 12
op idiv r4 25 12
op mod r5 25 12
op pow r6 25 12
op equal r7 25 12
op notEqual r8 25 12
op land r9 25 12
op lessThan r10 25 12
op lessThanEq r11 25 12
op greaterThan r12 25 12
op greaterThanEq r13 25 12
op strictEqual r14 25 12
op shl r15 25 12
op shr r16 25 12
op or r17 25 12
op and r18 25 12
op xor r19 25 12
op not r20 25 12
op max r21 25 12
op min r22 25 12
op angle r23 25 12
op len r24 25 12
op noise r25 25 12
op abs r26 25 12
op log r27 25 12
op log10 r28 25 12
op sin r29 25 12
op cos r30 25 12
op tan r31 25 12
op floor r32 25 12
op ceil r33 25 12
op sqrt r34 25 12
op rand r35 25 12`)
	testza.AssertNoError(t, err)
	context, counter := runtime.ConstructContext(map[string]interface{}{})
	err = runtime.ExecuteContext(operations, context, counter)
	testza.AssertNoError(t, err)

	testza.AssertEqual(t, float64(37), context.Variables["r0"].Value)
	testza.AssertEqual(t, float64(13), context.Variables["r1"].Value)
	testza.AssertEqual(t, float64(300), context.Variables["r2"].Value)
	testza.AssertEqual(t, float64(2.0833333333333335), context.Variables["r3"].Value)
	testza.AssertEqual(t, float64(2), context.Variables["r4"].Value)
	testza.AssertEqual(t, float64(1), context.Variables["r5"].Value)
	testza.AssertEqual(t, float64(5.960464477539062e+16), context.Variables["r6"].Value)
	testza.AssertEqual(t, false, context.Variables["r7"].Value)
	testza.AssertEqual(t, true, context.Variables["r8"].Value)
	testza.AssertEqual(t, true, context.Variables["r9"].Value)
	testza.AssertEqual(t, false, context.Variables["r10"].Value)
	testza.AssertEqual(t, false, context.Variables["r11"].Value)
	testza.AssertEqual(t, true, context.Variables["r12"].Value)
	testza.AssertEqual(t, true, context.Variables["r13"].Value)
	testza.AssertEqual(t, false, context.Variables["r14"].Value)
	testza.AssertEqual(t, int64(102400), context.Variables["r15"].Value)
	testza.AssertEqual(t, int64(0), context.Variables["r16"].Value)
	testza.AssertEqual(t, int64(29), context.Variables["r17"].Value)
	testza.AssertEqual(t, int64(8), context.Variables["r18"].Value)
	testza.AssertEqual(t, int64(21), context.Variables["r19"].Value)
	testza.AssertEqual(t, int64(-26), context.Variables["r20"].Value)
	testza.AssertEqual(t, float64(25), context.Variables["r21"].Value)
	testza.AssertEqual(t, float64(12), context.Variables["r22"].Value)
	testza.AssertEqual(t, float64(25.835288062773845), context.Variables["r23"].Value)
	testza.AssertEqual(t, float64(27.730849247724095), context.Variables["r24"].Value)
	testza.AssertEqual(t, float64(0.49890396713293766), context.Variables["r25"].Value)
	testza.AssertEqual(t, float64(25), context.Variables["r26"].Value)
	testza.AssertEqual(t, float64(3.2188758248682006), context.Variables["r27"].Value)
	testza.AssertEqual(t, float64(1.3979400086720375), context.Variables["r28"].Value)
	testza.AssertEqual(t, float64(0.42261826174069944), context.Variables["r29"].Value)
	testza.AssertEqual(t, float64(0.9063077870366499), context.Variables["r30"].Value)
	testza.AssertEqual(t, float64(0.4663076581549986), context.Variables["r31"].Value)
	testza.AssertEqual(t, float64(25), context.Variables["r32"].Value)
	testza.AssertEqual(t, float64(25), context.Variables["r33"].Value)
	testza.AssertEqual(t, float64(5), context.Variables["r34"].Value)

	rand := context.Variables["r35"].Value.(float64)
	testza.AssertTrue(t, rand > float64(0) && rand < float64(25))

	_, err = runtime.Parse(`op`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'op': expecting at least 2 arguments", err.Error())

	_, err = runtime.Parse(`op "hello" 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'op \"hello\" 1': target variable must not be a string", err.Error())
}

func TestOperationSet(t *testing.T) {
	operations, err := runtime.Parse(`set str "hello"
set int 1234
set float 1.234`)

	if err != nil {
		panic(err)
	}

	context, counter := runtime.ConstructContext(map[string]interface{}{})
	err = runtime.ExecuteContext(operations, context, counter)
	testza.AssertNil(t, err)

	testza.AssertEqual(t, "hello", context.Variables["str"].Value)
	testza.AssertEqual(t, int64(1234), context.Variables["int"].Value)
	testza.AssertEqual(t, float64(1.234), context.Variables["float"].Value)

	_, err = runtime.Parse(`set`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'set': expecting exactly 2 arguments", err.Error())

	_, err = runtime.Parse(`set "hello" 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'set \"hello\" 1': target variable must not be a string", err.Error())
}

func TestOperationJump(t *testing.T) {
	operations, err := runtime.Parse(`set a 0
set b 0
jump 12 always
set a 1
set b 1
end
jump 4 equal 1 1
jump 6 notEqual 1 2
jump 7 lessThan 1 2
jump 8 lessThanEq 1 1
jump 9 greaterThan 2 1
jump 10 greaterThanEq 1 1
jump 11 strictEqual 1 1`)

	if err != nil {
		panic(err)
	}

	context, counter := runtime.ConstructContext(map[string]interface{}{})
	err = runtime.ExecuteContext(operations, context, counter)
	testza.AssertNil(t, err)

	testza.AssertEqual(t, int64(0), context.Variables["a"].Value)
	testza.AssertEqual(t, int64(1), context.Variables["b"].Value)

	_, err = runtime.Parse(`jump`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'jump': expecting at least 2 arguments", err.Error())

	_, err = runtime.Parse(`jump 0 equal`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'jump 0 equal': expecting exactly 4 arguments", err.Error())

	_, err = runtime.Parse(`jump "hello" equal 1 1`)
	testza.AssertNotNil(t, err)
	testza.AssertEqual(t, "error on line 0: 'jump \"hello\" equal 1 1': jump target must not be a string", err.Error())
}
