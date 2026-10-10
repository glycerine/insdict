//go:build insdict_simd && amd64 && !purego && !insdict_swar

#include "textflag.h"

// func groupMasks(ctrl *[16]byte, tag byte) uint64
TEXT ·groupMasks(SB), NOSPLIT, $0-24
	MOVQ ctrl+0(FP), AX
	MOVOU (AX), X0
	MOVBLZX tag+8(FP), AX
	MOVD AX, X1
	PUNPCKLBW X1, X1
	PUNPCKLWL X1, X1
	PSHUFL $0, X1, X1
	PCMPEQB X0, X1
	PMOVMSKB X1, AX
	PXOR X2, X2
	PCMPEQB X0, X2
	PMOVMSKB X2, CX
	PMOVMSKB X0, DX
	XORL $0xffff, DX
	SHLQ $16, CX
	SHLQ $32, DX
	ORQ CX, AX
	ORQ DX, AX
	MOVQ AX, ret+16(FP)
	RET
