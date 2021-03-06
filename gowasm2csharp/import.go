// SPDX-License-Identifier: Apache-2.0

package gowasm2csharp

var importFuncBodies = map[string]string{
	// func wasmExit(code int32)
	"runtime.wasmExit": `    var code = go.mem.LoadInt32(local0 + 8);
    go.exited = true;
    go.inst = null;
    go.values = null;
    go.goRefCounts = null;
    go.ids = null;
    go.idPool = null;
    go.Exit(code);`,

	// func wasmWrite(fd uintptr, p unsafe.Pointer, n int32)
	"runtime.wasmWrite": `    var fd = go.mem.LoadInt64(local0 + 8);
    if (fd != 1 && fd != 2)
    {
        throw new NotImplementedException($"fd for runtime.wasmWrite must be 1 or 2 but {fd}");
    }
    var p = go.mem.LoadInt64(local0 + 16);
    var n = go.mem.LoadInt32(local0 + 24);

    // Note that runtime.wasmWrite is used only for print/println so far.
    // Write the buffer to the standard output regardless of fd.
    go.DebugWrite(go.mem.LoadSliceDirectly(p, n));`,

	// func resetMemoryDataView()
	"runtime.resetMemoryDataView": `    // Do nothing.`,

	// func nanotime1() int64
	"runtime.nanotime1": `    go.mem.StoreInt64(local0 + 8, go.PreciseNowInNanoseconds());`,

	// func walltime1() (sec int64, nsec int32)
	"runtime.walltime1": `    var now = go.UnixNowInMilliseconds();
    go.mem.StoreInt64(local0 + 8, (long)(now / 1000));
    go.mem.StoreInt32(local0 + 16, (int)((now % 1000) * 1_000_000));`,

	// func scheduleTimeoutEvent(delay int64) int32
	"runtime.scheduleTimeoutEvent": `    var interval = go.mem.LoadInt64(local0 + 8);
    var id = go.SetTimeout((double)interval);
    go.mem.StoreInt32(local0 + 16, id);`,

	// func clearTimeoutEvent(id int32)
	"runtime.clearTimeoutEvent": `    var id = go.mem.LoadInt32(local0 + 8);
    go.ClearTimeout(id);`,

	// func getRandomData(r []byte)
	"runtime.getRandomData": `    var slice = go.mem.LoadSlice(local0 + 8);
    var bytes = go.GetRandomBytes(slice.Count);
    for (int i = 0; i < slice.Count; i++) {
        slice.Array[slice.Offset + i] = bytes[i];
    }`,

	// func finalizeRef(v ref)
	"syscall/js.finalizeRef": `    int id = (int)go.mem.LoadUint32(local0 + 8);
    go.goRefCounts[id]--;
    if (go.goRefCounts[id] == 0)
    {
        var v = go.values[id];
        go.values[id] = null;
        go.ids.Remove(v);
        go.idPool.Push(id);
    }`,

	// func stringVal(value string) ref
	"syscall/js.stringVal": `    go.StoreValue(local0 + 24, go.mem.LoadString(local0 + 8));`,

	// func valueGet(v ref, p string) ref
	"syscall/js.valueGet": `    var result = JSObject.ReflectGet(go.LoadValue(local0 + 8), go.mem.LoadString(local0 + 16));
    local0 = go.inst.getsp();
    go.StoreValue(local0 + 32, result);`,

	// func valueSet(v ref, p string, x ref)
	"syscall/js.valueSet": `    JSObject.ReflectSet(go.LoadValue(local0 + 8), go.mem.LoadString(local0 + 16), go.LoadValue(local0 + 32));`,

	// func valueDelete(v ref, p string)
	"syscall/js.valueDelete": `    JSObject.ReflectDelete(go.LoadValue(local0 + 8), go.mem.LoadString(local0 + 16));`,

	// func valueIndex(v ref, i int) ref
	"syscall/js.valueIndex": `    go.StoreValue(local0 + 24, JSObject.ReflectGet(go.LoadValue(local0 + 8), go.mem.LoadInt64(local0 + 16).ToString()));`,

	// valueSetIndex(v ref, i int, x ref)
	"syscall/js.valueSetIndex": `    JSObject.ReflectSet(go.LoadValue(local0 + 8), go.mem.LoadInt64(local0 + 16).ToString(), go.LoadValue(local0 + 24));`,

	// func valueCall(v ref, m string, args []ref) (ref, bool)
	"syscall/js.valueCall": `    var v = go.LoadValue(local0 + 8);
    var m = JSObject.ReflectGet(v, go.mem.LoadString(local0 + 16));
    var args = go.LoadSliceOfValues(local0 + 32);
    var result = JSObject.ReflectApply(m, v, args);
    local0 = go.inst.getsp();
    go.StoreValue(local0 + 56, result);
    go.mem.StoreInt8(local0 + 64, 1);`,

	// func valueInvoke(v ref, args []ref) (ref, bool)
	"syscall/js.valueInvoke": `    var v = go.LoadValue(local0 + 8);
    var args = go.LoadSliceOfValues(local0 + 16);
    var result = JSObject.ReflectApply(v, JSObject.Undefined, args);
    local0 = go.inst.getsp();
    go.StoreValue(local0 + 40, result);
    go.mem.StoreInt8(local0 + 48, 1);`,

	// func valueNew(v ref, args []ref) (ref, bool)
	"syscall/js.valueNew": `    var v = go.LoadValue(local0 + 8);
    var args = go.LoadSliceOfValues(local0 + 16);
    var result = JSObject.ReflectConstruct(v, args);
    if (result != null)
    {
        local0 = go.inst.getsp();
        go.StoreValue(local0 + 40, result);
        go.mem.StoreInt8(local0 + 48, 1);
    }
    else
    {
        go.StoreValue(local0 + 40, null);
        go.mem.StoreInt8(local0 + 48, 0);
    }`,

	// func valueLength(v ref) int
	"syscall/js.valueLength": `    go.mem.StoreInt64(local0 + 16, ((Array)go.LoadValue(local0 + 8)).Length);`,

	// valuePrepareString(v ref) (ref, int)
	"syscall/js.valuePrepareString": `    byte[] str = Encoding.UTF8.GetBytes(go.LoadValue(local0 + 8).ToString());
    go.StoreValue(local0 + 16, str);
    go.mem.StoreInt64(local0 + 24, str.Length);`,

	// valueLoadString(v ref, b []byte)
	"syscall/js.valueLoadString": `    byte[] src = (byte[])go.LoadValue(local0 + 8);
    var dst = go.mem.LoadSlice(local0 + 16);
    int len = Math.Min(dst.Count, src.Length);
    for (int i = 0; i < len; i++)
    {
        dst.Array[dst.Offset + i] = src[i];
    }`,

	/*// func valueInstanceOf(v ref, t ref) bool
	"syscall/js.valueInstanceOf": (sp) => {
		this.mem.setUint8(sp + 24, loadValue(sp + 8) instanceof loadValue(sp + 16));
	},*/

	// func copyBytesToGo(dst []byte, src ref) (int, bool)
	"syscall/js.copyBytesToGo": `    var dst = go.mem.LoadSlice(local0 + 8);
    var src = go.LoadValue(local0 + 32);
    if (!(src is byte[]))
    {
        go.mem.StoreInt8(local0 + 48, 0);
        return;
    }
    var srcbs = (byte[])src;
    for (int i = 0; i < dst.Count; i++)
    {
        dst.Array[dst.Offset + i] = srcbs[i];
    }
    go.mem.StoreInt64(local0 + 40, (long)dst.Count);
    go.mem.StoreInt8(local0 + 48, 1);`,

	// func copyBytesToJS(dst ref, src []byte) (int, bool)
	"syscall/js.copyBytesToJS": `    var dst = go.LoadValue(local0 + 8);
    var src = go.mem.LoadSlice(local0 + 16);
    if (!(dst is byte[]))
    {
        go.mem.StoreInt8(local0 + 48, 0);
        return;
    }
    var dstbs = (byte[])dst;
    for (int i = 0; i < dstbs.Length; i++)
    {
        dstbs[i] = src.Array[src.Offset + i];
    }
    go.mem.StoreInt64(local0 + 40, (long)dstbs.Length);
    go.mem.StoreInt8(local0 + 48, 1);`,

	"debug": `    Console.WriteLine(local0);`,
}
