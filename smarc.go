package smarc

// #cgo LDFLAGS: -L${SRCDIR} -lsmarc
// #include <smarc.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

// PSFilter --
type PSFilter struct {
	L        int
	M        int
	fpass    float64
	fstop    float64
	rp       float64
	rs       float64
	rpFactor int
}

// PFilter --
type PFilter struct {
	fsin      int
	fsout     int
	fpass     float64
	fstop     float64
	rp        float64
	rs        float64
	nb_stages int
	filter    *PSFilter
}

// PStageBuffer --
type PStageBuffer struct {
	data []float64
	size int
	pos  int
}

// PState --
type PState struct {
	nb_stages int
	state     *PState
	buffer    *PStageBuffer
	// flush vars
	flush_buf   []float64
	flush_size  int
	flush_pos   int
	flush_stage int
}

// Resample --
func Resample(buf []float64, inRate int, outRate int, bandwidth float64, rippleFactor float64, stopbandAttenuation float64, tolerance float64) []float64 {
	// make an input buffer with size of buf
	inBuf := make([]C.double, len(buf))
	//	range over buf to convert to c type
	for idx, el := range buf {
		inBuf[idx] = C.double(el)
	}
	// scale the size of the output buffer
	outBufLen := len(buf) * ((outRate / inRate) + 1)
	outBuf := make([]C.double, outBufLen)
	fsin := C.int(inRate)   // input samplerate
	fsout := C.int(outRate) // output samplerate

	// hard coded for now
	bw := C.double(bandwidth)           // bandwidth
	rp := C.double(rippleFactor)        // passband ripple factor
	rs := C.double(stopbandAttenuation) // stopband attenuation
	tol := C.double(tolerance)          // tolerance

	pFilt := (*PFilter)(unsafe.Pointer(C.smarc_init_pfilter(fsin, fsout, bw, rp, rs, tol, nil, C.int(0))))
	pState := (*PState)(unsafe.Pointer(C.smarc_init_pstate((*C.struct_PFilter)(unsafe.Pointer(pFilt)))))

	written := C.smarc_resample((*C.struct_PFilter)(unsafe.Pointer(pFilt)), (*C.struct_PState)(unsafe.Pointer(pState)), &inBuf[0], C.int(len(buf)), &outBuf[0], C.int(outBufLen))
	written += C.smarc_resample_flush((*C.struct_PFilter)(unsafe.Pointer(pFilt)), (*C.struct_PState)(unsafe.Pointer(pState)), &outBuf[0], C.int(outBufLen))

	C.smarc_reset_pstate((*C.struct_PState)(unsafe.Pointer(pState)), (*C.struct_PFilter)(unsafe.Pointer(pFilt)))
	C.smarc_destroy_pfilter((*C.struct_PFilter)(unsafe.Pointer(pFilt)))
	C.smarc_destroy_pstate((*C.struct_PState)(unsafe.Pointer(pState)))
	out := make([]float64, outBufLen)
	for idx, el := range outBuf {
		out[idx] = float64(el)
	}

	return out[:written]
}
